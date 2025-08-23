package event

import (
	"encoding/json"
	"fmt"

	"github.com/mattsp1290/ag-ui/go-sdk/pkg/core/events"
)

func Parse(data []byte) (events.Event, error) {
	// Parse the SSE event
	var eventData map[string]interface{}

	err := json.Unmarshal(data, &eventData)
	if err != nil {
		return nil, fmt.Errorf("received non-JSON frame event data %w", err)
	}

	decoder := events.NewEventDecoder(nil)

	// Extract event type - the server sends it as "type" field directly
	eventType, _ := eventData["type"].(string)

	event, err := decoder.DecodeEvent(eventType, data)
	if err != nil {
		return nil, fmt.Errorf("failed to decode event %w", err)
	}

	return event, nil
}
