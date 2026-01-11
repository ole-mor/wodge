package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"wodge/internal/templates"

	"github.com/spf13/cobra"
)

var uiCmd = &cobra.Command{
	Use:   "ui [component]",
	Short: "Add a UI component to your project",
	Long:  `Add a UI component to your project. Available components: button, card, input, navbar, theme-provider, qast-test, secure-chat, token-manager`,
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
	case "secure-chat":
		addSecureChatComponent(appRoot)
	case "login":
		addLoginPage(appRoot)
	default:
		fmt.Printf("Unknown component/page: %s\n", componentName)
		fmt.Println("Available components: button, card, input, navbar, qast-test, secure-chat, login")
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
	fmt.Println("\nIMPORTANT NEXT STEPS:")
	fmt.Println("1. Wrap your App in AuthProvider in `src/App.tsx`:")
	fmt.Println("   import { AuthProvider } from '@/context/AuthProvider';")
	fmt.Println("   function App() { return <AuthProvider>...</AuthProvider> }")
	fmt.Println("\n2. Protect routes in `src/routes.generated.tsx` (or manually):")
	fmt.Println("   import { ProtectedRoute } from '@/components/ProtectedRoute';")
	fmt.Println("   { path: '/protected', element: <ProtectedRoute><YourPage /></ProtectedRoute> }")
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
