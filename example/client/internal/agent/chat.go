package agent

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/ag-ui-protocol/ag-ui/sdks/community/go/pkg/client/sse"
	"github.com/ag-ui-protocol/ag-ui/sdks/community/go/pkg/core/events"
	"github.com/mattsp1290/october-talks-2025/example/client/internal/event"
	"github.com/sirupsen/logrus"
)

func DefaultEndpoint() string {
	return "http://localhost:8000/agentic"
}

func Chat(ctx context.Context, inputMsg string, endpoint string, emit func(event events.Event)) error {
	logger := logrus.New()
	logger.SetLevel(logrus.FatalLevel)
	sseConfig := sse.Config{
		Endpoint:       endpoint,
		ConnectTimeout: 30 * time.Second,
		ReadTimeout:    5 * time.Minute,
		BufferSize:     100,
		Logger:         logger,
		AuthHeader:     "Authorization",
		AuthScheme:     "Bearer",
	}

	client := sse.NewClient(sseConfig)
	defer func() {
		client.Close()
	}()

	sessionID := "test-session-1755371887"
	runID := "run-1755744865857245000"

	payload := map[string]interface{}{
		"threadId": sessionID,
		"runId":    runID,
		"state":    map[string]interface{}{},
		"messages": []map[string]interface{}{
			{
				"id":      "msg-1",
				"role":    "user",
				"content": inputMsg,
			},
		},
		"tools":          []interface{}{},
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
		log.Fatal(err)
		//return errors.New("Failed to establish SSE connection")
	}

	// Parse SSE events
	for {
		select {
		case <-ctx.Done():
			break
		case frame, ok := <-frames:
			if !ok {
				return nil
			}

			rawEvent, err := event.Parse(frame.Data)
			if err != nil {
				return fmt.Errorf("failed to process SSE event %w", err)
			}
			emit(rawEvent)

		case err, ok := <-errorCh:
			if !ok {
				break
			}
			if err != nil {
				break
			}
		}
	}
}
