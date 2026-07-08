package controllers

import (
	"fmt"
	"log"
	"net/http"
	"samsungvoicebe/config"
	"samsungvoicebe/helper"
	"samsungvoicebe/models"
	"strings"

	"github.com/gin-gonic/gin"
)

type ChatControllerV2 struct {
	config *config.Config
}

func NewChatControllerV2(cfg *config.Config) *ChatControllerV2 {
	return &ChatControllerV2{
		config: cfg,
	}
}

func (cc *ChatControllerV2) Chat(c *gin.Context) {
	var req models.ChatRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ChatResponse{
			Error: "Invalid request format",
		})
		return
	}

	log.Printf("Chat V2 request: %s", req.Message)

	screen, err := cc.determineScreen(req.Message)
	if err != nil {
		log.Printf("Error determining screen: %v", err)
		c.JSON(http.StatusBadRequest, models.ChatResponse{
			Error: err.Error(),
		})
		return
	}

	log.Printf("Determined screen: %s", screen)

	c.JSON(http.StatusOK, models.ChatResponse{
		Response: screen,
		Screen:   screen,
	})
}

func (cc *ChatControllerV2) determineScreen(message string) (string, error) {
	msg := strings.ToLower(message)

	if cc.isInappropriateContent(msg) {
		return "", fmt.Errorf("inappropriate content detected")
	}

	if len(strings.TrimSpace(msg)) < 2 {
		return "", fmt.Errorf("message too short or empty")
	}

	if cc.isGibberish(msg) {
		return "", fmt.Errorf("unable to understand the message")
	}

	playKeywords := []string{
		"main", "bermain", "permainan", "game", "play", "mulai",
		"beranda", "home", "utama", "awal", "start", "menu",
	}

	scanKeywords := []string{
		"scan", "pindai", "kamera", "camera", "foto", "gambar",
		"ambil", "capture", "barcode", "qr", "scanner",
	}

	lessonKeywords := []string{
		"belajar", "pelajaran", "lesson", "materi", "kursus",
		"tutorial", "panduan", "edukasi", "pembelajaran", "study",
	}

	analyzeKeywords := []string{
		"analisis", "analyze", "analisa", "periksa", "cek",
		"evaluasi", "tinjau", "review", "laporan", "data",
	}

	settingKeywords := []string{
		"pengaturan", "setting", "konfigurasi", "config", "atur",
		"preferensi", "opsi", "options", "setup", "setelan",
	}

	if cc.containsAny(msg, playKeywords) {
		return "play", nil
	}

	if cc.containsAny(msg, scanKeywords) {
		return "scan", nil
	}

	if cc.containsAny(msg, lessonKeywords) {
		return "lesson", nil
	}

	if cc.containsAny(msg, analyzeKeywords) {
		return "analyze", nil
	}

	if cc.containsAny(msg, settingKeywords) {
		return "setting", nil
	}

	return "", fmt.Errorf("unable to determine destination from your message")
}

func (cc *ChatControllerV2) isInappropriateContent(message string) bool {
	inappropriateWords := []string{
		"kontol", "memek", "anjing", "bangsat", "babi", "tai", "shit", "fuck",
		"bitch", "asshole", "damn", "hell",
	}

	for _, word := range inappropriateWords {
		if strings.Contains(message, word) {
			return true
		}
	}
	return false
}

func (cc *ChatControllerV2) isGibberish(message string) bool {
	cleaned := strings.ReplaceAll(message, " ", "")

	if len(cleaned) > 3 {
		repeated := true
		firstChar := cleaned[0]
		for _, char := range cleaned {
			if char != rune(firstChar) {
				repeated = false
				break
			}
		}
		if repeated {
			return true
		}
	}

	vowels := "aeiouAEIOU"
	hasVowel := false
	for _, char := range cleaned {
		if strings.ContainsRune(vowels, char) {
			hasVowel = true
			break
		}
	}

	if !hasVowel && len(cleaned) > 2 {
		return true
	}

	return false
}

func (cc *ChatControllerV2) containsAny(text string, keywords []string) bool {
	for _, keyword := range keywords {
		if strings.Contains(text, keyword) {
			return true
		}
	}
	return false
}

func (cc *ChatControllerV2) ChatWithAI(c *gin.Context) {
	var req models.ChatRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ChatResponse{
			Error: "Invalid request format",
		})
		return
	}

	log.Printf("Chat V2 AI request: %s", req.Message)

	enhancedPrompt := fmt.Sprintf(`
Analyze this user input and determine which screen they want to navigate to.
Available screens: play, scan, lesson, analyze, setting

User input: "%s"

Rules:
- If user wants to go to main menu, game, or start something: return "play"
- If user wants to scan, take photo, or use camera: return "scan"
- If user wants to learn, study, or access lessons: return "lesson"
- If user wants to analyze, check, or review data: return "analyze"
- If user wants settings, configuration, or preferences: return "setting"

Respond with ONLY the screen name (play/scan/lesson/analyze/setting), no additional text.
`, req.Message)

	response, err := helper.PromptOllama(enhancedPrompt)
	if err != nil {
		log.Printf("Azure OpenAI API error, falling back to keyword matching: %v", err)
		// Fallback to keyword matching
		screen, fallbackErr := cc.determineScreen(req.Message)
		if fallbackErr != nil {
			log.Printf("Error in fallback determination: %v", fallbackErr)
			c.JSON(http.StatusBadRequest, models.ChatResponse{
				Error: fallbackErr.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, models.ChatResponse{
			Response: screen,
			Screen:   screen,
		})
		return
	}

	screen := strings.TrimSpace(strings.ToLower(response))
	validScreens := []string{"play", "scan", "lesson", "analyze", "setting"}

	isValid := false
	for _, validScreen := range validScreens {
		if screen == validScreen {
			isValid = true
			break
		}
	}

	if !isValid {
		fallbackScreen, err := cc.determineScreen(req.Message)
		if err != nil {
			log.Printf("Error in fallback determination: %v", err)
			c.JSON(http.StatusBadRequest, models.ChatResponse{
				Error: err.Error(),
			})
			return
		}
		screen = fallbackScreen
	}

	log.Printf("Determined screen: %s", screen)

	c.JSON(http.StatusOK, models.ChatResponse{
		Response: screen,
		Screen:   screen,
	})
}
