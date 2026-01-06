package cli

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add new routes or APIs to a Wodge app",
}

var addAPICmd = &cobra.Command{
	Use:   "api [name] OR api crud [name]",
	Short: "Add a new API route or service client to the app",
	Long: `Adds a new API. 
If 'crud [name]' is specified, it creates a simple CRUD skeleton.
If name is 'postgres', 'redis', or 'rabbitmq', it adds a client library for that service.`,
	Args: cobra.RangeArgs(1, 2),
	Run:  runAddAPI,
}

func init() {
	addCmd.AddCommand(addAPICmd)
}

func runAddAPI(cmd *cobra.Command, args []string) {
	// Handle special sub-case: crud
	if args[0] == "crud" {
		if len(args) < 2 {
			fmt.Println("Error: Please specify a name for the CRUD API. usage: wodge add api crud <name>")
			os.Exit(1)
		}
		apiName := args[1]
		appRoot, err := findAppRoot()
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}
		addCRUDRoute(appRoot, apiName)
		return
	}

	apiName := args[0]
	appRoot, err := findAppRoot()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		fmt.Println("Please run this command from within a Wodge app directory")
		os.Exit(1)
	}

	// Handle Predefined Services
	switch apiName {
	case "postgres":
		addPostgresClient(appRoot)
	case "redis":
		addRedisClient(appRoot)
	case "rabbitmq":
		addRabbitMQClient(appRoot)
	default:
		// Default behavior: Create a generic API route
		addGenericAPIRoute(appRoot, apiName)
	}
}

func addGenericAPIRoute(appRoot, apiName string) {
	fmt.Printf("Creating API: %s\n", apiName)

	apiDir := filepath.Join(appRoot, "src", "api")
	if err := os.MkdirAll(apiDir, 0755); err != nil {
		fmt.Printf("Error creating api directory: %v\n", err)
		os.Exit(1)
	}

	routeFileName := fmt.Sprintf("%s.route.ts", apiName)
	routePath := filepath.Join(apiDir, routeFileName)

	if _, err := os.Stat(routePath); !os.IsNotExist(err) {
		fmt.Printf("Error: API '%s' already exists\n", apiName)
		os.Exit(1)
	}

	routeContent := generateAPIRoute(apiName)
	if err := os.WriteFile(routePath, []byte(routeContent), 0644); err != nil {
		fmt.Printf("Error writing route file: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Created %s\n", filepath.Join("src/api", routeFileName))
	fmt.Println("The routes will be regenerated on next save")
}

func addCRUDRoute(appRoot, apiName string) {
	fmt.Printf("Creating CRUD API: %s\n", apiName)

	apiDir := filepath.Join(appRoot, "src", "api")
	if err := os.MkdirAll(apiDir, 0755); err != nil {
		fmt.Printf("Error creating api directory: %v\n", err)
		os.Exit(1)
	}

	routeFileName := fmt.Sprintf("%s.crud.route.ts", apiName)
	routePath := filepath.Join(apiDir, routeFileName)

	if _, err := os.Stat(routePath); !os.IsNotExist(err) {
		fmt.Printf("Error: CRUD API '%s' already exists\n", apiName)
		os.Exit(1)
	}

	routeContent := generateCRUDApiRoute(apiName)
	if err := os.WriteFile(routePath, []byte(routeContent), 0644); err != nil {
		fmt.Printf("Error writing route file: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Created %s\n", filepath.Join("src/api", routeFileName))
	fmt.Println("The routes will be regenerated on next save")
}

func addPostgresClient(appRoot string) {
	fmt.Println("Adding Postgres Client...")
	files := map[string]string{
		"src/lib/postgres.ts": `import { apiPost } from './wodge';

export interface QueryResult<T = any> {
  [key: string]: any;
}

export const postgres = {
  /**
   * Execute a SELECT query
   */
  async query<T = any>(query: string, args: any[] = []): Promise<T[]> {
    return apiPost('/postgres/query', { query, args });
  },

  /**
   * Execute an INSERT/UPDATE/DELETE query
   */
  async execute(query: string, args: any[] = []): Promise<{ rows_affected: number }> {
    return apiPost('/postgres/execute', { query, args });
  }
};
`,
	}
	writeFiles(appRoot, files)
	fmt.Println("Postgres client added to src/lib/postgres.ts")
	fmt.Println("Make sure POSTGRES_DSN is set in your environment variables.")
}

func addRedisClient(appRoot string) {
	fmt.Println("Adding Redis Client...")
	files := map[string]string{
		"src/lib/redis.ts": `import { apiGet, apiPost, apiDelete } from './wodge';

export const redis = {
  async get(key: string): Promise<string | null> {
    try {
        const res = await apiGet('/redis/' + encodeURIComponent(key));
        return res.value;
    } catch (e) {
        return null;
    }
  },

  async set(key: string, value: string, ttl: number = 0): Promise<void> {
    return apiPost('/redis', { key, value, ttl });
  },

  async delete(key: string): Promise<void> {
    return apiDelete('/redis/' + encodeURIComponent(key));
  }
};
`,
	}
	writeFiles(appRoot, files)
	fmt.Println("Redis client added to src/lib/redis.ts")
	fmt.Println("Make sure REDIS_ADDR (and optionally REDIS_PASSWORD) is set.")
}

func addRabbitMQClient(appRoot string) {
	fmt.Println("Adding RabbitMQ Client...")
	files := map[string]string{
		"src/lib/rabbitmq.ts": `import { apiPost } from './wodge';

export const rabbitmq = {
  async publish(topic: string, message: string): Promise<void> {
    return apiPost('/queue/publish', { topic, message });
  }
};
`,
	}
	writeFiles(appRoot, files)
	fmt.Println("RabbitMQ client added to src/lib/rabbitmq.ts")
	fmt.Println("Make sure RABBITMQ_URL is set.")
}

func writeFiles(root string, files map[string]string) {
	for path, content := range files {
		fullPath := filepath.Join(root, path)
		dir := filepath.Dir(fullPath)
		_ = os.MkdirAll(dir, 0755)
		if err := os.WriteFile(fullPath, []byte(content), 0644); err != nil {
			fmt.Printf("Error writing %s: %v\n", path, err)
		}
	}
}

func findAppRoot() (string, error) {
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
			return "", fmt.Errorf("could not find Wodge app root (no src/routes directory found)")
		}
		currentDir = parent
	}
}

func generateAPIRoute(name string) string {
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

func generateCRUDApiRoute(name string) string {
	// Generate a file with GET (List), GET (ByID), POST (Create), PUT (Update), DELETE handlers
	return fmt.Sprintf(`// Mock DB for demonstration (in real app, use postgres/redis clients)
const DB = new Map<string, any>();

// GET /api/%s
export async function GET(req: Request) {
  const url = new URL(req.url);
  const id = url.searchParams.get('id');

  try {
    if (id) {
       // Get One
       const item = DB.get(id);
       if (!item) return new Response(JSON.stringify({ error: 'Not Found' }), { status: 404 });
       return new Response(JSON.stringify(item));
    } else {
       // List All
       const items = Array.from(DB.values());
       return new Response(JSON.stringify(items));
    }
  } catch (error) {
    return new Response(JSON.stringify({ error: 'Internal Error' }), { status: 500 });
  }
}

// POST /api/%s
export async function POST(req: Request) {
  try {
    const body = await req.json();
    const id = Math.random().toString(36).substring(7);
    const item = { id, ...body, createdAt: new Date() };
    DB.set(id, item);
    return new Response(JSON.stringify(item), { status: 201 });
  } catch (error) {
    return new Response(JSON.stringify({ error: 'Invalid Data' }), { status: 400 });
  }
}

// DELETE /api/%s
export async function DELETE(req: Request) {
  const url = new URL(req.url);
  const id = url.searchParams.get('id');
  if (!id) return new Response(JSON.stringify({ error: 'ID required' }), { status: 400 });
  
  DB.delete(id);
  return new Response(JSON.stringify({ status: 'deleted' }));
}
`, name, name, name)
}
