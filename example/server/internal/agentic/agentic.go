package agentic

import (
	"bufio"
	"context"
	_ "embed"
	"fmt"

	"github.com/mattsp1290/october-talks-2025/example/server/internal/mcp"
	"github.com/tmc/langchaingo/agents"
	"github.com/tmc/langchaingo/chains"
	"github.com/tmc/langchaingo/llms/anthropic"

	"github.com/mattsp1290/ag-ui/go-sdk/pkg/core/events"
	"github.com/mattsp1290/ag-ui/go-sdk/pkg/encoding/sse"
	langchaingoTools "github.com/tmc/langchaingo/tools"
)

// reminder is a reminder for the AI to output in our expected format.
//
//go:embed data/reminder.md
var reminder string

func CallLLM(ctx context.Context, input string, tools []langchaingoTools.Tool) (string, error) {

	adapter, err := mcp.NewAdapter(fmt.Sprintf("http://127.0.0.1:%d/mcp", mcp.DefaultPort))

	if err != nil {
		return "", fmt.Errorf("new mcp adapter: %w", err)
	}

	mcpTools, err := adapter.Tools()
	if err != nil {
		return "", fmt.Errorf("append tools: %w", err)
	}

	llm, err := anthropic.New(anthropic.WithModel("claude-3-haiku-20240307"))
	if err != nil {
		return "", fmt.Errorf("failed to create LLM client: %w", err)
	}

	agent := agents.NewOneShotAgent(llm,
		append(tools, mcpTools...),
		agents.WithMaxIterations(50))

	executor := agents.NewExecutor(agent)

	inputMap := make(map[string]any)
	inputMap["input"] = input + "\n" + reminder

	result, err := chains.Call(ctx, executor, inputMap)
	if err != nil {
		return "", fmt.Errorf("run chain: %w", err)
	}
	output := result["output"].(string)
	return output, nil
}

func ProcessInput(ctx context.Context, w *bufio.Writer, sseWriter *sse.SSEWriter, input string) error {
	result, err := CallLLM(ctx, input, nil)
	if err != nil {
		return fmt.Errorf("failed to call LLM: %w", err)
	}
	message := events.NewTextMessageContentEvent("test", result)
	if err := sseWriter.WriteEvent(ctx, w, message); err != nil {
		return fmt.Errorf("failed to write event: %w", err)
	}

	return nil
}
