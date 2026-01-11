package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"unicode"

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
If name is 'health', it adds a health check client.
If name is 'postgres', 'redis', or 'rabbitmq', it adds a client library for that service.`,
	Args: cobra.RangeArgs(1, 2),
	Run:  runAddAPI,
}

func init() {
	addCmd.AddCommand(addAPICmd)
	addCmd.AddCommand(uiCmd)
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
	case "qast":
		addQastClient(appRoot)
	case "health":
		addHealthRoute(appRoot)
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
		"src/api/postgres.ts": `import { apiPost } from '@/lib/wodge';

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

	// Inject default env var
	updateEnvFile(appRoot, "POSTGRES_DSN", "postgres://postgres:postgres@127.0.0.1:5432/postgres?sslmode=disable")

	fmt.Println("Postgres client added to src/api/postgres.ts")
	fmt.Println("Added POSTGRES_DSN to .env")
	fmt.Println("\nTip: Run Postgres locally with Docker:")
	fmt.Println("  docker run --name postgres -e POSTGRES_PASSWORD=postgres -p 5432:5432 -d postgres")
}

func addHealthRoute(appRoot string) {
	fmt.Println("Adding Health Client...")
	files := map[string]string{
		"src/api/health.ts": `import { apiGet } from '@/lib/wodge';

export const HealthService = {
  async check(): Promise<{ status: string }> {
    return apiGet('/health');
  }
};
`,
	}
	writeFiles(appRoot, files)
	fmt.Println("Health client added to src/api/health.ts")
}

func addRedisClient(appRoot string) {
	fmt.Println("Adding Redis Client...")
	files := map[string]string{
		"src/api/redis.ts": `import { apiGet, apiPost, apiDelete } from '@/lib/wodge';

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

	updateEnvFile(appRoot, "REDIS_ADDR", "127.0.0.1:6379")
	updateEnvFile(appRoot, "REDIS_PASSWORD", "")

	fmt.Println("Redis client added to src/api/redis.ts")
	fmt.Println("Added REDIS_ADDR and REDIS_PASSWORD to .env")
	fmt.Println("\nTip: Run Redis locally with Docker:")
	fmt.Println("  docker run --name redis -p 6379:6379 -d redis")
}

func addRabbitMQClient(appRoot string) {
	fmt.Println("Adding RabbitMQ Client...")
	files := map[string]string{
		"src/api/rabbitmq.ts": `import { apiPost } from '@/lib/wodge';

export const rabbitmq = {
  async publish(topic: string, message: string): Promise<void> {
    return apiPost('/queue/publish', { topic, message });
  }
};
`,
	}
	writeFiles(appRoot, files)

	updateEnvFile(appRoot, "RABBITMQ_URL", "amqp://guest:guest@127.0.0.1:5672/")

	fmt.Println("RabbitMQ client added to src/api/rabbitmq.ts")
	fmt.Println("Added RABBITMQ_URL to .env")
	fmt.Println("\nTip: Run RabbitMQ locally with Docker:")
	fmt.Println("  docker run --name rabbitmq -p 5672:5672 -d rabbitmq")
}

func addQastClient(appRoot string) {
	fmt.Println("Adding QAST Client...")
	files := map[string]string{
		"src/api/qast.ts": `import { apiPost } from '@/lib/wodge';

export const qast = {
  // RAG Search via Composer
  async ask(query: string, userId: string = "default-user", expertise: string = "novice"): Promise<{ answer: string; context: string[] }> {
    return apiPost('/qast/ask', { query, user_id: userId, expertise_level: expertise });
  },

  // Secure PII Chat via Privacy (Rehydration handles the real PII)
  async chat(text: string, userId: string = "default-user"): Promise<{ llm_response: string; token_map: Record<string, string> }> {
    return apiPost('/qast/chat', { text, user_id: userId });
  },

  async ingest(text: string, userId: string = "default-user"): Promise<{ status: string; result: any }> {
    return apiPost('/qast/ingest', { text, user_id: userId });
  }
};
`,
	}
	writeFiles(appRoot, files)

	// Default QAST URL (Proxy)
	updateEnvFile(appRoot, "QAST_URL", "http://localhost:9988")

	fmt.Println("QAST client added to src/api/qast.ts")
	fmt.Println("Added QAST_URL to .env")
	fmt.Println("Make sure Qast-Link is running!")
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

func updateEnvFile(root, key, value string) {
	envPath := filepath.Join(root, ".env")

	// Create if not exists
	if _, err := os.Stat(envPath); os.IsNotExist(err) {
		os.WriteFile(envPath, []byte(""), 0644)
	}

	contentBytes, err := os.ReadFile(envPath)
	if err != nil {
		fmt.Printf("Warning: could not read .env file: %v\n", err)
		return
	}
	content := string(contentBytes)

	if strings.Contains(content, key+"=") {
		fmt.Printf("Note: %s already exists in .env, skipping.\n", key)
		return
	}

	f, err := os.OpenFile(envPath, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Printf("Warning: could not open .env file: %v\n", err)
		return
	}
	defer f.Close()

	if len(content) > 0 && !strings.HasSuffix(content, "\n") {
		f.WriteString("\n")
	}
	f.WriteString(fmt.Sprintf("%s=%s\n", key, value))
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
	// Generate a type-safe Service Object for the entity
	// This fits the client-side nature of Wodge (Vite) better than Request/Response handlers
	return fmt.Sprintf(`import { postgres } from '@/api/postgres';

export const %sService = {
  async list() {
    return postgres.query('SELECT * FROM %s');
  },

  async get(id: string) {
    const rows = await postgres.query('SELECT * FROM %s WHERE id = $1', [id]);
    return rows[0] || null;
  },

  async create(data: Record<string, any>) {
    const keys = Object.keys(data);
    const values = Object.values(data);
    const placeholders = keys.map((_, i) => '$' + (i + 1)).join(', ');
    const columns = keys.join(', ');
    
    const query = 'INSERT INTO %s (' + columns + ') VALUES (' + placeholders + ')';
    return postgres.execute(query, values);
  },

  async delete(id: string) {
    return postgres.execute('DELETE FROM %s WHERE id = $1', [id]);
  }
};
`, toPascalCase(name), name, name, name, name)
}

func toPascalCase(s string) string {
	parts := strings.FieldsFunc(s, func(r rune) bool {
		return !unicode.IsLetter(r) && !unicode.IsNumber(r)
	})
	var result string
	for _, part := range parts {
		if len(part) > 0 {
			runes := []rune(part)
			runes[0] = unicode.ToUpper(runes[0])
			result += string(runes)
		}
	}
	return result
}
