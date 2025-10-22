package agentic

import (
	"bufio"
	"context"
	_ "embed"
	"fmt"

	"github.com/ag-ui-protocol/ag-ui/sdks/community/go/pkg/core/events"
	"github.com/firebase/genkit/go/ai"
	"github.com/firebase/genkit/go/genkit"
	"github.com/firebase/genkit/go/plugins/googlegenai"
	"golang.org/x/sync/errgroup"

	aguigenkit "github.com/ag-ui-protocol/ag-ui/integrations/community/genkit/go/genkit"
	"github.com/ag-ui-protocol/ag-ui/sdks/community/go/pkg/encoding/sse"
)

// reminder is a reminder for the AI to output in our expected format.
//
//go:embed data/reminder.md
var reminder string

//func CallLLM(ctx context.Context, input string, tools []langchaingoTools.Tool, returnChan chan<- string) error {
//
//	adapter, err := internalmcp.NewAdapter(fmt.Sprintf("http://127.0.0.1:%d/mcp", internalmcp.DefaultPort))
//
//	if err != nil {
//		return fmt.Errorf("new mcp adapter: %w", err)
//	}
//
//	_, err = adapter.Tools()
//	if err != nil {
//		return fmt.Errorf("append tools: %w", err)
//	}
//
//	llm, err := tmcanthropic.New(tmcanthropic.WithModel("claude-3-haiku-20240307"))
//	if err != nil {
//		return fmt.Errorf("failed to create LLM client: %w", err)
//	}
//
//	agent := agents.NewOneShotAgent(llm,
//		tools,
//		agents.WithMaxIterations(50))
//
//	executor := agents.NewExecutor(agent, agents.WithCallbacksHandler(NewHandler(returnChan)))
//
//	inputMap := make(map[string]any)
//	inputMap["input"] = input + "\n" + reminder
//
//	result, err := chains.Call(ctx, executor, inputMap)
//	if err != nil {
//		return fmt.Errorf("run chain: %w", err)
//	}
//	output := result["output"].(string)
//
//	// Create a proper event for the final output
//	messageID := events.GenerateMessageID()
//	finalMessage := events.NewTextMessageContentEvent(messageID, output)
//	if jsonData, err := finalMessage.ToJSON(); err == nil {
//		returnChan <- string(jsonData)
//	}
//
//	return nil
//}

func ProcessInput(ctx context.Context, w *bufio.Writer, sseWriter *sse.SSEWriter, input string) error {
	resultChan := make(chan events.Event)
	g, groupCtx := errgroup.WithContext(ctx)

	g.Go(func() error {
		for {
			select {
			case result := <-resultChan:
				if result == nil {
					return nil
				}

				bytes, err := result.ToJSON()
				if err != nil {
					return fmt.Errorf("failed to convert event to JSON: %w", err)
				}
				if err = sseWriter.WriteBytes(ctx, w, bytes); err != nil {
					return fmt.Errorf("failed to write event: %w", err)
				}
			case <-ctx.Done():
				return ctx.Err()
			}
		}
	})

	g.Go(func() error {
		callLLMErr := CallLLM(groupCtx, input, resultChan)
		close(resultChan)
		return callLLMErr
	})

	return g.Wait()
}

func CallLLM(ctx context.Context, input string, returnChan chan<- events.Event) error {
	g := genkit.Init(
		ctx,
		genkit.WithDefaultModel("googleai/gemini-2.5-flash"),
		genkit.WithPlugins(&googlegenai.GoogleAI{}),
	)

	streamingFunc := aguigenkit.StreamingFunc("", "", returnChan)

	_, err := genkit.Generate(ctx, g,
		ai.WithPrompt(input),
		ai.WithStreaming(streamingFunc))
	return err
}
