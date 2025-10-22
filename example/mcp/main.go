package main

import (
	"bytes"
	_ "embed"
	"os"
	"os/exec"
	"path/filepath"
	"text/template"

	mcp "github.com/metoro-io/mcp-golang"
	"github.com/metoro-io/mcp-golang/transport/http"
)

//go:embed data/pr_review.md
var prPrompt string

//go:embed data/100mistakes.html
var commonGoMistakes string

type Content struct {
	Title       string  `json:"title" jsonschema:"required,description=The title to submit"`
	Description *string `json:"description" jsonschema:"description=The description to submit"`
}
type MyFunctionsArguments struct {
	Submitter string  `json:"submitter" jsonschema:"required,description=The name of the thing calling this tool (openai, google, claude, etc)"`
	Content   Content `json:"content" jsonschema:"required,description=The content of the message"`
}

type AppendFileArguments struct {
	Path string `json:"path" jsonschema:"required,description=Filepath to append to"`
	Text string `json:"text" jsonschema:"required,description=Text to append"`
}

type ReadFileArguments struct {
	Path string `json:"path" jsonschema:"required,description=Filepath to read from"`
}

type ExecCommandArguments struct {
	Command string   `json:"command" jsonschema:"required,description=Command to execute"`
	Args    []string `json:"args" jsonschema:"description=Arguments to pass to the command"`
}

type PRArguments struct {
	Repo string `json:"repo" jsonschema:"required,description=The root folder of the repository for review"`
}

func readFile(path string) (string, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func appendFile(path, content string) (string, error) {
	// Ensure parent directory exists if a directory is specified
	dir := filepath.Dir(path)
	if dir != "." && dir != "" {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return "", err
		}
	}

	// Open the file in append mode, create it if it doesn't exist
	f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return "", err
	}
	defer f.Close()

	// Append the content
	if _, err := f.WriteString(content); err != nil {
		return "", err
	}

	// Read and return the current contents
	return readFile(path)
}

func execCommand(command string, args []string) (string, error) {
	cmd := exec.Command(command, args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", err
	}
	return string(output), nil
}

func main() {
	done := make(chan struct{})

	transport := http.NewHTTPTransport("/mcp")
	transport.WithAddr(":8080")

	server := mcp.NewServer(transport)
	err := server.RegisterTool("append_file", "Append content to a file, creates the file if doesn't exist. Returns the final contents of the file.", func(arguments AppendFileArguments) (*mcp.ToolResponse, error) {
		content, err := appendFile(arguments.Path, arguments.Text)
		if err != nil {
			return nil, err
		}
		return mcp.NewToolResponse(mcp.NewTextContent(content)), nil
	})
	if err != nil {
		panic(err)
	}

	err = server.RegisterTool("exec_command", "Execute a command and return the output.", func(arguments ExecCommandArguments) (*mcp.ToolResponse, error) {
		content, err := execCommand(arguments.Command, arguments.Args)
		if err != nil {
			return nil, err
		}
		return mcp.NewToolResponse(mcp.NewTextContent(content)), nil
	})
	if err != nil {
		panic(err)
	}

	err = server.RegisterTool("read_file", "Read the contents of a file.", func(arguments ReadFileArguments) (*mcp.ToolResponse, error) {
		content, err := readFile(arguments.Path)
		if err != nil {
			return nil, err
		}
		return mcp.NewToolResponse(mcp.NewTextContent(content)), nil
	})
	if err != nil {
		panic(err)
	}

	createPrDescription := "Prompt for creating a pull request in markdown files"
	err = server.RegisterPrompt("create_pr", createPrDescription, func(arguments PRArguments) (*mcp.PromptResponse, error) {
		tpl := template.New("prPrompt")
		tpl, err := tpl.Parse(prPrompt)
		if err != nil {
			return nil, err
		}

		var buff bytes.Buffer
		if err := tpl.Execute(&buff, arguments); err != nil {
			return nil, err
		}
		return mcp.NewPromptResponse(createPrDescription, mcp.NewPromptMessage(mcp.NewTextContent(buff.String()), mcp.RoleUser)), nil
	})
	if err != nil {
		panic(err)
	}

	err = server.RegisterResource("file://go_mistakes.html", "go_mistakes", "Common Go Mistakes", "text/html", func() (*mcp.ResourceResponse, error) {
		return mcp.NewResourceResponse(mcp.NewTextEmbeddedResource("file://go_mistakes.html", commonGoMistakes, "text/html")), nil
	})
	if err != nil {
		panic(err)
	}

	err = server.Serve()
	if err != nil {
		panic(err)
	}

	<-done
}
