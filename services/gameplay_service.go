package services

import (
	"fmt"
	"strings"

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

func (s *GameplayService) PlayerMove(gameID *string, fen, move, botLevel string) (models.BotMove, error) {
	if gameID != nil {
		err := s.gameplayRepo.GameMove(*gameID, fen, move)
		if err != nil {
			err = fmt.Errorf("GameplayService-PlayerMove-GameMove: %w", err)
			fmt.Printf("gameID: %s, fen: %s, move: %s", gameID, fen, move)
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
			fmt.Printf("gameID: %s, fen: %s, move: %s", gameID, fen, move)
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
	hint := helper.PromptGemini(prompt)

	return hint, nil
}

func (s *GameplayService) PlayerMoveByVoiceTranscription(fen, transcription string) (models.PlayerMoveByTranscription, error) {
	prompt := fmt.Sprintf(models.MoveFromDescriptionPrompt, fen, transcription)
	move := helper.PromptGemini(prompt)
	move = strings.TrimSpace(move)

	if move == models.InvalidMove {
		err := fmt.Errorf("GameplayService-PlayerMoveByVoiceTranscription-PromptGemini: invalid move from transcription")
		return models.PlayerMoveByTranscription{}, err
	}

	position, err := chess.FEN(fen)
	if err != nil {
		err = fmt.Errorf("GameplayService-PlayerMoveByVoiceTranscription-chess.FEN: %w", err)
		return models.PlayerMoveByTranscription{}, err
	}

	game := chess.NewGame(position)
	err = game.MoveStr(move)
	if err != nil {
		err = fmt.Errorf("GameplayService-PlayerMoveByVoiceTranscription-game.Move: %w", err)
		return models.PlayerMoveByTranscription{}, err
	}

	playerMoveFEN := game.FEN()

	playerMove := models.PlayerMoveByTranscription{
		Move: move,
		Fen:  playerMoveFEN,
	}

	return playerMove, nil
}
