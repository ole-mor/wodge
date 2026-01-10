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
	Long:  `Add a UI component to your project. Available components: button, card, input, navbar, theme-provider`,
	Args:  cobra.ExactArgs(1),
	Run:   addUIComponent,
}

func addUIComponent(cmd *cobra.Command, args []string) {
	component := args[0]
	appRoot, _ := os.Getwd()

	componentMap := map[string]struct {
		path    string
		content string
	}{
		"button":         {"src/components/ui/Button.tsx", templates.ComponentButton},
		"card":           {"src/components/ui/Card.tsx", templates.ComponentCard},
		"input":          {"src/components/ui/Input.tsx", templates.ComponentInput},
		"navbar":         {"src/components/ui/Navbar.tsx", templates.ComponentNavbar},
		"theme-provider": {"src/context/ThemeProvider.tsx", templates.ThemeProvider},
		"qast-test":      {"src/components/ui/QastTest.tsx", templates.ComponentQastTest},
	}

	comp, exists := componentMap[component]
	if !exists {
		fmt.Printf("Unknown component '%s'. Available: button, card, input, navbar, theme-provider\n", component)
		return
	}

	fullPath := filepath.Join(appRoot, comp.path)
	dir := filepath.Dir(fullPath)

	// Create directory if it doesn't exist
	if err := os.MkdirAll(dir, 0755); err != nil {
		fmt.Printf("Error creating directory %s: %v\n", dir, err)
		return
	}

	// Check if file already exists
	if _, err := os.Stat(fullPath); err == nil {
		fmt.Printf("Component '%s' already exists at %s\n", component, comp.path)
		return
	}

	// Write component file
	if err := os.WriteFile(fullPath, []byte(comp.content), 0644); err != nil {
		fmt.Printf("Error writing file %s: %v\n", fullPath, err)
		return
	}

	fmt.Printf("✓ Added %s to %s\n", component, comp.path)
}
