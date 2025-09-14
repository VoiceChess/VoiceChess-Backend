package helper

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
	"samsungvoicebe/models"
)

func PromptGemini(prompt string) string {
	ctx := context.Background()

	client, err := genai.NewClient(ctx, option.WithAPIKey(os.Getenv("GEMINI_API_KEY")))
	if err != nil {
		log.Fatalf("helper-PromptGemini-genai.NewClient %v", err)
	}
	defer client.Close()

	model := client.GenerativeModel(models.GeminiModel)

	resp, err := model.GenerateContent(ctx, genai.Text(prompt))
	if err != nil {
		log.Fatalf("helper-PromptGemini-model.GenerateContent: %v", err)
	}

	if len(resp.Candidates) == 0 || len(resp.Candidates[0].Content.Parts) == 0 {
		return ""
	}

	return fmt.Sprintf("%v", resp.Candidates[0].Content.Parts[0])
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
		genai.ImageData("image/png", imageFile),
	)

	if err != nil {
		return "", fmt.Errorf("helper-AnalyzePictureWithGemini-model.GenerateContent: %w", err)
	}

	if len(resp.Candidates) == 0 || len(resp.Candidates[0].Content.Parts) == 0 {
		return "", fmt.Errorf("no response generated")
	}

	return fmt.Sprintf("%v", resp.Candidates[0].Content.Parts[0]), nil
}
