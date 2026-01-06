package cli

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"wodge/internal/templates"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "wodge",
	Short: "Wodge is a tool for creating and managing NIS2-compliant web applications",
	Long:  `Wodge is a comprehensive framework and CLI tool designed to streamline the creation of React + Go web applications with built-in compliance features.`,
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.AddCommand(newCmd)
	rootCmd.AddCommand(devCmd)
	rootCmd.AddCommand(addCmd)
	rootCmd.AddCommand(updateCmd)
}

var newCmd = &cobra.Command{
	Use:   "new [name]",
	Short: "Create a new Wodge application",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		appName := args[0]
		fmt.Printf("Creating new Wodge app: %s\n", appName)
		createApp(appName)
	},
}

// devCmd is defined in dev.go

func createApp(name string) {
	fmt.Printf("Scaffolding %s...\n", name)
	if _, err := os.Stat(name); !os.IsNotExist(err) {
		fmt.Printf("Error: Directory '%s' already exists\n", name)
		return
	}

	if err := os.Mkdir(name, 0755); err != nil {
		fmt.Printf("Error creating directory: %v\n", err)
		return
	}

	files := map[string]string{
		"package.json":                templates.GetPackageJSON(name),
		"vite.config.ts":              templates.ViteConfig,
		"tsconfig.json":               templates.TsConfig,
		"tsconfig.node.json":          templates.TsConfigNode,
		"index.html":                  templates.IndexHTML,
		"src/entry-client.tsx":        templates.EntryClient,
		"src/entry-server.tsx":        templates.EntryServer,
		"src/App.tsx":                 templates.AppTSX,
		"src/routes.generated.tsx":    templates.RoutesGenerated,
		"src/routes/home.route.tsx":   templates.HomeRoute,
		"go.mod":                      templates.GetGoMod(name),
		"cmd/server/main.go":          fmt.Sprintf(templates.BackendMain, name),
		"internal/handlers/routes.go": templates.BackendHandlers,
	}

	for path, content := range files {
		fullPath := filepath.Join(name, path)
		dir := filepath.Dir(fullPath)
		if err := os.MkdirAll(dir, 0755); err != nil {
			fmt.Printf("Error creating directory %s: %v\n", dir, err)
			continue
		}
		if err := os.WriteFile(fullPath, []byte(content), 0644); err != nil {
			fmt.Printf("Error writing file %s: %v\n", fullPath, err)
		} else {
			fmt.Printf("Created %s\n", path)
		}
	}

	fmt.Printf("\nInstalling dependencies...\n")
	cmd := exec.Command("npm", "install")
	cmd.Dir = name
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		fmt.Printf("Error installing dependencies: %v\n", err)
		return
	}

	fmt.Printf("\nInitializing Go module...\n")
	goCmd := exec.Command("go", "mod", "tidy")
	goCmd.Dir = name
	goCmd.Stdout = os.Stdout
	goCmd.Stderr = os.Stderr
	if err := goCmd.Run(); err != nil {
		fmt.Printf("Error tidying go module: %v\n", err)
	}

	fmt.Printf("\nSuccess! Created %s\n", name)
	fmt.Printf("cd %s\n", name)
	fmt.Printf("wodge dev\n")
}
