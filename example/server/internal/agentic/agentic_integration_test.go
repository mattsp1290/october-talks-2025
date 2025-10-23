package agentic

import (
	"context"
	_ "embed"
	"testing"

	"github.com/ag-ui-protocol/ag-ui/sdks/community/go/pkg/core/events"
	"github.com/stretchr/testify/require"
	"golang.org/x/sync/errgroup"
)

func TestToolCalls(t *testing.T) {
	ctx := context.Background()
	resultChan := make(chan events.Event)
	g, groupCtx := errgroup.WithContext(ctx)
	var results []events.Event

	g.Go(func() error {
		for {
			select {
			case result := <-resultChan:
				if result == nil {
					return nil
				}
				results = append(results, result)
			case <-groupCtx.Done():
				return groupCtx.Err()
			case <-ctx.Done():
				return ctx.Err()
			}
		}
	})

	g.Go(func() error {
		callErr := CallLLM(groupCtx, "/Users/punk1290/git/october-talks-2025/example/birb-client", resultChan)
		close(resultChan)
		if callErr != nil {
			return callErr
		}
		return nil
	})
	
	err := g.Wait()
	require.NoError(t, err)
	require.NotEmpty(t, results)
}
