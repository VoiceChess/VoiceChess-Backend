package services

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/notnil/chess"
	"samsungvoicebe/helper"
	"samsungvoicebe/models"
	"samsungvoicebe/repo"
)

type GameplayService struct {
	gameplayRepo    *repo.GameplayRepo
	analysisService *AnalysisService
}

func NewGameplayService(gameplayRepo *repo.GameplayRepo, analysisService *AnalysisService) *GameplayService {
	return &GameplayService{
		gameplayRepo:    gameplayRepo,
		analysisService: analysisService,
	}
}

func (s *GameplayService) GameBelongsToUser(gameID, userID string) (bool, error) {
	return s.gameplayRepo.GameBelongsToUser(gameID, userID)
}

func (s *GameplayService) PlayerMove(gameID *string, fen, move, botLevel string) (models.BotMove, error) {
	if gameID != nil {
		err := s.gameplayRepo.GameMove(*gameID, fen, move)
		if err != nil {
			err = fmt.Errorf("GameplayService-PlayerMove-GameMove: %w", err)
			fmt.Printf("gameID: %s, fen: %s, move: %s", *gameID, fen, move)
			return models.BotMove{}, err
		}
	}

	analysisResult, err := s.analysisService.StockfishAnalyze(fen, botLevel)
	if err != nil {
		err = fmt.Errorf("GameplayService-PlayerMove-GameMove: %w", err)
		return models.BotMove{}, err
	}

	if gameID != nil {
		err = s.gameplayRepo.GameMove(*gameID, analysisResult.Fen, analysisResult.BestMove)
		if err != nil {
			err = fmt.Errorf("GameplayService-PlayerMove-GameMove: %w", err)
			fmt.Printf("gameID: %s, fen: %s, move: %s", *gameID, fen, move)
			return models.BotMove{}, err
		}
	}

	var botMove models.BotMove

	botMove = models.BotMove{
		Fen:  analysisResult.Fen,
		Move: analysisResult.BestMove,
	}

	return botMove, nil
}

func (s *GameplayService) CreateGame(userID string) (string, error) {
	gameID, err := s.gameplayRepo.CreateGame(userID)
	if err != nil {
		err = fmt.Errorf("GameplayService-CreateGame-CreateGame: %w", err)
		return "", err
	}
	return gameID, nil
}

func (s *GameplayService) GetHint(fen string) (string, error) {
	prompt := fmt.Sprintf(models.HintPrompt, fen)
	hint, err := helper.PromptOllamaWithTimeout(prompt, 5*time.Second)
	if err != nil {
		analysis, stockfishErr := s.analysisService.StockfishAnalyze(fen, "medium")
		if stockfishErr != nil {
			return "Look for checks, captures, and threats before moving.", nil
		}
		return fmt.Sprintf("Try improving your position with a move like %s.", analysis.BestMove), nil
	}

	return hint, nil
}

func (s *GameplayService) PlayerMoveByVoiceTranscription(fen, transcription string) (models.PlayerMoveByTranscription, error) {
	if move, ok := coordinateMoveFromText(transcription); ok {
		return applyVoiceMove(fen, move)
	}

	prompt := fmt.Sprintf(models.MoveFromDescriptionPrompt, fen, transcription)
	move, err := helper.PromptOllamaWithTimeout(prompt, 5*time.Second)
	if err != nil {
		return models.PlayerMoveByTranscription{}, fmt.Errorf("GameplayService-PlayerMoveByVoiceTranscription-PromptOllama: %w", err)
	}
	move = strings.TrimSpace(move)

	if move == models.InvalidMove {
		err := fmt.Errorf("GameplayService-PlayerMoveByVoiceTranscription-PromptAzureOpenAI: invalid move from transcription")
		return models.PlayerMoveByTranscription{}, err
	}

	if strings.Contains(move, "location") {
		playerMove := models.PlayerMoveByTranscription{
			Status: "Location",
			Move:   strings.Replace(move, "location:", "", -1),
			Fen:    fen,
		}

		return playerMove, nil
	}

	return applyVoiceMove(fen, move)
}

func coordinateMoveFromText(transcription string) (string, bool) {
	text := strings.ToLower(transcription)
	replacements := map[string]string{
		"one": "1", "satu": "1",
		"two": "2", "dua": "2",
		"three": "3", "tiga": "3",
		"four": "4", "empat": "4",
		"five": "5", "lima": "5",
		"six": "6", "enam": "6",
		"seven": "7", "tujuh": "7",
		"eight": "8", "delapan": "8",
	}
	for from, to := range replacements {
		text = regexp.MustCompile(`\b`+from+`\b`).ReplaceAllString(text, to)
	}
	text = regexp.MustCompile(`\b([a-h])\s+([1-8])\b`).ReplaceAllString(text, `$1$2`)
	match := regexp.MustCompile(`\b([a-h][1-8])\b\s*(?:to|tu|ke|menuju|pindah ke)\s*\b([a-h][1-8])\b`).FindStringSubmatch(text)
	if len(match) != 3 {
		return "", false
	}
	return match[1] + match[2], true
}

func applyVoiceMove(fen, move string) (models.PlayerMoveByTranscription, error) {
	position, err := chess.FEN(fen)
	if err != nil {
		return models.PlayerMoveByTranscription{}, fmt.Errorf("GameplayService-PlayerMoveByVoiceTranscription-chess.FEN: %w", err)
	}

	game := chess.NewGame(position)
	if err := game.MoveStr(move); err != nil {
		return models.PlayerMoveByTranscription{}, fmt.Errorf("GameplayService-PlayerMoveByVoiceTranscription-game.Move: %w", err)
	}

	return models.PlayerMoveByTranscription{
		Status: "Move",
		Move:   move,
		Fen:    game.FEN(),
	}, nil
}

func (s *GameplayService) UndoMove(gameID string) error {
	err := s.gameplayRepo.UndoMove(gameID)
	if err != nil {
		err = fmt.Errorf("GameplayService-UndoMove-UndoMove: %w", err)
		return err
	}
	return nil
}
