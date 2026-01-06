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
	appName := args[0]
	command := "dev" // Default command
	if len(args) > 1 {
		command = args[1]
	}

	// 1. Switch Directory
	if err := os.Chdir(appName); err != nil {
		fmt.Printf("Error: Could not find application directory '%s': %v\n", appName, err)
		os.Exit(1)
	}

	// 2. Dispatch Command
	switch command {
	case "dev":
		fmt.Printf("Running 'dev' for %s...\n", appName)
		// We can reuse the logic from dev.go, but need to be careful about Cobra context.
		// Since runDev expects *cobra.Command and []string, we might need to adapt it
		// or refactor dev.go to expose a standalone function.
		// For now, let's call runDev directly with the original args.
		runDev(cmd, []string{})
	default:
		fmt.Printf("Unknown command '%s'. Available commands: dev\n", command)
		os.Exit(1)
	}
}
