package templates

import "fmt"

func GetPackageJSON(appName string) string {
	return fmt.Sprintf(`{
  "name": "%s",
  "version": "0.0.0",
  "type": "module",
  "scripts": {
    "dev": "wodge run",
    "build": "vite build && vite build --ssr src/entry-server.tsx --outDir dist/server",
    "lint": "eslint ."
  },
  "dependencies": {
    "react": "^18.3.1",
    "react-dom": "^18.3.1",
    "react-router-dom": "^6.23.0",
    "framer-motion": "^11.0.0",
    "@mui/material": "^5.15.0",
    "@mui/icons-material": "^5.15.0",
    "@emotion/react": "^11.11.0",
    "@emotion/styled": "^11.11.0",
    "@fontsource/assistant": "^5.0.0",
    "clsx": "^2.1.0",
    "tailwind-merge": "^2.2.0"
  },
  "devDependencies": {
    "@types/react": "^18.3.3",
    "@types/react-dom": "^18.3.0",
    "@typescript-eslint/eslint-plugin": "^7.13.1",
    "@typescript-eslint/parser": "^7.13.1",
    "@vitejs/plugin-react": "^4.3.1",
    "@types/node": "^20.12.7",
    "eslint": "^8.57.0",
    "eslint-plugin-react-hooks": "^4.6.2",
    "eslint-plugin-react-refresh": "^0.4.7",
    "typescript": "^5.2.2",
    "vite": "^5.3.1",
    "tailwindcss": "^3.4.0",
    "postcss": "^8.4.0",
    "autoprefixer": "^10.4.0"
  }
}`, appName)
}

const ViteConfig = `import { defineConfig } from 'vite';
import react from '@vitejs/plugin-react';
import path from 'path';

export default defineConfig({
  plugins: [react()],
  resolve: {
    alias: {
      '@': path.resolve(__dirname, './src'),
    },
  },
  build: {
    rollupOptions: {
      input: {
        main: './index.html',
        client: './src/entry-client.tsx',
        server: './src/entry-server.tsx'
      }
    }
  }
});`

const TsConfig = `{
  "compilerOptions": {
    "target": "ES2020",
    "useDefineForClassFields": true,
    "lib": ["ES2020", "DOM", "DOM.Iterable"],
    "module": "ESNext",
    "skipLibCheck": true,

    /* Bundler mode */
    "moduleResolution": "bundler",
    "allowImportingTsExtensions": true,
    "resolveJsonModule": true,
    "isolatedModules": true,
    "noEmit": true,
    "jsx": "react-jsx",

    /* Linting */
    "strict": true,
    "noUnusedLocals": true,
    "noUnusedParameters": true,
    "noFallthroughCasesInSwitch": true,
    
    /* Paths */
    "baseUrl": ".",
    "paths": {
      "@/*": ["src/*"]
    }
  },
  "include": ["src"],
  "references": [{ "path": "./tsconfig.node.json" }]
}`

const TsConfigNode = `{
  "compilerOptions": {
    "composite": true,
    "skipLibCheck": true,
    "module": "ESNext",
    "moduleResolution": "bundler",
    "allowSyntheticDefaultImports": true
  },
  "include": ["vite.config.ts", "tailwind.config.ts"]
}`

const IndexHTML = `<!doctype html>
<html lang="en">
  <head>
    <meta charset="UTF-8" />
    <meta name="viewport" content="width=device-width, initial-scale=1.0" />
    <title>Wodge App</title>
    <!--app-head-->
  </head>
  <body>
    <div id="root"><!--app-html--></div>
    <script type="module" src="/src/entry-client.tsx"></script>
  </body>
</html>`

const EntryClient = `import React from 'react';
import ReactDOM from 'react-dom/client';
import { BrowserRouter } from 'react-router-dom';
import App from './App';

ReactDOM.createRoot(document.getElementById('root')!).render(
  <React.StrictMode>
    <BrowserRouter>
      <App />
    </BrowserRouter>
  </React.StrictMode>
);`

const EntryServer = `import React from 'react';
import ReactDOMServer from 'react-dom/server';
import { StaticRouter } from 'react-router-dom/server';
import App from './App';

export function render(url: string) {
  return ReactDOMServer.renderToString(
    <React.StrictMode>
      <StaticRouter location={url}>
        <App />
      </StaticRouter>
    </React.StrictMode>
  );
}`

const AppTSX = `import React from 'react';
import { GeneratedRoutes } from './routes.generated';
import '@fontsource/assistant';
import './index.css';

function App() {
  return (
    <div className="min-h-screen bg-background font-sans text-foreground">
      <main className="container mx-auto p-8">
        <GeneratedRoutes />
      </main>
    </div>
  );
}

export default App;
`

const RoutesGenerated = `import React from 'react';
import { useRoutes } from 'react-router-dom';
import Home from './routes/home.route';

export const routes = [
  { path: '/', element: <Home /> }
];

export function GeneratedRoutes() {
  return useRoutes(routes);
}
`

const HomeRoute = `import React from 'react';

export default function Home() {
  return (
    <div className="max-w-6xl mx-auto space-y-12">
      {/* Hero Section */}
      <div className="flex flex-col items-center text-center space-y-6 py-12">
        <div className="text-6xl font-bold text-primary">
          WODGE
        </div>
        <h1 className="text-5xl font-extrabold tracking-tight">
          Welcome to Wodge
        </h1>
        <p className="text-xl text-muted-foreground max-w-2xl">
          A modern, compliant web application platform. Build faster with integrated services and beautiful UI components.
        </p>
      </div>

      {/* Quick Start */}
      <div className="rounded-lg border-2 border-primary/20 bg-card/50 p-8">
        <h2 className="text-2xl font-bold mb-4 flex items-center gap-2">
          <span className="text-3xl">🚀</span> Quick Start
        </h2>
        <p className="mb-4 text-foreground/80">
          Edit <code className="bg-primary/10 px-2 py-1 rounded text-primary">src/routes/home.route.tsx</code> and see changes instantly.
        </p>
      </div>

      {/* Features Grid */}
      <div>
        <h2 className="text-3xl font-bold mb-8 text-center">Add Features On-Demand</h2>
        <div className="grid md:grid-cols-2 gap-6">
          
          {/* UI Components Card */}
          <div className="rounded-lg border border-border bg-card p-6 hover:border-primary/50 transition-colors">
            <div className="flex items-start gap-4">
              <span className="text-4xl">🎨</span>
              <div className="flex-1">
                <h3 className="text-xl font-semibold mb-3">UI Components</h3>
                <p className="text-sm text-muted-foreground mb-4">
                  Add beautiful, pre-styled components with Tailwind CSS and Framer Motion
                </p>
                <div className="space-y-2">
                  <code className="block bg-muted px-3 py-2 rounded text-sm">wodge add ui button</code>
                  <code className="block bg-muted px-3 py-2 rounded text-sm">wodge add ui card</code>
                  <code className="block bg-muted px-3 py-2 rounded text-sm">wodge add ui navbar</code>
                  <code className="block bg-muted px-3 py-2 rounded text-sm">wodge add ui theme-provider</code>
                </div>
              </div>
            </div>
          </div>

          {/* Backend Services Card */}
          <div className="rounded-lg border border-border bg-card p-6 hover:border-primary/50 transition-colors">
            <div className="flex items-start gap-4">
              <span className="text-4xl">💾</span>
              <div className="flex-1">
                <h3 className="text-xl font-semibold mb-3">Backend Services</h3>
                <p className="text-sm text-muted-foreground mb-4">
                  Integrate production-ready services with a single command
                </p>
                <div className="space-y-2">
                  <code className="block bg-muted px-3 py-2 rounded text-sm">wodge add api postgres</code>
                  <code className="block bg-muted px-3 py-2 rounded text-sm">wodge add api redis</code>
                  <code className="block bg-muted px-3 py-2 rounded text-sm">wodge add api rabbitmq</code>
                  <code className="block bg-muted px-3 py-2 rounded text-sm">wodge add api qast</code>
                </div>
              </div>
            </div>
          </div>

          {/* CRUD APIs Card */}
          <div className="rounded-lg border border-border bg-card p-6 hover:border-primary/50 transition-colors">
            <div className="flex items-start gap-4">
              <span className="text-4xl">⚡</span>
              <div className="flex-1">
                <h3 className="text-xl font-semibold mb-3">CRUD APIs</h3>
                <p className="text-sm text-muted-foreground mb-4">
                  Generate complete CRUD endpoints with frontend integration
                </p>
                <div className="space-y-2">
                  <code className="block bg-muted px-3 py-2 rounded text-sm">wodge add api crud users</code>
                  <code className="block bg-muted px-3 py-2 rounded text-sm">wodge add api crud products</code>
                </div>
              </div>
            </div>
          </div>

          {/* Health Check Card */}
          <div className="rounded-lg border border-border bg-card p-6 hover:border-primary/50 transition-colors">
            <div className="flex items-start gap-4">
              <span className="text-4xl">❤️</span>
              <div className="flex-1">
                <h3 className="text-xl font-semibold mb-3">Health Monitoring</h3>
                <p className="text-sm text-muted-foreground mb-4">
                  Add health check endpoints for monitoring
                </p>
                <div className="space-y-2">
                  <code className="block bg-muted px-3 py-2 rounded text-sm">wodge add api health</code>
                </div>
              </div>
            </div>
          </div>

        </div>
      </div>

      {/* Tech Stack */}
      <div className="rounded-lg border border-border bg-card p-8">
        <h2 className="text-2xl font-bold mb-6 text-center">Built With Modern Tech</h2>
        <div className="grid grid-cols-2 md:grid-cols-4 gap-4 text-center">
          <div className="p-4">
            <div className="text-3xl mb-2">⚛️</div>
            <div className="font-semibold">React 18</div>
            <div className="text-xs text-muted-foreground">Frontend</div>
          </div>
          <div className="p-4">
            <div className="text-3xl mb-2">🎨</div>
            <div className="font-semibold">Tailwind CSS</div>
            <div className="text-xs text-muted-foreground">Styling</div>
          </div>
          <div className="p-4">
            <div className="text-3xl mb-2">🔷</div>
            <div className="font-semibold">TypeScript</div>
            <div className="text-xs text-muted-foreground">Type Safety</div>
          </div>
          <div className="p-4">
            <div className="text-3xl mb-2">⚙️</div>
            <div className="font-semibold">Go</div>
            <div className="text-xs text-muted-foreground">Backend</div>
          </div>
        </div>
      </div>

    </div>
  );
}
`

const WodgeClientTS = `const API_BASE = 'http://localhost:8080/api';

export async function apiGet<T = any>(path: string): Promise<T> {
  const res = await fetch(API_BASE + path, {
    method: 'GET',
    headers: { 'Content-Type': 'application/json' },
  });
  if (!res.ok) {
    const err = await res.json().catch(() => ({ error: res.statusText }));
    throw new Error(err.error || res.statusText);
  }
  return res.json();
}

export async function apiPost<T = any>(path: string, body: any): Promise<T> {
  const res = await fetch(API_BASE + path, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(body),
  });
  if (!res.ok) {
    const err = await res.json().catch(() => ({ error: res.statusText }));
    throw new Error(err.error || res.statusText);
  }
  return res.json();
}

export async function apiDelete<T = any>(path: string): Promise<T> {
  const res = await fetch(API_BASE + path, {
    method: 'DELETE',
    headers: { 'Content-Type': 'application/json' },
  });
  if (!res.ok) {
    const err = await res.json().catch(() => ({ error: res.statusText }));
    throw new Error(err.error || res.statusText);
  }
  return res.json();
}
`

func GetGoMod(appName string) string {
	return fmt.Sprintf(`module %s

go 1.21

require (
	github.com/gin-gonic/gin v1.9.1
)
`, appName)
}

const BackendMain = `package main

import (
	"log"
	
	"github.com/gin-gonic/gin"
	
	"%s/internal/handlers"
)

func main() {
	r := gin.Default()
	
	// Register generated routes
	handlers.RegisterRoutes(r)

	log.Println("Starting backend on :8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
`

const BackendHandlers = `package handlers

import (
	"net/http"
	"github.com/gin-gonic/gin"
)

// RegisterRoutes is called by main to setup routes
// This file will be updated by wodge generator
func RegisterRoutes(r *gin.Engine) {
	api := r.Group("/api")
	
	api.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})
}
`

const GitIgnore = `# Logs
logs
*.log
npm-debug.log*
yarn-debug.log*
yarn-error.log*
pnpm-debug.log*
lerna-debug.log*

# Diagnostic reports (https://nodejs.org/api/report.html)
report.[0-9]*.[0-9]*.[0-9]*.[0-9]*.json

# Runtime data
pids
*.pid
*.seed
*.pid.lock

# Directory for instrumented libs generated by jscoverage/JSCover
lib-cov

# Coverage directory used by tools like istanbul
coverage
*.lcov

# nyc test coverage
.nyc_output

# Grunt intermediate storage (https://gruntjs.com/creating-plugins#storing-task-files)
.grunt

# Bower dependency directory (https://bower.io/)
bower_components

# node-waf configuration
.lock-wscript

# Compiled binary addons (https://nodejs.org/api/addons.html)
build/Release

# Dependency directories
node_modules/
jspm_packages/

# TypeScript v1 declaration files
typings/

# TypeScript cache
*.tsbuildinfo

# Optional npm cache directory
.npm

# Optional eslint cache
.eslintcache

# Microbundle cache
.rpt2_cache/
.rts2_cache_cjs/
.rts2_cache_es/
.rts2_cache_umd/

# Optional REPL history
.node_repl_history

# Output of 'npm pack'
*.tgz

# Yarn Integrity file
.yarn-integrity

# dotenv environment variables file
.env
.env.test

# parcel-bundler cache (https://parceljs.org/)
.cache
.parcel-cache

# Next.js build output
.next
out

# Nuxt.js build / generate output
.nuxt
dist

# Gatsby files
.cache/
# Comment in the public line in if your project uses Gatsby and not Next.js
# public

# vuepress build output
.vuepress/dist

# Serverless directories
.serverless/

# FuseBox cache
.fusebox/

# DynamoDB Local files
.dynamodb/

# TernJS port file
.tern-port

# Stores VSCode versions used for testing VSCode extensions
.vscode-test

# yarn v2
.yarn/cache
.yarn/unplugged
.yarn/build-state.yml
.yarn/install-state.gz
.pnp.*

# Wodge specific
wodge
tmp/
`

const EnvFile = `# Wodge Environment Variables

# Backend Configuration
PORT=8080

# Add service configurations below via 'wodge add api ...'
`
