package providers

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"

	"github.com/sashabaranov/go-openai"
)

const (
	ContentTypeHeaderKey   = "Content-Type"
	ApplicationJsonValue   = "application/json"
	AuthorizationHeaderKey = "Authorization"
	OpenAICompletionsURL   = "https://api.openai.com/v1/chat/completions"
)

var (
	ErrOpenAIResponse = errors.New("error response from OpenAI")
	ErrNoOutput       = errors.New("no output received from OpenAI")
)

type OpenAI struct {
	client *openai.Client
}

// OpenAIChatRequest represents the request payload for OpenAI API
type OpenAIChatRequest struct {
	Model    string            `json:"model"`
	Messages map[string]string `json:"messages"`
}

// OpenAIChatResponse represents the response payload for OpenAI API
type OpenAIChatResponse struct {
	Choices []struct {
		Message []struct {
			Role    string `json:"role"`
			Content string `json:"content"`
		} `json:"message"`
	} `json:"outputs"`
}

func NewOpenAI(apiKey string) *OpenAI {
	client := openai.NewClient(apiKey)
	return &OpenAI{client: client}
}

func (o *OpenAI) GenerateText(ctx context.Context, prompt string) (string, error) {

	request := openai.ChatCompletionRequest{
		Model: openai.GPT4oMini, // Testing with the mini model
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleUser,
				Content: prompt,
			},
		},
		MaxTokens: 150, //For now just testing
		Stream:    true,
	}

	stream, err := o.client.CreateChatCompletionStream(ctx, request)
	if err != nil {
		slog.ErrorContext(ctx, "Error while returning completion stream", "error", err)
		return "", err
	}
	defer stream.Close()

	fmt.Printf("Response: ")
	var responseText string
	for {
		response, err := stream.Recv()
		if errors.Is(err, io.EOF) {
			break
		}

		if err != nil {
			slog.ErrorContext(ctx, "Error while streaming result", "error", err)
			return "", err
		}

		fmt.Printf(response.Choices[0].Delta.Content)
		responseText += response.Choices[0].Delta.Content
	}
	return responseText, nil
}
