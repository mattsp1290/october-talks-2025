package message

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/charmbracelet/lipgloss"
	"github.com/mattsp1290/ag-ui/go-sdk/pkg/core/events"
)

var serverStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("21"))

type Message struct {
	contents []string
}

func (m *Message) Strings() []string {
	return m.contents
}

func NewMessage(event events.Event) *Message {
	return getMessageFromEvent(event)
}

func getNameplate() string {
	return serverStyle.Render("Assistant: ")
}

func getMessageFromEvent(event events.Event) *Message {
	eventType := event.Type()
	switch eventType {
	case events.EventTypeRunStarted:
		_, ok := event.(*events.RunStartedEvent)
		if !ok {
			return nil
		}
		content := getNameplate() + "Run started"
		return &Message{
			contents: []string{content},
		}

	case events.EventTypeRunFinished:
		_, ok := event.(*events.RunFinishedEvent)
		if !ok {
			return nil
		}
		content := getNameplate() + "Run finished"
		return &Message{
			contents: []string{content},
		}
	case events.EventTypeRunError:
		_, ok := event.(*events.RunErrorEvent)
		if !ok {
			return nil
		}
		log.Fatalf("Event type: %s is not yet supported. \n", eventType)
		return &Message{}
	case events.EventTypeTextMessageStart:
		_, ok := event.(*events.TextMessageStartEvent)
		if !ok {
			return nil
		}
		curMsg := getNameplate() + "text message started"
		return &Message{
			contents: []string{curMsg},
		}
	case events.EventTypeTextMessageContent:
		msg, ok := event.(*events.TextMessageContentEvent)
		if !ok {
			return nil
		}
		return &Message{
			contents: []string{getNameplate() + msg.Delta},
		}
	case events.EventTypeTextMessageEnd:
		_, ok := event.(*events.TextMessageEndEvent)
		if !ok {
			return nil
		}
		curMsg := getNameplate() + "text message ended"
		return &Message{
			contents: []string{curMsg},
		}
	case events.EventTypeToolCallStart:
		_, ok := event.(*events.ToolCallStartEvent)
		if !ok {
			return nil
		}
		curMsg := getNameplate() + "tool call started"
		return &Message{
			contents: []string{curMsg},
		}
	case events.EventTypeToolCallArgs:
		args, ok := event.(*events.ToolCallArgsEvent)
		if !ok {
			return nil
		}
		curMsg := fmt.Sprintf("%stool call args: %s", getNameplate(), args.Delta)
		return &Message{
			contents: []string{curMsg},
		}
	case events.EventTypeToolCallEnd:
		_, ok := event.(*events.ToolCallEndEvent)
		if !ok {
			return nil
		}
		curMsg := getNameplate() + "tool call ended"
		return &Message{
			contents: []string{curMsg},
		}
	case events.EventTypeToolCallResult:
		result, ok := event.(*events.ToolCallResultEvent)
		if !ok {
			return nil
		}
		curMsg := getNameplate() + result.Content
		return &Message{
			contents: []string{curMsg},
		}
	case events.EventTypeStateSnapshot:
		snapshot, ok := event.(*events.StateSnapshotEvent)
		if !ok {
			return nil
		}
		var contents []string
		if snapshot.Snapshot != nil {
			jsonData, err := json.Marshal(snapshot.Snapshot)
			if err != nil {
				fmt.Println("Error marshaling JSON:", err)
				return nil
			}
			contents = append(contents, getNameplate()+string(jsonData))

		}
		return &Message{
			contents: contents,
		}
	case events.EventTypeStateDelta:
		delta, ok := event.(*events.StateDeltaEvent)
		if !ok {
			return nil
		}
		var contents []string
		for _, op := range delta.Delta {
			currOp := fmt.Sprintf("%s Operation: %s, Path: %s, Value: %s", serverStyle.Render("Server:"), op.Op, op.Path, op.Value)
			contents = append(contents, currOp)
		}
		return &Message{
			contents: contents,
		}
	case events.EventTypeMessagesSnapshot:
		snapshot, ok := event.(*events.MessagesSnapshotEvent)
		if !ok {
			return nil
		}
		var contents []string
		for _, msg := range snapshot.Messages {
			if msg.Content != nil && msg.Role != "user" {
				contents = append(contents, *msg.Content)
			}
			if msg.ToolCalls != nil {
				for _, toolCall := range msg.ToolCalls {
					toolCallContent := serverStyle.Render("Tool Call: ") + toolCall.Function.Name + " - " + toolCall.Function.Arguments
					contents = append(contents, toolCallContent)
				}
			}
		}

		return &Message{
			contents: contents,
		}
	case events.EventTypeStepStarted:
		_, ok := event.(*events.StepStartedEvent)
		if !ok {
			return nil
		}
		log.Fatalf("Event type: %s is not yet supported. \n", eventType)
		return &Message{}

	case events.EventTypeStepFinished:
		_, ok := event.(*events.StepFinishedEvent)
		if !ok {
			return nil
		}
		log.Fatalf("Event type: %s is not yet supported. \n", eventType)
		return &Message{}
	case events.EventTypeThinkingStart:
		_, ok := event.(*events.ThinkingStartEvent)
		if !ok {
			return nil
		}
		log.Fatalf("Event type: %s is not yet supported. \n", eventType)
		return &Message{}
	case events.EventTypeThinkingEnd:
		_, ok := event.(*events.ThinkingEndEvent)
		if !ok {
			return nil
		}
		log.Fatalf("Event type: %s is not yet supported. \n", eventType)
		return &Message{}
	case events.EventTypeThinkingTextMessageStart:
		txtMsg, ok := event.(*events.ThinkingTextMessageStartEvent)
		if !ok {
			return nil
		}
		fmt.Println(txtMsg)
		log.Fatalf("Event type: %s is not yet supported. \n", eventType)
		return &Message{}
	case events.EventTypeThinkingTextMessageContent:
		_, ok := event.(*events.ThinkingTextMessageContentEvent)
		if !ok {
			return nil
		}
		log.Fatalf("Event type: %s is not yet supported. \n", eventType)
		return &Message{}

	case events.EventTypeThinkingTextMessageEnd:
		_, ok := event.(*events.ThinkingTextMessageEndEvent)
		if !ok {
			return nil
		}
		log.Fatalf("Event type: %s is not yet supported. \n", eventType)
		return &Message{}

	case events.EventTypeCustom:
		evt, ok := event.(*events.CustomEvent)
		if !ok {
			return nil
		}
		jsonData, err := json.Marshal(evt.Value)
		if err != nil {
			fmt.Println("Error marshaling JSON:", err)
			return nil
		}
		fmt.Println(evt)
		return &Message{
			contents: []string{getNameplate() + string(jsonData)},
		}

	case events.EventTypeRaw:
		_, ok := event.(*events.RawEvent)
		if !ok {
			return nil
		}
		log.Fatalf("Event type: %s is not yet supported. \n", eventType)
		return &Message{}

	default:
		// For any other event types, return a raw event
		log.Fatalf("Unknown event type: %s", eventType)
		return nil
	}
}
