package services

import (
	"encoding/json"
	"fmt"
	"log"
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

func logGameplayEvent(event string, fields map[string]any) {
	fields["event"] = event
	payload, err := json.Marshal(fields)
	if err != nil {
		log.Printf("gameplay_event_marshal_failed: %v", err)
		return
	}
	log.Println(string(payload))
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
	logGameplayEvent("voice_move_received", map[string]any{
		"transcription": transcription,
		"text_length":   len(transcription),
	})

	if move, ok := coordinateMoveFromText(transcription); ok {
		logGameplayEvent("voice_move_parser_matched", map[string]any{
			"strategy": "coordinate",
			"move":     move,
		})
		return applyVoiceMove(fen, move)
	}

	logGameplayEvent("voice_move_parser_missed", map[string]any{
		"strategy":      "coordinate",
		"fallback":      "ollama",
		"transcription": transcription,
	})
	prompt := fmt.Sprintf(models.MoveFromDescriptionPrompt, fen, transcription)
	logGameplayEvent("voice_move_ollama_started", map[string]any{
		"timeout_ms": 5000,
	})
	move, err := helper.PromptOllamaWithTimeout(prompt, 5*time.Second)
	if err != nil {
		logGameplayEvent("voice_move_ollama_failed", map[string]any{
			"error": err.Error(),
		})
		return models.PlayerMoveByTranscription{}, fmt.Errorf("GameplayService-PlayerMoveByVoiceTranscription-PromptOllama: %w", err)
	}
	move = strings.TrimSpace(move)
	logGameplayEvent("voice_move_ollama_completed", map[string]any{
		"move": move,
	})

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
		"ah": "a", "ha": "h",
		"be": "b", "bi": "b",
		"ce": "c", "si": "c", "see": "c",
		"de": "d", "di": "d", "the": "d",
		"eh": "e", "i": "e",
		"ef": "f",
		"ge": "g", "ji": "g",
		"one": "1", "satu": "1",
		"two": "2", "dua": "2",
		"three": "3", "tiga": "3",
		"four": "4", "empat": "4",
		"five": "5", "lima": "5", "rima": "5", "delima": "5", "de lima": "5",
		"six": "6", "enam": "6",
		"seven": "7", "tujuh": "7",
		"eight": "8", "delapan": "8",
		"kadeh": "ke d", "kade": "ke d", "ke de": "ke d",
	}
	for from, to := range replacements {
		text = regexp.MustCompile(`\b`+from+`\b`).ReplaceAllString(text, to)
	}
	text = regexp.MustCompile(`\bke\s+b\s+e\s*([1-8])\b`).ReplaceAllString(text, `ke d$1`)
	text = regexp.MustCompile(`\bke\s+b\s+e([1-8])\b`).ReplaceAllString(text, `ke d$1`)
	text = regexp.MustCompile(`\b([a-h])\s+([1-8])\b`).ReplaceAllString(text, `$1$2`)
	text = regexp.MustCompile(`\b([a-h])([1-8])\s+(?:2|to|tu|ke|menuju|pindah ke)\s+([a-h])\s*([1-8])\b`).ReplaceAllString(text, `$1$2 ke $3$4`)
	text = regexp.MustCompile(`\b([a-h])\s*([1-8])\s+([a-h])\s*([1-8])\b`).ReplaceAllString(text, `$1$2 ke $3$4`)
	match := regexp.MustCompile(`\b([a-h][1-8])\b\s*(?:2|to|tu|ke|menuju|pindah ke)\s*\b([a-h][1-8])\b`).FindStringSubmatch(text)
	if len(match) != 3 {
		return "", false
	}
	return match[1] + match[2], true
}

func applyVoiceMove(fen, move string) (models.PlayerMoveByTranscription, error) {
	position, err := chess.FEN(fen)
	if err != nil {
		logGameplayEvent("voice_move_apply_failed", map[string]any{
			"move":  move,
			"stage": "fen",
			"error": err.Error(),
		})
		return models.PlayerMoveByTranscription{}, fmt.Errorf("GameplayService-PlayerMoveByVoiceTranscription-chess.FEN: %w", err)
	}

	game := chess.NewGame(position)
	if err := game.MoveStr(move); err != nil {
		logGameplayEvent("voice_move_apply_failed", map[string]any{
			"move":  move,
			"stage": "move",
			"error": err.Error(),
		})
		return models.PlayerMoveByTranscription{}, fmt.Errorf("GameplayService-PlayerMoveByVoiceTranscription-game.Move: %w", err)
	}

	logGameplayEvent("voice_move_apply_completed", map[string]any{
		"move": move,
	})
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
