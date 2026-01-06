package cli

import (
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"syscall"
	"wodge/internal/generator"
	"wodge/internal/registry"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/cobra"
)

func findWodgeRoot() (string, error) {
	// Start from current directory and walk up looking for cmd/api-server/main.go
	currentDir, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("could not get current directory: %v", err)
	}

	for {
		apiServerPath := filepath.Join(currentDir, "cmd", "api-server", "main.go")
		if _, err := os.Stat(apiServerPath); !os.IsNotExist(err) {
			return currentDir, nil
		}

		parent := filepath.Dir(currentDir)
		if parent == currentDir {
			// Reached root directory
			return "", fmt.Errorf("could not find Wodge root (no cmd/api-server/main.go found)")
		}
		currentDir = parent
	}
}

var devCmd = &cobra.Command{
	Use:   "dev [app_name]",
	Short: "Run the Wodge application in development mode",
	Long:  `Run the Wodge application in development mode. If [app_name] is provided, it switches to that directory first.`,
	Args:  cobra.MaximumNArgs(1),
	Run:   runDev,
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
	reg, err := registry.Load()
	if err == nil {
		if err := reg.Register(appName, 8080, cwd); err != nil {
			fmt.Printf("Warning: Failed to register app: %v\n", err)
		} else {
			fmt.Printf("Registered app '%s' in Wodge registry\n", appName)
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
		// Cleanup logic is in defers, we need to manually trigger them if we exit here,
		// but since we are blocking on viteCmd.Run() below, we might not reach here unless we architecture differently.
		// Actually, standard Ctrl+C propagates to child processes (vite) which exit, causing viteCmd.Run() to return,
		// allowing main function defers to run.
	}()

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

	// Watch src/routes for changes
	// We need to ensure the directory exists first
	routesDir := "src/routes"
	if _, err := os.Stat(routesDir); os.IsNotExist(err) {
		// Try to create it if it doesn't exist, though scaffolding should have made it
		_ = os.MkdirAll(routesDir, 0755)
	}

	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				// We care about creating, removing, or renaming files in routes dir
				// modifying might also change the export name if we parse it, but standard route structure relies on filename
				// For now, let's just regenerate on any event in that folder for simplicity
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
	go func() {
		fmt.Println("Starting API server...")
		// Find the wodge root by looking for cmd/api-server/main.go
		wodgeRoot, err := findWodgeRoot()
		if err != nil {
			fmt.Printf("Warning: Could not find Wodge root, skipping API server: %v\n", err)
			return
		}

		apiCmd := exec.Command("go", "run", "cmd/api-server/main.go")
		apiCmd.Dir = wodgeRoot
		apiCmd.Stdout = os.Stdout
		apiCmd.Stderr = os.Stderr
		if err := apiCmd.Run(); err != nil {
			fmt.Printf("API server error: %v\n", err)
		}
	}()

	// 4. Start Vite
	// We assume we are in the project root, so we check for node_modules/.bin/vite
	// or try npx vite
	viteCmd := exec.Command("npx", "vite")
	viteCmd.Stdout = os.Stdout
	viteCmd.Stderr = os.Stderr
	viteCmd.Stdin = os.Stdin

	fmt.Println("Running Vite...")
	if err := viteCmd.Run(); err != nil {
		fmt.Printf("Error running vite: %v\n", err)
		os.Exit(1)
	}
}
