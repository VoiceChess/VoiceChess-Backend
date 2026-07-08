package helper

import (
	"context"
	"fmt"
	"os"

	"samsungvoicebe/models"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

func PromptGemini(prompt string) (string, error) {
	ctx := context.Background()

	client, err := genai.NewClient(ctx, option.WithAPIKey(os.Getenv("GEMINI_API_KEY")))
	if err != nil {
		return "", fmt.Errorf("helper-PromptGemini-genai.NewClient: %w", err)
	}
	defer client.Close()

	model := client.GenerativeModel(models.GeminiModel)

	resp, err := model.GenerateContent(ctx, genai.Text(prompt))
	if err != nil {
		return "", fmt.Errorf("helper-PromptGemini-model.GenerateContent: %w", err)
	}

	if len(resp.Candidates) == 0 || len(resp.Candidates[0].Content.Parts) == 0 {
		return "", nil
	}

	return fmt.Sprintf("%v", resp.Candidates[0].Content.Parts[0]), nil
}

func AnalyzePictureWithGemini(imageFile []byte, prompt string) (string, error) {
	ctx := context.Background()

	client, err := genai.NewClient(ctx, option.WithAPIKey(os.Getenv("GEMINI_API_KEY")))
	if err != nil {
		return "", fmt.Errorf("helper-AnalyzePictureWithGemini-genai.NewClient: %w", err)
	}
	defer client.Close()

	model := client.GenerativeModel(models.GeminiModel)

	resp, err := model.GenerateContent(ctx,
		genai.Text(prompt),
		genai.ImageData("png", imageFile),
	)

	if err != nil {
		return "", fmt.Errorf("helper-AnalyzePictureWithGemini-model.GenerateContent: %w", err)
	}

	if len(resp.Candidates) == 0 || len(resp.Candidates[0].Content.Parts) == 0 {
		return "", fmt.Errorf("no response generated")
	}

	return fmt.Sprintf("%v", resp.Candidates[0].Content.Parts[0]), nil
}
