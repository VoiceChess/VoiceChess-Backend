package services

import (
	"fmt"
	"strings"

	"github.com/notnil/chess"
	"github.com/notnil/chess/uci"
	"samsungvoicebe/helper"
	"samsungvoicebe/models"
	"samsungvoicebe/repo"
)

type AnalysisService struct {
	analysisRepo *repo.AnalysisRepo
}

func NewAnalysisService(analysisRepo *repo.AnalysisRepo) *AnalysisService {
	return &AnalysisService{analysisRepo: analysisRepo}
}

func (a *AnalysisService) StockfishAnalyze(fen string, botLevel string) (models.StockfishAnalysisResult, error) {
	var analysisResult models.StockfishAnalysisResult

	engine, err := uci.New("stockfish")
	if err != nil {
		err = fmt.Errorf("AnalysisService-StockfishAnalyze-uci.New: %w", err)
		return models.StockfishAnalysisResult{}, err
	}
	defer engine.Close()

	err = engine.Run(uci.CmdUCI, uci.CmdIsReady)
	if err != nil {
		err = fmt.Errorf("AnalysisService-StockfishAnalyze-engine.Run-uci.CmdUCI-uci.CmdIsReady: %w", err)
		return models.StockfishAnalysisResult{}, err
	}

	position, err := chess.FEN(fen)
	if err != nil {
		err = fmt.Errorf("AnalysisService-StockfishAnalyze-chess.FEN: %w", err)
		return models.StockfishAnalysisResult{}, err
	}

	game := chess.NewGame(position)

	err = engine.Run(uci.CmdPosition{Position: game.Position()})
	if err != nil {
		err = fmt.Errorf("AnalysisService-StockfishAnalyze-engine.Run-uci.CmdPosition: %w", err)
		return models.StockfishAnalysisResult{}, err
	}

	var depth int
	switch botLevel {
	case "easy":
		depth = models.BotLevelEasy
	case "medium":
		depth = models.BotLevelMedium
	case "hard":
		depth = models.BotLevelhard
	default:
		depth = models.BotLevelMedium
	}

	searchBestMove := uci.CmdGo{Depth: depth}

	err = engine.Run(searchBestMove)
	if err != nil {
		err = fmt.Errorf("AnalysisService-StockfishAnalyze-engine.Run-searchBestMove: %w", err)
		return models.StockfishAnalysisResult{}, err
	}

	bestMove := engine.SearchResults().BestMove

	err = game.Move(bestMove)
	if err != nil {
		err = fmt.Errorf("AnalysisService-StockfishAnalyze-game.Move: %w", err)
		return models.StockfishAnalysisResult{}, err
	}

	analysisResult = models.StockfishAnalysisResult{
		BestMove: bestMove.String(),
		Fen:      game.FEN(),
	}

	return analysisResult, nil
}

func (a *AnalysisService) GetGameHistoryList(userID string) ([]models.Game, error) {
	games, err := a.analysisRepo.GetGameHistoryList(userID)
	if err != nil {
		err = fmt.Errorf("AnalysisService-GetGameHistoryList-GetGameHistoryList: %w", err)
		return []models.Game{}, err
	}

	return games, nil
}

func (a *AnalysisService) GetAnalyzedMoveByOrder(moveOrder int, gameID string) (models.MoveAnalysis, error) {
	var analyzedMove models.MoveAnalysis

	move, err := a.analysisRepo.GetMoveByOrder(moveOrder, gameID)
	if err != nil {
		err = fmt.Errorf("AnalysisService-GetAnalyzedMoveByOrder-GetMoveByOrder: %w", err)
		return models.MoveAnalysis{}, err
	}

	analyzedMove = models.MoveAnalysis{
		Move: move.Move,
		Fen:  move.Fen,
	}

	stockfishResult, err := a.StockfishAnalyze(move.Fen, "hard")
	if err != nil {
		err = fmt.Errorf("AnalysisService-GetAnalyzedMoveByOrder-StockfishAnalyze: %w", err)
		return models.MoveAnalysis{}, err
	}

	analyzedMove.BestMove = stockfishResult.BestMove

	return analyzedMove, nil
}

func (a *AnalysisService) GetFenFromPicture(imageFile []byte) (string, error) {
	fen, err := helper.AnalyzePictureWithGemini(imageFile, models.GetFenFromPicturePrompt)
	fen = strings.TrimSpace(fen)
	if err != nil || fen == models.InvalidImage {
		err = fmt.Errorf("AnalysisService-GetFenFromPicture-AnalyzePictureWithGemini: %w", err)
		return "", err
	}

	return fen, nil
}
