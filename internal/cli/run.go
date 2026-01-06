package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var runCmd = &cobra.Command{
	Use:   "run [app_name] [command]",
	Short: "Run a command in a Wodge application context",
	Long: `Run a command for a specific Wodge application.
	
Examples:
  wodge run my-app dev      # Start dev server for my-app
  wodge run my-app          # Defaults to 'dev'`,
	Args: cobra.RangeArgs(1, 2),
	Run:  executeRun,
}

func init() {
	rootCmd.AddCommand(runCmd)
}

func executeRun(cmd *cobra.Command, args []string) {
	var appName string
	var command string

	// Case 1: wodge run (no args) -> defaulting to current dir, dev
	if len(args) == 0 {
		appName = "."
		command = "dev"
	} else if len(args) == 1 {
		// Case 2: wodge run arg1
		// Check if arg1 is a directory
		if isDir(args[0]) {
			// It is a directory: wodge run my-app -> dev my-app
			appName = args[0]
			command = "dev"
		} else {
			// It is NOT a directory (e.g. 'dev', 'build'): wodge run dev -> dev .
			appName = "."
			command = args[0]
		}
	} else {
		// Case 3: wodge run arg1 arg2 -> wodge run my-app dev
		appName = args[0]
		command = args[1]
	}

	// 1. Switch Directory if specific app provided
	if appName != "." {
		if err := os.Chdir(appName); err != nil {
			fmt.Printf("Error: Could not find application directory '%s': %v\n", appName, err)
			os.Exit(1)
		}
	} else {
		// Verify we are in a wodge app
		if _, err := findAppRoot(); err != nil {
			fmt.Println("Error: Current directory is not a Wodge application.")
			fmt.Println("Usage: wodge run <app_name> [command] OR wodge run [command] (inside app)")
			os.Exit(1)
		}
	}

	// 2. Dispatch Command
	switch command {
	case "dev":
		// Get absolute path for logging
		cwd, _ := os.Getwd()
		fmt.Printf("Starting Wodge Dev Server in: %s\n", cwd)
		runDev(cmd, []string{})
	default:
		fmt.Printf("Unknown command '%s'. Available commands: dev\n", command)
		os.Exit(1)
	}
}

func isDir(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return info.IsDir()
}
