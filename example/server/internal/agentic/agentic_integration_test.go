package agentic

import (
	"context"
	_ "embed"
	"os"
	"testing"

	"github.com/mattsp1290/october-talks-2025/example/server/internal/mcp"
	"github.com/stretchr/testify/require"
)

//go:embed data/client_prompt.md
var client_prompt string

//go:embed data/languages_prompt.md
var languages_prompt string

func TestToolCalls(t *testing.T) {
	if os.Getenv("RUN_INTEGRATION_TESTS") != "true" {
		t.Skip("Skipping integration test")
	}
	mcpServer, err := mcp.NewServer(mcp.DefaultPort)
	require.NoError(t, err)
	go func() {
		mcpErr := mcpServer.Start()
		if mcpErr != nil {
			require.NoError(t, mcpErr)
		}
	}()

	ctx := context.Background()
	results, err := CallLLM(ctx, languages_prompt, nil)
	require.NoError(t, err)
	require.NotEmpty(t, results)
}
