package providers

import "fmt"

// AIProviderFactory returns the appropriate AI provider based on the model name
func AIProviderFactory(apiKey, modelName string) (AIModel, error) {
	switch modelName {
	case "gpt-4o-mini":
		return NewOpenAI(apiKey), nil
	case "llama3.1":
		return &Ollama{modelName: modelName}, nil
	default:
		return nil, fmt.Errorf("unsupported model name: %s", modelName)
	}
}
