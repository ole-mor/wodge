package cli

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"syscall"
	"time"
	"wodge/internal/monitor"
	"wodge/internal/registry"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
)

var monitorCmd = &cobra.Command{
	Use:   "monitor [app_name]",
	Short: "Monitor a running Wodge backend",
	Long: `Monitor a running Wodge backend. 
If no app name is specified, it will try to find one in the current directory or list available apps.

Available subcommands:
  monitor list           List all registered apps
  monitor [app]          Connect to the app's event stream
  monitor stop [app]     Stop the running app instance
  monitor remove [app]   Remove form registry (stops first if running)`,
	Args: cobra.RangeArgs(0, 2),
	Run: func(cmd *cobra.Command, args []string) {
		reg, err := registry.Load()
		if err != nil {
			fmt.Printf("Error loading registry: %v\n", err)
			os.Exit(1)
		}

		// Handle subcommands
		if len(args) > 0 {
			command := args[0]

			if command == "list" {
				listApps(reg)
				return
			}

			if len(args) == 2 {
				appName := args[1]
				app, exists := reg.Apps[appName]
				if !exists {
					fmt.Printf("App '%s' not found.\n", appName)
					os.Exit(1)
				}

				if command == "stop" {
					if app.Status != "running" {
						fmt.Printf("App '%s' is already stopped.\n", appName)
						return
					}
					stopApp(app)
					return
				}

				if command == "remove" {
					// Stop if running
					if app.Status == "running" {
						stopApp(app)
					}
					reg.Remove(appName)
					fmt.Printf("App '%s' removed from registry.\n", appName)
					return
				}

				// For Start/Restart, we need more logic to execute commands
				// This is tricky as we need to spawn separate processes
				if command == "start" {
					if app.Status == "running" {
						fmt.Printf("App '%s' is already running.\n", appName)
						return
					}
					// TODO: Implement start logic (spawn wodge run in detached mode)
					fmt.Println("Start command not yet fully implemented. Please run 'wodge run' in the app directory.")
					return
				}
			}
		}

		var targetApp registry.WodgeApp

		if len(args) == 1 {
			appName := args[0]
			var ok bool
			targetApp, ok = reg.Apps[appName]
			if !ok {
				fmt.Printf("App '%s' not found running in registry.\n", appName)
				listApps(reg)
				os.Exit(1)
			}
		} else {
			// No args, try to find one or list
			if len(reg.Apps) == 1 {
				for _, app := range reg.Apps {
					targetApp = app
				}
				fmt.Printf("auto-selecting only running app: %s\n", targetApp.Name)
			} else if len(reg.Apps) == 0 {
				fmt.Println("No running Wodge apps found.")
				os.Exit(0)
			} else {
				fmt.Println("Multiple apps running. Please specify one:")
				listApps(reg)
				os.Exit(0)
			}
		}

		fmt.Printf("Connecting to %s on port %d...\n", targetApp.Name, targetApp.Port)

		// Configure global eventChan for valid app
		currentPort = targetApp.Port

		p := tea.NewProgram(initialModel(targetApp.Name))
		if _, err := p.Run(); err != nil {
			fmt.Printf("Alas, there's been an error: %v", err)
			os.Exit(1)
		}
	},
}

func listApps(reg *registry.Registry) {
	fmt.Println("\nWodge Apps Registry:")
	fmt.Printf("%-20s %-10s %-10s %-10s %s\n", "NAME", "STATUS", "PORT", "PID", "PATH")
	fmt.Println(strings.Repeat("-", 80))
	for _, app := range reg.Apps {
		status := app.Status
		if status == "" {
			status = "running" // Backwards compat
		}

		portStr := fmt.Sprintf("%d", app.Port)
		pidStr := fmt.Sprintf("%d", app.PID)

		if status != "running" {
			portStr = "-"
			pidStr = "-"
		}

		fmt.Printf("%-20s %-10s %-10s %-10s %s\n", app.Name, status, portStr, pidStr, app.Path)
	}
	fmt.Println()
}

func init() {
	rootCmd.AddCommand(monitorCmd)
}

// Global to pass to event listener (not elegant but works for this structure)
var currentPort int = 8080

// -- Bubble Tea Model --

type model struct {
	appName string
	events  []monitor.Event
	table   table.Model
	err     error
}

type eventMsg monitor.Event
type errMsg error

func initialModel(appName string) model {
	columns := []table.Column{
		{Title: "Time", Width: 10},
		{Title: "Type", Width: 10},
		{Title: "Details", Width: 60},
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithFocused(true),
		table.WithHeight(20),
	)

	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		Bold(false)
	s.Selected = s.Selected.
		Foreground(lipgloss.Color("229")).
		Background(lipgloss.Color("57")).
		Bold(false)
	t.SetStyles(s)

	return model{
		appName: appName,
		events:  []monitor.Event{},
		table:   t,
	}
}

func (m model) Init() tea.Cmd {
	return listenForEvents
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		}
	case eventMsg:
		m.events = append(m.events, monitor.Event(msg))
		// Keep only last 100 events
		if len(m.events) > 100 {
			m.events = m.events[1:]
		}
		m.updateTable()
		return m, listenForEvents
	case errMsg:
		m.err = msg
		return m, tea.Quit
	}

	m.table, cmd = m.table.Update(msg)
	return m, cmd
}

func (m *model) updateTable() {
	rows := []table.Row{}
	// Show in reverse order (newest top)
	for i := len(m.events) - 1; i >= 0; i-- {
		e := m.events[i]
		var payloadStr string

		if e.Type == "REQUEST" {
			// Try to parse the map
			if data, ok := e.Payload.(map[string]interface{}); ok {
				status := fmt.Sprintf("%v", data["status"])
				method := fmt.Sprintf("%v", data["method"])
				path := fmt.Sprintf("%v", data["path"])
				duration := fmt.Sprintf("%v", data["duration_ms"])

				// Gin-like format: | 200 | GET /path | 123ms |
				payloadStr = fmt.Sprintf("| %s | %s %s | %sms |", status, method, path, duration)
			} else {
				payloadStr = fmt.Sprintf("%v", e.Payload)
			}
		} else {
			payloadStr = fmt.Sprintf("%v", e.Payload)
		}

		// Truncate payload for display
		if len(payloadStr) > 57 {
			payloadStr = payloadStr[:57] + "..."
		}

		rows = append(rows, table.Row{
			e.Timestamp.Format("15:04:05"),
			string(e.Type),
			payloadStr,
		})
	}
	m.table.SetRows(rows)
}

func (m model) View() string {
	if m.err != nil {
		return fmt.Sprintf("Error: %v\n", m.err)
	}

	header := fmt.Sprintf("Monitoring: %s (Port %d)", m.appName, currentPort)
	return baseStyle.Render(header+"\n\n"+m.table.View()) + "\n  Press 'q' to quit.\n"
}

var baseStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("240")).
	Padding(1)

// -- Event Listener --

func listenForEvents() tea.Msg {
	if eventChan == nil {
		eventChan = make(chan monitor.Event)
		go startEventStream()
		time.Sleep(100 * time.Millisecond) // Give it a sec to connect
	}

	val := <-eventChan
	return eventMsg(val)
}

var eventChan chan monitor.Event

func startEventStream() {
	url := fmt.Sprintf("http://localhost:%d/wodge/monitor/events", currentPort)
	resp, err := http.Get(url)
	if err != nil {
		// Send error maybe? for now retry
		time.Sleep(1 * time.Second)
		startEventStream()
		return
	}
	defer resp.Body.Close()

	reader := bufio.NewReader(resp.Body)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			break
		}
		if strings.HasPrefix(line, "data:") {
			jsonStr := strings.TrimPrefix(line, "data:")
			var evt monitor.Event
			if err := json.Unmarshal([]byte(jsonStr), &evt); err == nil {
				eventChan <- evt
			}
		}
	}
}

func stopApp(app registry.WodgeApp) {
	fmt.Printf("Stopping %s (PID %d)...\n", app.Name, app.PID)
	proc, err := os.FindProcess(app.PID)
	if err == nil {
		err = proc.Signal(syscall.SIGINT) // SIGINT allows graceful shutdown we implemented
		if err != nil {
			fmt.Printf("Error stopping process: %v\n", err)
		} else {
			fmt.Println("Stop signal sent.")
		}
	} else {
		fmt.Printf("Could not find process: %v\n", err)
	}
}
