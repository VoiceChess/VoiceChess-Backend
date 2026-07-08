package helper

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

type fenDetectRequest struct {
	ImgB64   string `json:"img_b64"`
	NumTries int    `json:"num_tries,omitempty"`
}

type fenDetectResponse struct {
	Success bool   `json:"success"`
	Fen     string `json:"fen"`
	Detail  string `json:"detail,omitempty"`
}

// DetectFenFromPicture sends an image to the self-hosted chess-diagram-to-FEN
// model (AI-VoiceChess /api/detect) and returns the detected FEN string.
func DetectFenFromPicture(imageFile []byte) (string, error) {
	url := os.Getenv("FEN_DETECT_URL")
	if url == "" {
		return "", fmt.Errorf("FEN_DETECT_URL not configured")
	}

	reqBody := fenDetectRequest{
		ImgB64:   base64.StdEncoding.EncodeToString(imageFile),
		NumTries: 5,
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

	// ponytail: CNN inference on CPU is ~5-30s depending on num_tries
	client := &http.Client{Timeout: 60 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to call FEN detect API: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("FEN detect API returned status %d: %s", resp.StatusCode, string(body))
	}

	var out fenDetectResponse
	if err := json.Unmarshal(body, &out); err != nil {
		return "", fmt.Errorf("failed to parse response: %w", err)
	}
	if !out.Success || out.Fen == "" {
		return "", fmt.Errorf("FEN detection failed: %s", out.Detail)
	}

	return out.Fen, nil
}
