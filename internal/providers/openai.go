package providers

import (
	"context"
	"fmt"
)

type OpenAI struct{}

func (o *OpenAI) GenerateText(ctx context.Context, prompt string) (string, error) {
	return fmt.Sprintf("OpenAI response to '%s'", prompt), nil
}
