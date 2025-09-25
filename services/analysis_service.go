package services

import (
	"fmt"
	"math"
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
	score := engine.SearchResults().Info.Score

	if score.Mate != 0 {
		analysisResult.Eval = score.Mate
		analysisResult.IsMate = true
	} else {
		analysisResult.Eval = score.CP
		analysisResult.IsMate = false
	}

	err = game.Move(bestMove)
	if err != nil {
		err = fmt.Errorf("AnalysisService-StockfishAnalyze-game.Move: %w", err)
		return models.StockfishAnalysisResult{}, err
	}

	analysisResult.Fen = game.FEN()
	analysisResult.BestMove = bestMove.String()

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
	var currMoveAnalysis models.StockfishAnalysisResult

	firstMove, err := a.analysisRepo.GetMoveByOrder(1, gameID)
	if err != nil {
		err = fmt.Errorf("AnalysisService-GetAnalyzedMoveByOrder-GetMoveByOrder: %w", err)
		return models.MoveAnalysis{}, err
	}

	currMove, err := a.analysisRepo.GetMoveByOrder(moveOrder, gameID)
	if err != nil {
		err = fmt.Errorf("AnalysisService-GetAnalyzedMoveByOrder-GetMoveByOrder: %w", err)
		return models.MoveAnalysis{}, err
	}

	if moveOrder > 2 || firstMove.Move != "black" {
		prevMoveAnalysis, err := a.getMoveAnalysis(moveOrder-1, gameID)
		if err != nil {
			err = fmt.Errorf("AnalysisService-GetAnalyzedMoveByOrder-getMoveAnalysis: %w", err)
			return models.MoveAnalysis{}, err
		}
		currMoveAnalysis, err = a.getMoveAnalysis(moveOrder, gameID)
		if err != nil {
			err = fmt.Errorf("AnalysisService-GetAnalyzedMoveByOrder-getMoveAnalysis: %w", err)
			return models.MoveAnalysis{}, err
		}

		if prevMoveAnalysis.IsMate && currMoveAnalysis.IsMate {
			analyzedMove.MoveGrade = a.classifyMove(nil, nil, &prevMoveAnalysis.Eval, &currMoveAnalysis.Eval)
		} else if prevMoveAnalysis.IsMate && !currMoveAnalysis.IsMate {
			analyzedMove.MoveGrade = a.classifyMove(nil, &currMoveAnalysis.Eval, &prevMoveAnalysis.Eval, nil)
		} else if !prevMoveAnalysis.IsMate && currMoveAnalysis.IsMate {
			analyzedMove.MoveGrade = a.classifyMove(&prevMoveAnalysis.Eval, nil, nil, &currMoveAnalysis.Eval)
		} else {
			analyzedMove.MoveGrade = a.classifyMove(&prevMoveAnalysis.Eval, &currMoveAnalysis.Eval, nil, nil)
		}
	} else if (moveOrder == 2 && firstMove.Move == "black") || moveOrder == 1 {
		currMoveAnalysis, err = a.getMoveAnalysis(moveOrder, gameID)
		if err != nil {
			err = fmt.Errorf("AnalysisService-GetAnalyzedMoveByOrder-getMoveAnalysis: %w", err)
			return models.MoveAnalysis{}, err
		}
		zeroPtr := 0
		analyzedMove.MoveGrade = a.classifyMove(&zeroPtr, &currMoveAnalysis.Eval, nil, nil)
	}

	analyzedMove.Move = currMove.Move
	analyzedMove.Fen = currMove.Fen
	analyzedMove.BestMove = currMoveAnalysis.BestMove
	analyzedMove.IsEvalMate = currMoveAnalysis.IsMate

	if currMoveAnalysis.IsMate {
		analyzedMove.EvalGraph = float64(currMoveAnalysis.Eval)
	} else {
		analyzedMove.EvalGraph = float64(currMoveAnalysis.Eval) / 100.0
	}

	return analyzedMove, nil
}

func (a *AnalysisService) getMoveAnalysis(moveOrder int, gameID string) (models.StockfishAnalysisResult, error) {
	move, err := a.analysisRepo.GetMoveByOrder(moveOrder, gameID)
	if err != nil {
		err = fmt.Errorf("AnalysisService-GetAnalyzedMoveByOrder-GetMoveByOrder: %w", err)
		return models.StockfishAnalysisResult{}, err
	}

	stockfishResult, err := a.StockfishAnalyze(move.Fen, "hard")
	if err != nil {
		err = fmt.Errorf("AnalysisService-GetAnalyzedMoveByOrder-StockfishAnalyze: %w", err)
		return models.StockfishAnalysisResult{}, err
	}

	return stockfishResult, nil
}

func (a *AnalysisService) classifyMove(
	evalBefore *int,
	evalAfter *int,
	mateBefore *int,
	mateAfter *int,
) string {
	if mateBefore == nil && mateAfter == nil {
		if evalBefore == nil || evalAfter == nil {
			return "Unknown"
		}

		diff := *evalAfter - *evalBefore

		if diff <= 0 {
			return "Excellent"
		}

		switch {
		case diff <= 20:
			return "Excellent"
		case diff <= 50:
			return "Good"
		case diff <= 150:
			return "Inaccuracy"
		case diff <= 300:
			return "Mistake"
		default:
			return "Blunder"
		}
	}

	if mateBefore != nil && mateAfter != nil {
		if (*mateBefore > 0 && *mateAfter > 0) || (*mateBefore < 0 && *mateAfter < 0) {
			if math.Abs(float64(*mateAfter)) <= math.Abs(float64(*mateBefore)) {
				return "Excellent"
			}
			return "Inaccuracy"
		}
		return "Blunder"
	}

	// no longer has checkmates
	if mateBefore != nil && mateAfter == nil {
		if *mateBefore > 0 {
			return "Blunder"
		}
		return "Excellent"
	}

	// got checkmates
	if mateBefore == nil && mateAfter != nil {
		if *mateAfter > 0 {
			return "Excellent"
		}
		return "Blunder"
	}

	return "Unknown"
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
