package main

// A simple example that shows how to send messages to a Bubble Tea program
// from outside the program using Program.Send(Msg).

import (
	"context"
	"log"

	"github.com/ag-ui-protocol/ag-ui/sdks/community/go/pkg/core/events"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/google/uuid"
	"github.com/mattsp1290/october-talks-2025/example/client/internal/agent"
	"golang.org/x/sync/errgroup"
)

var (
	spinnerStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("63"))
	helpStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("241")).Margin(1, 0)
	msgStyle     = helpStyle.UnsetMargins()
	appStyle     = lipgloss.NewStyle().Margin(1, 2, 0, 2)
)

func main() {
	p := tea.NewProgram(newModel())

	emit := func(event events.Event) {
		p.Send(event)
	}

	eg, groupCtx := errgroup.WithContext(context.Background())
	go func() {
		if _, err := p.Run(); err != nil {
			log.Fatalf("Error running program: %v", err)
		}
	}()

	prompt := "/Users/punk1290/git/october-talks-2025/example/birb-client"

	eg.Go(func() error {
		return agent.Chat(groupCtx, prompt, agent.DefaultEndpoint(), emit)
	})

	err := eg.Wait()
	if err != nil {
		log.Fatal(err)
	}
}

type model struct {
	spinner  spinner.Model
	order    []string
	messages map[string]string
	quitting bool
}

func newModel() *model {
	s := spinner.New()
	s.Style = spinnerStyle
	return &model{
		spinner:  s,
		order:    make([]string, 0),
		messages: make(map[string]string),
	}
}

func (m *model) Init() tea.Cmd {
	return m.spinner.Tick
}

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		m.quitting = true
		return m, tea.Quit
	case events.Event:
		event := events.Event(msg)
		if event.Type() == events.EventTypeRunFinished {
			return m, tea.Quit
		}
		m.processEvent(event)
		return m, nil
	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	default:
		return m, nil
	}
}

func (m *model) View() string {
	var s string

	s += "\n\n"

	for _, currID := range m.order {
		s += m.messages[currID] + "\n"
	}

	if m.quitting {
		s += "Shutting down..."
	} else {
		s += m.spinner.View()
	}

	if !m.quitting {
		s += helpStyle.Render("Press any key to exit")
	}

	if m.quitting {
		s += "\n"
	}

	return appStyle.Render(s)
}

func (m *model) processEvent(event events.Event) {
	eventType := event.Type()
	switch eventType {
	case events.EventTypeRunStarted:
		_, ok := event.(*events.RunStartedEvent)
		if !ok {
			log.Fatal("RunStarted event not valid\n")
		}
		currID := uuid.New().String()
		m.order = append(m.order, currID)
		m.messages[currID] = "Run started"
	case events.EventTypeTextMessageStart:
		_, ok := event.(*events.TextMessageStartEvent)
		if !ok {
			log.Fatal("TextMessageStart event not valid\n")
		}
		currID := uuid.New().String()
		m.order = append(m.order, currID)
		m.messages[currID] = ""
	case events.EventTypeTextMessageChunk:
		evt, ok := event.(*events.TextMessageChunkEvent)
		if !ok {
			log.Fatal("TextMessageChunkEvent event not valid\n")
		}
		if evt.Delta == nil {
			log.Print("Got nil delta")
		}
		content := *evt.Delta

		currID := m.order[len(m.order)-1]
		m.messages[currID] = m.messages[currID] + content
	case events.EventTypeToolCallResult:
		res, ok := event.(*events.ToolCallResultEvent)
		if !ok {
			log.Fatal("ToolCallResult event not valid\n")
		}
		currID := uuid.New().String()
		m.order = append(m.order, currID)
		m.messages[currID] = res.Content
	case events.EventTypeTextMessageContent:
		content, ok := event.(*events.TextMessageContentEvent)
		if !ok {
			log.Fatal("TextMessageContent event not valid\n")
		}
		currID := uuid.New().String()
		m.order = append(m.order, currID)
		m.messages[currID] = content.Delta
	case events.EventTypeTextMessageEnd:
		_, ok := event.(*events.TextMessageEndEvent)
		if !ok {
		}
	default:
		// For any unknown event types, hard fail
		log.Fatalf("Unhandled event type: %s\n", eventType)
	}
}
