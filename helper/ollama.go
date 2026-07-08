package helper

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

type ollamaRequest struct {
	Model    string          `json:"model"`
	Messages []ollamaMessage `json:"messages"`
	Stream   bool            `json:"stream"`
}

type ollamaMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ollamaResponse struct {
	Message ollamaMessage `json:"message"`
	Error   string        `json:"error,omitempty"`
}

// PromptOllama sends a single-prompt chat request to the Ollama gateway and
// returns the model's text reply. Drop-in replacement for PromptAzureOpenAI.
func PromptOllama(prompt string) (string, error) {
	url := os.Getenv("OLLAMA_API_URL")
	model := os.Getenv("OLLAMA_MODEL_NAME")
	if url == "" {
		return "", fmt.Errorf("OLLAMA_API_URL not configured")
	}
	if model == "" {
		model = "qwen2.5:3b"
	}

	reqBody := ollamaRequest{
		Model:    model,
		Stream:   false,
		Messages: []ollamaMessage{{Role: "user", Content: prompt}},
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	if key := os.Getenv("OLLAMA_API_KEY"); key != "" {
		req.Header.Set("Authorization", "Bearer "+key)
	}

	// ponytail: 30s covers CPU inference on qwen2.5:3b; raise if you host a bigger model
	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to call Ollama API: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("ollama API returned status %d: %s", resp.StatusCode, string(body))
	}

	var out ollamaResponse
	if err := json.Unmarshal(body, &out); err != nil {
		return "", fmt.Errorf("failed to parse response: %w", err)
	}
	if out.Error != "" {
		return "", fmt.Errorf("ollama API error: %s", out.Error)
	}

	return out.Message.Content, nil
}
