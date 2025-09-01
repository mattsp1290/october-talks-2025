package main

// Useful for quickly testing the client

import (
	"context"
	"fmt"
	"os"

	"github.com/mattsp1290/october-talks-2025/example/client/internal/agent"
	"github.com/mattsp1290/october-talks-2025/example/client/internal/message"
)

func onMsg(msg *message.Message) {
	for _, s := range msg.Strings() {
		fmt.Println(s)
	}
}

func Host() string {
	host := os.Getenv("SERVER_HOST")
	if host == "" {
		host = "localhost"
	}
	return host
}

func main() {
	if len(os.Args) < 1 {
		panic("No arguments provided")
	}
	endpoint := os.Args[1]
	url := "http://" + Host() + ":8000/" + endpoint
	err := agent.Chat(context.Background(), "test", url, onMsg)
	if err != nil {
		panic(err)
	}
}
