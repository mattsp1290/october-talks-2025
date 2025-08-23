package main

// A simple program demonstrating the text area component from the Bubbles
// component library.

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/mattsp1290/ag-ui/go-sdk/pkg/client/sse"
	"github.com/mattsp1290/october-talks-2025/example/client/internal/event"
	"github.com/mattsp1290/october-talks-2025/example/client/internal/message"
	"github.com/sirupsen/logrus"
)

const gap = "\n\n"

type (
	errMsg error
)

type model struct {
	viewport    viewport.Model
	messages    []string
	textarea    textarea.Model
	senderStyle lipgloss.Style
	userInput   chan string
	err         error
}

func initialModel(userInput chan string) model {
	ta := textarea.New()
	ta.Placeholder = "Send a message..."
	ta.Focus()

	ta.Prompt = "┃ "
	ta.CharLimit = 280

	ta.SetWidth(30)
	ta.SetHeight(3)

	// Remove cursor line styling
	ta.FocusedStyle.CursorLine = lipgloss.NewStyle()

	ta.ShowLineNumbers = false

	vp := viewport.New(30, 5)
	vp.SetContent(`Welcome to the chat room!
Type a message and press Enter to send.`)

	ta.KeyMap.InsertNewline.SetEnabled(false)

	return model{
		textarea:    ta,
		messages:    []string{},
		viewport:    vp,
		userInput:   userInput,
		senderStyle: lipgloss.NewStyle().Foreground(lipgloss.Color("5")),
		err:         nil,
	}
}

func (m model) Init() tea.Cmd {
	return textarea.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		tiCmd tea.Cmd
		vpCmd tea.Cmd
	)

	m.textarea, tiCmd = m.textarea.Update(msg)
	m.viewport, vpCmd = m.viewport.Update(msg)

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.viewport.Width = msg.Width
		m.textarea.SetWidth(msg.Width)
		m.viewport.Height = msg.Height - m.textarea.Height() - lipgloss.Height(gap)

	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			fmt.Println(m.textarea.Value())
			return m, tea.Quit
		case tea.KeyEnter:
			m.userInput <- m.textarea.Value()
			m.messages = append(m.messages, m.senderStyle.Render("You: ")+m.textarea.Value())
			//m.viewport.SetContent(lipgloss.NewStyle().Width(m.viewport.Width).Render(strings.Join(m.messages, "\n")))
			m.textarea.Reset()
			//m.viewport.GotoBottom()
		}
	case *message.Message:
		m.messages = append(m.messages, msg.Strings()...)

	// We handle errors just like any other message
	case errMsg:
		m.err = msg
		return m, nil
	}

	if len(m.messages) > 0 {
		// Wrap content before setting it.
		m.viewport.SetContent(lipgloss.NewStyle().Width(m.viewport.Width).Render(strings.Join(m.messages, "\n")))
		m.viewport.GotoBottom()
	}

	return m, tea.Batch(tiCmd, vpCmd)
}

func (m model) View() string {
	return fmt.Sprintf(
		"%s%s%s",
		m.viewport.View(),
		gap,
		m.textarea.View(),
	)
}

func chat(ctx context.Context, msg string, p *tea.Program) error {
	logger := logrus.New()
	logger.SetLevel(logrus.FatalLevel)
	endpoint := "http://localhost:8000/agentic"
	sseConfig := sse.Config{
		Endpoint:       endpoint,
		ConnectTimeout: 30 * time.Second,
		ReadTimeout:    5 * time.Minute,
		BufferSize:     100,
		Logger:         nil,
		AuthHeader:     "Authorization",
		AuthScheme:     "Bearer",
	}

	client := sse.NewClient(sseConfig)
	defer func() {
		client.Close()
	}()

	sessionID := "test-session-1755371887"
	runID := "run-1755744865857245000"

	//logger.WithFields(logrus.Fields{
	//	"endpoint":   endpoint,
	//	"session_id": sessionID,
	//	"run_id":     runID,
	//}).Debug("Connecting to SSE stream")

	//payload := sse.RunAgentInput{
	//	SessionID: sessionID,
	//}

	payload := map[string]interface{}{
		"threadId": sessionID,
		"runId":    runID,
		"state":    map[string]interface{}{},
		"messages": []map[string]interface{}{
			{
				"id":      "msg-1",
				"role":    "user",
				"content": msg,
			},
		},
		"tools":          []interface{}{}, // Request tool discovery
		"context":        []interface{}{},
		"forwardedProps": map[string]interface{}{},
	}

	// Start the SSE stream
	var err error
	frames, errorCh, err := client.Stream(sse.StreamOptions{
		Context: ctx,
		Payload: payload,
	})

	if err != nil {
		//logger.WithError(err).Error("Failed to establish SSE connection")
		return errors.New("Failed to establish SSE connection")
	}

	// Parse SSE events
	for {
		select {
		case frame, ok := <-frames:
			if !ok {
				return nil
			}

			rawEvent, err := event.Parse(frame.Data)
			if err != nil {
				//logger.WithError(err).Error("Failed to process SSE event")
				return fmt.Errorf("failed to process SSE event %w", err)
			}
			currMsg := message.NewMessage(rawEvent)
			if currMsg == nil {
				return fmt.Errorf("failed to parse message %w", err)
			}
			p.Send(currMsg)

		case err, ok := <-errorCh:
			if !ok {
				break
			}
			if err != nil {
				//logger.WithError(err).Error("SSE stream error")
				break
			}

		case <-ctx.Done():
			//logger.Debug("Context cancelled, closing stream")
			break
		}
	}

	return nil
}

func runTea(p *tea.Program, userInputCh chan string) error {
	defer close(userInputCh)
	_, err := p.Run()
	return err
}

func main() {
	userInputCh := make(chan string)
	p := tea.NewProgram(initialModel(userInputCh), tea.WithAltScreen())

	go func() {
		for msg := range userInputCh {
			err := chat(context.Background(), msg, p)
			if err != nil {
				log.Fatal(err)
			}
		}
	}()

	teaErr := runTea(p, userInputCh)
	if teaErr != nil {
		log.Fatal(teaErr)
	}
}
