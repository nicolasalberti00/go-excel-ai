package providers

import "context"

type AIModel interface {
	GenerateText(ctx context.Context, prompt string) (string, error)
}
