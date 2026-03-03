package cli

import (
	"fmt"
	"os"
	"strconv"
	"wodge/internal/server"

	"github.com/spf13/cobra"
)

var serveCmd = &cobra.Command{
	Use:   "serve [port]",
	Short: "Start the Wodge API server and serve static files from ./dist",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		port := 8080
		if len(args) > 0 {
			p, err := strconv.Atoi(args[0])
			if err != nil {
				fmt.Printf("Invalid port: %v\n", err)
				os.Exit(1)
			}
			port = p
		} else {
			// Check for PORT env var
			if pStr := os.Getenv("PORT"); pStr != "" {
				if p, err := strconv.Atoi(pStr); err == nil {
					port = p
				}
			}
		}

		fmt.Printf("Starting Wodge Production Server on port %d...\n", port)
		server.Start(port)
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)
}
