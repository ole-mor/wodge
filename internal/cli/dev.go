package cli

import (
	"fmt"
	"net"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"wodge/internal/generator"
	"wodge/internal/registry"
	"wodge/internal/server"

	// Added import
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/cobra"
)

var devCmd = &cobra.Command{
	Use:    "dev [app_name]",
	Short:  "Run the Wodge application in development mode",
	Long:   `Run the Wodge application in development mode. If [app_name] is provided, it switches to that directory first.`,
	Args:   cobra.MaximumNArgs(1),
	Hidden: true, // Prefer 'wodge run [app] dev'
	Run:    runDev,
}

func runDev(cmd *cobra.Command, args []string) {
	if len(args) > 0 {
		targetDir := args[0]
		if err := os.Chdir(targetDir); err != nil {
			fmt.Printf("Error switching to directory '%s': %v\n", targetDir, err)
			os.Exit(1)
		}
		fmt.Printf("Switched to directory: %s\n", targetDir)
	}

	fmt.Println("Starting Wodge development server...")

	// 0. Register App
	cwd, _ := os.Getwd()
	appName := filepath.Base(cwd)

	// 1. Find a free port for Wodge backend (start at 8080)
	// If 8080 is taken (e.g., by AstAuth), it will try 8081, etc.
	port := findAvailablePort(8080)
	fmt.Printf("DEBUG: findAvailablePort returned port %d\n", port)

	// Write port to wodge client file so frontend knows where to look
	updateEnvPort(cwd, port)
	updateWodgeClient(cwd, port)

	reg, err := registry.Load()
	if err == nil {
		if err := reg.Register(appName, port, cwd); err != nil {
			fmt.Printf("Warning: Failed to register app: %v\n", err)
		} else {
			fmt.Printf("Registered app '%s' on port %d\n", appName, port)
			defer func() {
				reg.Unregister(appName)
				fmt.Printf("Unregistered app '%s'\n", appName)
			}()
		}
	}

	// Handle graceful shutdown to ensure defer runs
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
	}()

	// ... continue with generation and watcher ...

	// 1. Generate routes initially
	err = generator.GenerateRoutes("src")
	if err != nil {
		fmt.Printf("Warning: Failed to generate routes: %v\n", err)
	} else {
		fmt.Println("Routes generated successfully.")
	}

	// 2. Setup File Watcher
	var watcher *fsnotify.Watcher
	watcher, err = fsnotify.NewWatcher()
	if err != nil {
		fmt.Printf("Error creating file watcher: %v\n", err)
		os.Exit(1)
	}
	defer watcher.Close()

	routesDir := "src/routes"
	if _, err := os.Stat(routesDir); os.IsNotExist(err) {
		_ = os.MkdirAll(routesDir, 0755)
	}

	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				if event.Op&fsnotify.Write == fsnotify.Write ||
					event.Op&fsnotify.Create == fsnotify.Create ||
					event.Op&fsnotify.Remove == fsnotify.Remove ||
					event.Op&fsnotify.Rename == fsnotify.Rename {

					fmt.Println("Route change detected. Regenerating routes...")
					if err := generator.GenerateRoutes("src"); err != nil {
						fmt.Printf("Error generating routes: %v\n", err)
					}
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				fmt.Printf("Watcher error: %v\n", err)
			}
		}
	}()

	if err := watcher.Add(routesDir); err != nil {
		fmt.Printf("Error adding watcher to %s: %v\n", routesDir, err)
	} else {
		fmt.Printf("Watching %s for changes...\n", routesDir)
	}

	// 3. Start Go API Server
	fmt.Printf("Starting API server on port %d...\n", port)
	startBackend(cwd, port)

	// 4. Start Vite dev server
	// In dev mode, always start Vite regardless of dist/ folder existence.
	// The dist/ folder is used for production builds, not dev mode.
	viteCmd := exec.Command("npx", "vite")
	viteCmd.Stdout = os.Stdout
	viteCmd.Stderr = os.Stderr
	viteCmd.Stdin = os.Stdin

	fmt.Println("Running Vite...")
	if err := viteCmd.Run(); err != nil {
		fmt.Printf("Vite exited: %v\n", err)
	}
}

func findAvailablePort(startPort int) int {
	port := startPort
	for {
		// Try to listen on the port to see if it's available
		listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
		if err == nil {
			listener.Close()
			return port
		}
		port++
		if port > startPort+100 {
			return startPort // Fallback to default
		}
	}
}

func startBackend(appPath string, port int) *os.Process {
	// Load environment variables from app's .env
	loadEnv(appPath)

	// Force PORT env var for the server to pick up
	os.Setenv("PORT", fmt.Sprintf("%d", port))

	go func() {
		server.Start(port)
	}()

	return nil
}

func loadEnv(appPath string) {
	envFile := filepath.Join(appPath, ".env")
	content, err := os.ReadFile(envFile)
	if err == nil {
		lines := strings.Split(string(content), "\n")
		for _, line := range lines {
			parts := strings.SplitN(line, "=", 2)
			if len(parts) == 2 {
				key := strings.TrimSpace(parts[0])
				val := strings.TrimSpace(parts[1])
				// Expand value? For now raw
				val = strings.Trim(val, `"'`)
				os.Setenv(key, val)
			}
		}
	}
}

func updateEnvPort(appPath string, port int) {
	envFile := filepath.Join(appPath, ".env")
	content, err := os.ReadFile(envFile)
	var lines []string

	if err == nil {
		lines = strings.Split(string(content), "\n")
	}

	found := false
	newLines := []string{}
	for _, line := range lines {
		if strings.HasPrefix(strings.TrimSpace(line), "PORT=") {
			newLines = append(newLines, fmt.Sprintf("PORT=%d", port))
			found = true
		} else {
			newLines = append(newLines, line)
		}
	}

	if !found {
		newLines = append(newLines, fmt.Sprintf("PORT=%d", port))
	}

	// Write back
	finalContent := strings.Join(newLines, "\n")
	if finalContent != "" && !strings.HasSuffix(finalContent, "\n") {
		finalContent += "\n"
	}
	os.WriteFile(envFile, []byte(finalContent), 0644)
}

func updateWodgeClient(appPath string, port int) {
	// Deprecated: We now use environment variables (VITE_API_BASE)
	// which are handled by .env files and Vite config.
	// We no longer overwrite the source code.
}
