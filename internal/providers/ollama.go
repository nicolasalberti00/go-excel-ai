package providers

import (
	"context"
	"fmt"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/ollama"
	"log/slog"
)

type Ollama struct {
	llm *ollama.LLM
}

func NewOllama(modelName string) (*Ollama, error) {
	llm, err := ollama.New(ollama.WithModel(modelName))
	if err != nil {
		slog.Error("Error when creating new ollama instance", "Error", err.Error())
		return &Ollama{}, err
	}
	return &Ollama{llm: llm}, nil
}

func (o *Ollama) GenerateText(ctx context.Context, prompt string) (string, error) {
	resp, err := llms.GenerateFromSinglePrompt(ctx, o.llm, prompt,
		llms.WithStreamingFunc(func(ctx context.Context, chunk []byte) error {
			fmt.Printf("chunk len=%d: %s\n", len(chunk), chunk)
			return nil
		}))
	if err != nil {
		slog.ErrorContext(ctx, "Error when generating text", "Error", err.Error())
		return "", err
	}
	return resp, nil
}
