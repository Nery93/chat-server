package ai

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

type AIRequest struct {
	Model              string             `json:"model"`
	Messages           []AIMessage        `json:"messages"`
	Stream             bool               `json:"stream"`
	MaxTokens          int                `json:"max_tokens"`
	ChatTemplateKwargs ChatTemplateKwargs `json:"chat_template_kwargs"`
}

type ChatTemplateKwargs struct {
	EnableThinking bool `json:"enable_thinking"`
}

type AIMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type AIResponse struct {
	Choices []struct {
		Message AIMessage `json:"message"`
	} `json:"choices"`
	Usage struct {
		PromptTokens     int `json:"prompt_tokens"`
		CompletionTokens int `json:"completion_tokens"`
		TotalTokens      int `json:"total_tokens"`
	} `json:"usage"`
}

func AskGenerateResponse(prompt string) (string, error) {
	requestJSON, err := json.Marshal(NewAIRequest(prompt))
	if err != nil {
		return "", err
	}

	request, err := http.NewRequest(http.MethodPost, "https://integrate.api.nvidia.com/v1/chat/completions",
		bytes.NewReader(requestJSON),
	)
	if err != nil {
		return "", err
	}

	apiKey := os.Getenv("NVIDIA_API_KEY")
	if apiKey == "" {
		return "", errors.New("NVIDIA_API_KEY não está configurada")
	}

	request.Header.Set("Authorization", "Bearer "+apiKey)
	request.Header.Set("Content-Type", "application/json")

	response, err := (&http.Client{}).Do(request)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		return "", err
	}

	if response.StatusCode < 200 || response.StatusCode >= 300 {
		return "", fmt.Errorf("NVIDIA respondeu com status %d: %s", response.StatusCode, responseBody)
	}

	var result AIResponse
	if err := json.Unmarshal(responseBody, &result); err != nil {
		return "", err
	}

	if len(result.Choices) == 0 {
		return "", errors.New("NVIDIA não devolveu escolhas")
	}

	log.Printf(
		"Uso NVIDIA: prompt=%d resposta=%d total=%d",
		result.Usage.PromptTokens,
		result.Usage.CompletionTokens,
		result.Usage.TotalTokens,
	)

	return result.Choices[0].Message.Content, nil
}

func NewAIRequest(prompt string) *AIRequest {
	return &AIRequest{
		Model: "nvidia/nemotron-3-ultra-550b-a55b",
		Messages: []AIMessage{
			{
				Role:    "user",
				Content: prompt,
			},
		},
		Stream:             false,
		MaxTokens:          1500,
		ChatTemplateKwargs: ChatTemplateKwargs{EnableThinking: false},
	}
}
