package services

import (
	"bytes"
	"fmt"
	"log"
	"math"
	"strings"
	"text/template"

	"samsungvoicebe/helper"
	"samsungvoicebe/models"
	"samsungvoicebe/repo"

	"github.com/notnil/chess"
	"github.com/notnil/chess/uci"
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

	var prevMove models.Move
	var prevMoveBestMove string

	if moveOrder > 2 || firstMove.Move != "black" {
		prevMove, err = a.analysisRepo.GetMoveByOrder(moveOrder-1, gameID)
		if err != nil {
			err = fmt.Errorf("AnalysisService-GetAnalyzedMoveByOrder-GetMoveByOrder: %w", err)
			return models.MoveAnalysis{}, err
		}

		prevMoveAnalysis, err := a.getMoveAnalysis(moveOrder-1, gameID)
		if err != nil {
			err = fmt.Errorf("AnalysisService-GetAnalyzedMoveByOrder-getMoveAnalysis: %w", err)
			return models.MoveAnalysis{}, err
		}

		prevMoveBestMove = prevMoveAnalysis.BestMove

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
		prevMove, err = a.analysisRepo.GetMoveByOrder(moveOrder, gameID)
		if err != nil {
			err = fmt.Errorf("AnalysisService-GetAnalyzedMoveByOrder-GetMoveByOrder: %w", err)
			return models.MoveAnalysis{}, err
		}

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

	analyzedMove.OverviewSection, err = a.GetOverviewSectionData(
		prevMove.Fen,
		currMove.Fen,
		currMove.Move,
		analyzedMove.MoveGrade,
		analyzedMove.EvalGraph,
		currMoveAnalysis.IsMate,
	)
	if err != nil {
		return models.MoveAnalysis{}, err
	}

	analyzedMove.ThreatsSection, err = a.GetThreatsSectionData(currMove.Fen)
	if err != nil {
		return models.MoveAnalysis{}, err
	}

	analyzedMove.BestMoveSection, err = a.GetBestMoveSectionData(prevMove.Fen, prevMoveBestMove)
	if err != nil {
		return models.MoveAnalysis{}, err
	}

	analyzedMove.StrategySections, err = a.GetStrategySectionData(currMove.Fen)
	if err != nil {
		return models.MoveAnalysis{}, err
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

func (a *AnalysisService) GetOverviewSectionData(
	startingFEN,
	resultingFEN,
	playerMove,
	moveGrade string,
	evalGraph float64,
	IsEvalMate bool,
) (models.OverviewSection, error) {
	promptData := models.GetGradeExplanationData{
		StartFEN:     startingFEN,
		PlayerMove:   playerMove,
		ResultingFEN: resultingFEN,
		EvalGraph:    evalGraph,
		IsEvalMate:   IsEvalMate,
		MoveGrade:    moveGrade,
	}

	prompt, err := a.buildGradeExplanationPrompt(promptData)
	if err != nil {
		log.Fatalf("AnalysisService-GetOverviewSectionData-buildGradeExplanationPrompt: %v", err)
		return models.OverviewSection{}, err
	}
	gradeExplanation := helper.PromptGemini(prompt)

	prompt = fmt.Sprintf(models.GetThreatPrompt, resultingFEN)
	threats := helper.PromptGemini(prompt)

	var colorWithAdvantage string
	stringifiedEvalGraph := fmt.Sprintf("%.2f", evalGraph)
	if evalGraph > 0 && !IsEvalMate {
		colorWithAdvantage = "Advantage: Winning for White " + stringifiedEvalGraph
	} else if evalGraph < 0 && !IsEvalMate {
		colorWithAdvantage = "Advantage: Winning for Black " + stringifiedEvalGraph
	} else if evalGraph > 0 && IsEvalMate {
		colorWithAdvantage = "Advantage: Winning for White, mate in " + stringifiedEvalGraph
	} else if evalGraph < 0 && IsEvalMate {
		colorWithAdvantage = "Advantage: Winning for Black, mate in " + stringifiedEvalGraph
	} else {
		colorWithAdvantage = "No Advantage: Equal Position"
	}

	prompt = fmt.Sprintf(models.GetAdvantageExplanationPrompt, stringifiedEvalGraph, resultingFEN, colorWithAdvantage)
	advantageExplanation := helper.PromptGemini(prompt)

	prompt = fmt.Sprintf(models.GetNotableMovesPrompt, resultingFEN)
	notableMove := helper.PromptGemini(prompt)

	overviewData := models.OverviewSection{
		GradeExplanation:     gradeExplanation,
		Threats:              threats,
		AdvantageExplanation: advantageExplanation,
		NotableMoves:         notableMove,
		ColorWithAdvantage:   colorWithAdvantage,
	}

	return overviewData, nil
}

func (a *AnalysisService) GetThreatsSectionData(FEN string) (models.ThreatsSection, error) {

	FENAnalysis, err := a.StockfishAnalyze(FEN, "hard")
	if err != nil {
		err = fmt.Errorf("AnalysisService-GetThreatsSectionData-StockfishAnalyze: %w", err)
		return models.ThreatsSection{}, err
	}

	threateningMove := FENAnalysis.BestMove

	prompt := fmt.Sprintf(models.GetThreateningMoveExplanationPrompt, FEN, threateningMove)
	explanation := helper.PromptGemini(prompt)

	threatsData := models.ThreatsSection{
		ThreateningMove: threateningMove,
		Explanation:     explanation,
	}

	return threatsData, nil
}

func (a *AnalysisService) GetBestMoveSectionData(FEN, bestMove string) (models.BestMoveSection, error) {

	prompt := fmt.Sprintf(models.GetBestMoveExplanationPrompt, FEN, bestMove)
	explanation := helper.PromptGemini(prompt)

	bestMoveData := models.BestMoveSection{
		BestMove:    bestMove,
		Explanation: explanation,
	}

	return bestMoveData, nil
}

func (a *AnalysisService) GetStrategySectionData(FEN string) (models.StrategySection, error) {

	prompt := fmt.Sprintf(models.GetStrategyTitlePrompt, FEN)
	title := helper.PromptGemini(prompt)

	prompt = fmt.Sprintf(models.GetStrategyExplanationPrompt, FEN, title)
	explanation := helper.PromptGemini(prompt)

	strategyData := models.StrategySection{
		Title:       title,
		Explanation: explanation,
	}

	return strategyData, nil
}

func (a *AnalysisService) buildGradeExplanationPrompt(data models.GetGradeExplanationData) (string, error) {
	tmpl, err := template.New("gradeExplanation").Parse(models.GetGradeExplanationPrompt)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		log.Fatal(err)
	}

	return buf.String(), nil
}
