package main

import (
	"context"
	"fmt"
	"log"

	"github.com/firebase/genkit/go/ai"
	"github.com/firebase/genkit/go/genkit"
	"github.com/firebase/genkit/go/plugins/googlegenai"
)

func main() {
	ctx := context.Background()
	g := genkit.Init(
		ctx,
		genkit.WithDefaultModel("googleai/gemini-2.5-flash"),
		genkit.WithPlugins(&googlegenai.GoogleAI{}),
	)

	finalResult, err := genkit.Generate(ctx, g,
		ai.WithPrompt("Tell me about why GenKit is awesome"))

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(finalResult.Text())
}
