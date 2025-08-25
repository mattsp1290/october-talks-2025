package main

// A simple program demonstrating the text area component from the Bubbles
// component library.

import (
	"context"
	"log"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/mattsp1290/october-talks-2025/example/client/internal/agent"
	"github.com/mattsp1290/october-talks-2025/example/client/internal/ui"
)

func runTea(p *tea.Program, userInputCh chan string) error {
	defer close(userInputCh)
	_, err := p.Run()
	return err
}

func main() {
	userInputCh := make(chan string)
	p := tea.NewProgram(ui.InitialModel(userInputCh), tea.WithAltScreen())

	go func() {
		for msg := range userInputCh {
			err := agent.Chat(context.Background(), msg, p)
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
