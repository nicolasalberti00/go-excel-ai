package providers

import "fmt"

// AIProviderFactory returns the appropriate AI provider based on the model name
func AIProviderFactory(apiKey, modelName string) (AIModel, error) {
	switch modelName {
	case "gpt-4o-mini":
		return NewOpenAI(apiKey), nil
	case "llama3.1":
		llm, err := NewOllama(modelName)
		if err != nil {
			return nil, fmt.Errorf("error during ollama model initialization: %s", err)
		}
		return llm, nil
	default:
		return nil, fmt.Errorf("unsupported model name: %s", modelName)
	}
}
