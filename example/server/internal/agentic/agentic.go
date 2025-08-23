package agentic

import (
	"bufio"
	"context"
	"fmt"

	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/anthropic"

	"github.com/mattsp1290/ag-ui/go-sdk/pkg/core/events"
	"github.com/mattsp1290/ag-ui/go-sdk/pkg/encoding/sse"
)

func ProcessInput(ctx context.Context, w *bufio.Writer, sseWriter *sse.SSEWriter, input string) error {
	llm, err := anthropic.New(anthropic.WithModel("claude-3-haiku-20240307"))
	if err != nil {
		return fmt.Errorf("failed to create LLM client: %w", err)
	}

	// Sending initial message to the model, with a list of available tools.
	messageHistory := []llms.MessageContent{
		llms.TextParts(llms.ChatMessageTypeHuman, input),
	}

	fmt.Println("Querying for weather in Boston..")
	resp, err := llm.GenerateContent(ctx, messageHistory)
	if err != nil {
		return fmt.Errorf("failed to generate content: %w", err)
	}

	for _, choice := range resp.Choices {
		message := events.NewTextMessageContentEvent("test", choice.Content)
		if err := sseWriter.WriteEvent(ctx, w, message); err != nil {
			return fmt.Errorf("failed to write event: %w", err)
		}
	}

	return nil
}
