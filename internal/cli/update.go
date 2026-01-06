package cli

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/spf13/cobra"
)

// TODO: Replace with your actual GitHub repository URL
const repoURL = "https://github.com/ole-mor/wodge.git"

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update wodge to the latest version from GitHub",
	Run:   runUpdate,
}

func runUpdate(cmd *cobra.Command, args []string) {
	fmt.Println("Updating wodge...")

	// 1. Create a temporary directory
	tempDir, err := os.MkdirTemp("", "wodge-update-*")
	if err != nil {
		fmt.Printf("Error creating temp directory: %v\n", err)
		os.Exit(1)
	}
	defer os.RemoveAll(tempDir) // Clean up

	// 2. Clone the repository
	fmt.Printf("Cloning %s...\n", repoURL)
	gitCmd := exec.Command("git", "clone", repoURL, tempDir)
	gitCmd.Stdout = os.Stdout
	gitCmd.Stderr = os.Stderr
	if err := gitCmd.Run(); err != nil {
		fmt.Printf("Error cloning repository: %v\n", err)
		fmt.Println("Make sure you have git installed and access to the repository.")
		os.Exit(1)
	}

	// 3. Build the binary
	fmt.Println("Building wodge...")
	buildCmd := exec.Command("go", "build", "-o", "wodge", "./cmd/wodge")
	buildCmd.Dir = tempDir
	buildCmd.Stdout = os.Stdout
	buildCmd.Stderr = os.Stderr
	if err := buildCmd.Run(); err != nil {
		fmt.Printf("Error building wodge: %v\n", err)
		os.Exit(1)
	}

	// 4. Find current executable path
	exePath, err := os.Executable()
	if err != nil {
		fmt.Printf("Error determining executable path: %v\n", err)
		os.Exit(1)
	}

	// Resolve symlinks if any (e.g. if installed via homebrew or checking actual binary)
	realPath, err := filepath.EvalSymlinks(exePath)
	if err != nil {
		fmt.Printf("Warning: Could not resolve symlinks for %s: %v\n", exePath, err)
		realPath = exePath
	}

	fmt.Printf("Replacing binary at %s...\n", realPath)

	// 5. Replace the binary
	// We might need sudo if it's in a protected directory like /usr/local/bin
	// If we can't write, we try with sudo
	srcBinary := filepath.Join(tempDir, "wodge")

	// Try direct move first
	err = os.Rename(srcBinary, realPath)
	if err != nil {
		// If permission denied, try with sudo
		if os.IsPermission(err) {
			fmt.Println("Permission denied. Trying with sudo...")
			sudoCmd := exec.Command("sudo", "mv", srcBinary, realPath)
			sudoCmd.Stdout = os.Stdout
			sudoCmd.Stderr = os.Stderr
			sudoCmd.Stdin = os.Stdin // Allow sudo prompt
			if err := sudoCmd.Run(); err != nil {
				fmt.Printf("Error replacing binary with sudo: %v\n", err)
				os.Exit(1)
			}
		} else {
			// Cross-device rename might fail (e.g. /tmp to /usr/local/bin), so we copy
			// or if it's strictly a rename error not permission related
			// Let's force move via shell command which handles cross-device
			mvCmd := exec.Command("mv", srcBinary, realPath)
			if err := mvCmd.Run(); err != nil {
				// Retry with sudo if ordinary move fails
				fmt.Println("Move failed. Trying with sudo...")
				sudoCmd := exec.Command("sudo", "mv", srcBinary, realPath)
				sudoCmd.Stdout = os.Stdout
				sudoCmd.Stderr = os.Stderr
				sudoCmd.Stdin = os.Stdin
				if err := sudoCmd.Run(); err != nil {
					fmt.Printf("Error replacing binary: %v\n", err)
					os.Exit(1)
				}
			}
		}
	}

	fmt.Println("Successfully updated wodge!")
}
