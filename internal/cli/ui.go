package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"wodge/internal/templates"

	"github.com/spf13/cobra"
)

var uiCmd = &cobra.Command{
	Use:   "ui [component]",
	Short: "Add a UI component to your project",
	Long:  `Add a UI component to your project. Available components: button, card, input, navbar, theme-provider, qast-test, secure-chat (llm-chat), token-manager`,
	Args:  cobra.ExactArgs(1),
	Run:   addUIComponent,
}

func addUIComponent(cmd *cobra.Command, args []string) {
	componentName := args[0]
	appRoot, err := findAppRoot() // Reuse helper from add.go if possible, or duplicate safely
	if err != nil {
		// Fallback if findAppRoot is not exported or available in this package scope (it is in same package)
		appRoot, _ = os.Getwd()
	}

	switch componentName {
	case "button":
		addComponentFile(appRoot, "button", "src/components/ui/Button.tsx", templates.ComponentButton)
	case "card":
		addComponentFile(appRoot, "card", "src/components/ui/Card.tsx", templates.ComponentCard)
	case "input":
		addComponentFile(appRoot, "input", "src/components/ui/Input.tsx", templates.ComponentInput)
	case "navbar":
		addComponentFile(appRoot, "navbar", "src/components/ui/Navbar.tsx", templates.ComponentNavbar)
	case "theme-provider":
		addComponentFile(appRoot, "theme-provider", "src/context/ThemeProvider.tsx", templates.ThemeProvider)
	case "qast-test":
		addComponentFile(appRoot, "qast-test", "src/components/ui/QastTest.tsx", templates.ComponentQastTest)
	case "token-manager":
		addComponentFile(appRoot, "token-manager", "src/utils/TokenManager.ts", templates.ComponentTokenManager)
	case "sidebar":
		addComponentFile(appRoot, "sidebar", "src/components/ui/Sidebar.tsx", templates.ComponentSidebar)
		// Sidebar needs API clients
		addComponentFile(appRoot, "users-api", "src/api/users.ts", templates.ComponentUsersAPI)
		addComponentFile(appRoot, "history-api", "src/api/history.ts", templates.ComponentHistoryAPI)
	case "secure-chat", "llm-chat":
		addSecureChatComponent(appRoot)
	case "llmwrapper":
		addLLMWrapperComponent(appRoot)
	case "login":
		addLoginPage(appRoot)
	default:
		fmt.Printf("Unknown component/page: %s\n", componentName)
		fmt.Println("Available components: button, card, input, navbar, qast-test, secure-chat (llm-chat), login, llmwrapper")
		os.Exit(1)
	}
}

func addLoginPage(appRoot string) {
	fmt.Println("Adding Login Page and Auth Context...")

	// 1. Ensure UI components exist (Button, Input, Card)
	addComponentFile(appRoot, "button", "src/components/ui/Button.tsx", templates.ComponentButton)
	addComponentFile(appRoot, "input", "src/components/ui/Input.tsx", templates.ComponentInput)
	addComponentFile(appRoot, "card", "src/components/ui/Card.tsx", templates.ComponentCard)

	files := map[string]string{
		"src/context/AuthProvider.tsx":      templates.ComponentAuthProvider,
		"src/components/ProtectedRoute.tsx": templates.ComponentProtectedRoute,
		"src/routes/login.route.tsx":        templates.ComponentLoginPage,
	}

	for path, content := range files {
		fullPath := filepath.Join(appRoot, path)
		dir := filepath.Dir(fullPath)
		_ = os.MkdirAll(dir, 0755)
		if err := os.WriteFile(fullPath, []byte(content), 0644); err != nil {
			fmt.Printf("Error writing %s: %v\n", path, err)
		} else {
			fmt.Printf("Created %s\n", path)
		}
	}

	fmt.Println("\nLogin Page added successfully!")

	// 2. Wrap App.tsx with AuthProvider
	appTsxPath := filepath.Join(appRoot, "src", "App.tsx")
	if content, err := os.ReadFile(appTsxPath); err == nil {
		sContent := string(content)
		if !strings.Contains(sContent, "AuthProvider") {
			// Add import
			if !strings.Contains(sContent, "@/context/AuthProvider") {
				sContent = "import { AuthProvider } from '@/context/AuthProvider';\n" + sContent
			}
			// Wrap <GeneratedRoutes />
			// Simple heuristic: find <GeneratedRoutes /> and wrap it?
			// Or wrap the entire content returned by App?
			// The user requirement: "wrap the content of generated routes ... with the authprovider wrapper"
			// Actually, wrapping GeneratedRoutes is safest to ensure context is available to all routes.

			if strings.Contains(sContent, "<GeneratedRoutes />") {
				sContent = strings.Replace(sContent, "<GeneratedRoutes />", "<AuthProvider><GeneratedRoutes /></AuthProvider>", 1)
				if err := os.WriteFile(appTsxPath, []byte(sContent), 0644); err != nil {
					fmt.Printf("Warning: Failed to update App.tsx automatically: %v\n", err)
				} else {
					fmt.Println("✓ Automatically wrapped generated routes with AuthProvider in App.tsx")
				}
			} else {
				fmt.Println("Note: Could not find <GeneratedRoutes /> in App.tsx. Please manually wrap your routes with <AuthProvider>.")
			}
		}
	}

	fmt.Println("\nIMPORTANT NEXT STEPS:")
	fmt.Println("1. Ensure your backend AstAuth server is running.")
	fmt.Println("2. Run `wodge run dev` - your routes (except login) should now be protected automatically!")
}

func addComponentFile(appRoot, component, path, content string) {
	fullPath := filepath.Join(appRoot, path)
	dir := filepath.Dir(fullPath)

	if err := os.MkdirAll(dir, 0755); err != nil {
		fmt.Printf("Error creating directory %s: %v\n", dir, err)
		return
	}

	if _, err := os.Stat(fullPath); err == nil {
		fmt.Printf("Component '%s' already exists at %s\n", component, path)
		return
	}

	if err := os.WriteFile(fullPath, []byte(content), 0644); err != nil {
		fmt.Printf("Error writing file %s: %v\n", fullPath, err)
		return
	}

	fmt.Printf("✓ Added %s to %s\n", component, path)
}

func addSecureChatComponent(appRoot string) {
	fmt.Println("Adding SecureChat Component and Dependencies...")
	// Dependencies: Card, Input, Button, TokenManager
	addComponentFile(appRoot, "card", "src/components/ui/Card.tsx", templates.ComponentCard)
	addComponentFile(appRoot, "input", "src/components/ui/Input.tsx", templates.ComponentInput)
	addComponentFile(appRoot, "button", "src/components/ui/Button.tsx", templates.ComponentButton)
	addComponentFile(appRoot, "token-manager", "src/utils/TokenManager.ts", templates.ComponentTokenManager)

	// Add the Chat component
	addComponentFile(appRoot, "secure-chat", "src/components/ui/SecureChat.tsx", templates.ComponentSecureChat)
}

func addLLMWrapperComponent(appRoot string) {
	fmt.Println("Adding LLMWrapper Component and Dependencies...")

	// 1. Ensure Dependencies exist
	addComponentFile(appRoot, "button", "src/components/ui/Button.tsx", templates.ComponentButton)
	addComponentFile(appRoot, "sidebar", "src/components/ui/Sidebar.tsx", templates.ComponentSidebar)

	// LLMWrapper relies on AuthContext, ensure it is set up or assume it exists if user adds this?
	// Best practice: Ensure generic UI components exist, but AuthContext is usually added via 'wodge add ui login'.
	// We will warn if missing, or just add the wrapper file.
	// It basically needs SecureChat.
	addSecureChatComponent(appRoot)

	// 2. Add the Layout
	addComponentFile(appRoot, "llm-layout", "src/components/layout/LLMLayout.tsx", templates.ComponentLLMWrapper)

	fmt.Println("\nLLMLayout added successfully!")
	fmt.Println("Usage: Import LLMLayout and use it as a wrapper for your routes or as a standalone page.")
}
