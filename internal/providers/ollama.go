package providers

import (
	"context"
	"fmt"
)

// Ollama struct represents the Ollama implementation
type Ollama struct {
	modelName string
}

// GenerateText generates text using the Ollama model
func (o *Ollama) GenerateText(ctx context.Context, prompt string) (string, error) {
	// Mocked response for illustration
	return fmt.Sprintf("Ollama (%s) response to '%s'", o.modelName, prompt), nil
}
