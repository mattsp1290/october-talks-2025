package agentic

import (
	"bufio"
	"context"
	_ "embed"
	"fmt"
	"log"
	"os/exec"

	"github.com/ag-ui-protocol/ag-ui/sdks/community/go/pkg/core/events"
	"github.com/firebase/genkit/go/ai"
	"github.com/firebase/genkit/go/genkit"
	"github.com/firebase/genkit/go/plugins/googlegenai"
	"github.com/firebase/genkit/go/plugins/mcp"
	mcp_golang "github.com/metoro-io/mcp-golang"
	"github.com/metoro-io/mcp-golang/transport/stdio"
	"golang.org/x/sync/errgroup"

	aguigenkit "github.com/ag-ui-protocol/ag-ui/integrations/community/genkit/go/genkit"
	"github.com/ag-ui-protocol/ag-ui/sdks/community/go/pkg/encoding/sse"
)

const mcpPath = "/Users/punk1290/git/october-talks-2025/example/mcp/mcp"

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

func getPrompt(input string) (string, error) {
	// Start the server process
	cmd := exec.Command(mcpPath)
	stdin, err := cmd.StdinPipe()
	if err != nil {
		log.Fatalf("Failed to get stdin pipe: %v", err)
	}
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatalf("Failed to get stdout pipe: %v", err)
	}

	if err := cmd.Start(); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
	defer cmd.Process.Kill()

	clientTransport := stdio.NewStdioServerTransportWithIO(stdout, stdin)
	client := mcp_golang.NewClient(clientTransport)

	if _, err := client.Initialize(context.Background()); err != nil {
		log.Fatalf("Failed to initialize client: %v", err)
	}

	// List available tools
	prompt, err := client.GetPrompt(context.Background(), "create_pr", map[string]string{"repo": input})
	if err != nil {
		return "", err
	}
	if len(prompt.Messages) == 0 {
		return "", fmt.Errorf("no prompt found")
	}
	return prompt.Messages[0].Content.TextContent.Text, nil
}

func CallLLM(ctx context.Context, input string, returnChan chan<- events.Event) error {
	g := genkit.Init(
		ctx,
		genkit.WithDefaultModel("googleai/gemini-2.5-flash"),
		genkit.WithPlugins(&googlegenai.GoogleAI{}),
	)

	mcpClient, err := mcp.NewGenkitMCPClient(mcp.MCPClientOptions{
		Name: "localMcp",
		Stdio: &mcp.StdioConfig{
			Command: mcpPath,
		},
	})
	if err != nil {
		return err
	}

	mcpTools, err := mcpClient.GetActiveTools(ctx, g)
	if err != nil {
		return err
	}

	tools := make([]ai.ToolRef, len(mcpTools))
	for i, tool := range mcpTools {
		tools[i] = tool
	}

	prompt, err := getPrompt(input)
	if err != nil {
		return err
	}

	streamingFunc := aguigenkit.StreamingFunc("", "", returnChan)

	_, err = genkit.Generate(ctx, g,
		ai.WithPrompt(prompt),
		ai.WithTools(tools...),
		ai.WithMaxTurns(1000),
		ai.WithStreaming(streamingFunc))
	if err != nil {
		returnChan <- events.NewTextMessageContentEvent("", fmt.Sprintf("Error: %v", err))
	} else {
		returnChan <- events.NewTextMessageContentEvent("", "finished processing.")
	}
	returnChan <- events.NewRunFinishedEvent("", "")
	return err
}
