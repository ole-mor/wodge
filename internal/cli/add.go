package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add new routes or APIs to a Wodge app",
}

var addAPICmd = &cobra.Command{
	Use:   "api [name]",
	Short: "Add a new API route to the app",
	Args:  cobra.ExactArgs(1),
	Run:   runAddAPI,
}

func init() {
	addCmd.AddCommand(addAPICmd)
}

func runAddAPI(cmd *cobra.Command, args []string) {
	apiName := args[0]

	// Find the app root by looking for src/routes directory
	appRoot, err := findAppRoot()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		fmt.Println("Please run this command from within a Wodge app directory (where src/routes exists)")
		os.Exit(1)
	}

	fmt.Printf("Found app at: %s\n", appRoot)
	fmt.Printf("Creating API: %s\n", apiName)

	// Create the /api directory if it doesn't exist
	apiDir := filepath.Join(appRoot, "src", "api")
	if err := os.MkdirAll(apiDir, 0755); err != nil {
		fmt.Printf("Error creating api directory: %v\n", err)
		os.Exit(1)
	}

	// Create the route file
	routeFileName := fmt.Sprintf("%s.route.ts", apiName)
	routePath := filepath.Join(apiDir, routeFileName)

	// Check if file already exists
	if _, err := os.Stat(routePath); !os.IsNotExist(err) {
		fmt.Printf("Error: API '%s' already exists\n", apiName)
		os.Exit(1)
	}

	// Generate the route template
	routeContent := generateAPIRoute(apiName)

	// Write the file
	if err := os.WriteFile(routePath, []byte(routeContent), 0644); err != nil {
		fmt.Printf("Error writing route file: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Created %s\n", filepath.Join("src/api", routeFileName))
	fmt.Println("The routes will be regenerated on next save")
}

func findAppRoot() (string, error) {
	// Start from current directory and walk up looking for src/routes
	currentDir, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("could not get current directory: %v", err)
	}

	for {
		routesPath := filepath.Join(currentDir, "src", "routes")
		if _, err := os.Stat(routesPath); !os.IsNotExist(err) {
			return currentDir, nil
		}

		parent := filepath.Dir(currentDir)
		if parent == currentDir {
			// Reached root directory
			return "", fmt.Errorf("could not find Wodge app root (no src/routes directory found)")
		}
		currentDir = parent
	}
}

func generateAPIRoute(name string) string {
	// Create a simple API route handler that delegates to the backend
	return fmt.Sprintf(`import { apiGet, apiPost } from '@/lib/wodge';

// Delegate GET requests to the backend API
export async function GET(req: Request) {
  try {
    const data = await apiGet('/%s');
    return new Response(JSON.stringify(data), {
      status: 200,
      headers: { 'Content-Type': 'application/json' },
    });
  } catch (error) {
    return new Response(
      JSON.stringify({ error: error instanceof Error ? error.message : 'Internal Server Error' }),
      { status: 500, headers: { 'Content-Type': 'application/json' } }
    );
  }
}

// Delegate POST requests to the backend API
export async function POST(req: Request) {
  try {
    const body = await req.json();
    const data = await apiPost('/%s', body);
    return new Response(JSON.stringify(data), {
      status: 200,
      headers: { 'Content-Type': 'application/json' },
    });
  } catch (error) {
    return new Response(
      JSON.stringify({ error: error instanceof Error ? error.message : 'Invalid request' }),
      { status: 400, headers: { 'Content-Type': 'application/json' } }
    );
  }
}
`, name, name)
}

func toPascalCase(s string) string {
	parts := strings.Split(s, "-")
	for i, part := range parts {
		if len(part) > 0 {
			parts[i] = strings.ToUpper(string(part[0])) + part[1:]
		}
	}
	return strings.Join(parts, "")
}
