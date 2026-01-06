package cli

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"
	"wodge/internal/monitor"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
)

var monitorCmd = &cobra.Command{
	Use:   "monitor",
	Short: "Monitor the running Wodge backend",
	Run: func(cmd *cobra.Command, args []string) {
		p := tea.NewProgram(initialModel())
		if _, err := p.Run(); err != nil {
			fmt.Printf("Alas, there's been an error: %v", err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(monitorCmd)
}

// -- Bubble Tea Model --

type model struct {
	events []monitor.Event
	table  table.Model
	err    error
}

type eventMsg monitor.Event
type errMsg error

func initialModel() model {
	columns := []table.Column{
		{Title: "Time", Width: 10},
		{Title: "Type", Width: 10},
		{Title: "Details", Width: 50},
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
		events: []monitor.Event{},
		table:  t,
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
		return m, listenForEvents // Continue listening (actually this should be a subscription loop)
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
		payloadStr := fmt.Sprintf("%v", e.Payload)
		// Truncate payload for display
		if len(payloadStr) > 47 {
			payloadStr = payloadStr[:47] + "..."
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
	return baseStyle.Render(m.table.View()) + "\n  Press 'q' to quit.\n"
}

var baseStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("240"))

// -- Event Listener --

func listenForEvents() tea.Msg {
	// This is a simplified "poll" for the sake of the TUI command loop.
	// In a real TUI, we'd start a separate goroutine that pumps msgs to the program.
	// But since Init() expects a single Cmd that returns a single Msg, we need a slight adapter.
	// For this proof of concept, we'll try to connect once and block until an event comes.

	// Check connection existence (simple hack for now: better struct needed)
	if eventChan == nil {
		go startEventStream()
		time.Sleep(100 * time.Millisecond) // Give it a sec to connect
	}

	val := <-eventChan
	return eventMsg(val)
}

var eventChan = make(chan monitor.Event)

func startEventStream() {
	resp, err := http.Get("http://localhost:8080/wodge/monitor/events")
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
