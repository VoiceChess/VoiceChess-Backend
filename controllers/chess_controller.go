package controllers

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"samsungvoicebe/config"
	"samsungvoicebe/helper"
	"samsungvoicebe/models"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/notnil/chess"
)

type ChessController struct {
	config *config.Config
}

func NewChessController(cfg *config.Config) *ChessController {
	rand.Seed(time.Now().UnixNano())
	return &ChessController{
		config: cfg,
	}
}

func (cc *ChessController) PlayChess(c *gin.Context) {
	var req models.ChessRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ChessResponse{
			Error: "Invalid request format: " + err.Error(),
		})
		return
	}

	fenNotation, err := chess.FEN(req.Fen)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ChessResponse{
			Error: "Invalid FEN notation: " + err.Error(),
		})
		return
	}

	game := chess.NewGame(fenNotation)

	if game.Outcome() != chess.NoOutcome {
		c.JSON(http.StatusBadRequest, models.ChessResponse{
			Error:     "Game is already finished",
			Status:    cc.getGameStatus(game),
			Winner:    cc.getWinner(game),
			IsGameEnd: true,
		})
		return
	}

	if strings.TrimSpace(req.Message) == "" {
		cc.handleAIMove(c, game, req.Mode)
		return
	}

	cc.handlePlayerMove(c, req, game)
}

func (cc *ChessController) handleAIMove(c *gin.Context, game *chess.Game, mode string) {
	validMoves := game.ValidMoves()
	if len(validMoves) == 0 {
		c.JSON(http.StatusOK, models.ChessResponse{
			Status:    cc.getGameStatus(game),
			Winner:    cc.getWinner(game),
			IsGameEnd: true,
		})
		return
	}

	strategy := models.GetAIStrategy(mode)
	bestMove := cc.selectBestMoveWithStrategy(validMoves, game, strategy)

	err := game.Move(bestMove)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ChessResponse{
			Error: "Failed to apply AI move: " + err.Error(),
		})
		return
	}

	newFen := game.FEN()
	status := cc.getGameStatus(game)
	winner := cc.getWinner(game)

	c.JSON(http.StatusOK, models.ChessResponse{
		Move:      bestMove.String(),
		NewFen:    newFen,
		Status:    status,
		Winner:    winner,
		IsGameEnd: game.Outcome() != chess.NoOutcome,
	})
}

func (cc *ChessController) handlePlayerMove(c *gin.Context, req models.ChessRequest, game *chess.Game) {
	analysis, err := cc.analyzeMoveTWithGemini(req.Message, req.Fen)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ChessResponse{
			Error: "Failed to analyze move: " + err.Error(),
		})
		return
	}

	if !analysis.IsValidRequest {
		c.JSON(http.StatusBadRequest, models.ChessResponse{
			Error: "Invalid move request: " + analysis.Explanation,
		})
		return
	}

	move, err := cc.findMoveFromAnalysis(analysis, game)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ChessResponse{
			Error: "Invalid move: " + err.Error(),
		})
		return
	}

	err = game.Move(move)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ChessResponse{
			Error: "Illegal move: " + err.Error(),
		})
		return
	}

	newFen := game.FEN()
	status := cc.getGameStatus(game)
	winner := cc.getWinner(game)

	c.JSON(http.StatusOK, models.ChessResponse{
		Move:      move.String(),
		NewFen:    newFen,
		Status:    status,
		Winner:    winner,
		IsGameEnd: game.Outcome() != chess.NoOutcome,
	})
}

func (cc *ChessController) selectBestMoveWithStrategy(moves []*chess.Move, game *chess.Game, strategy models.AIStrategy) *chess.Move {
	if len(moves) == 0 {
		return nil
	}

	if strategy.RandomFactor > 0.3 && rand.Float64() < strategy.RandomFactor {
		return moves[rand.Intn(len(moves))]
	}

	type ScoredMove struct {
		move  *chess.Move
		score int
	}

	var scoredMoves []ScoredMove

	for _, move := range moves {
		score := cc.scoreMove(move, game, strategy)
		scoredMoves = append(scoredMoves, ScoredMove{move: move, score: score})
	}

	bestScore := scoredMoves[0].score
	for _, sm := range scoredMoves {
		if sm.score > bestScore {
			bestScore = sm.score
		}
	}

	var bestMoves []*chess.Move
	for _, sm := range scoredMoves {
		if sm.score == bestScore {
			bestMoves = append(bestMoves, sm.move)
		}
	}

	if len(bestMoves) > 1 && rand.Float64() < strategy.RandomFactor {
		return bestMoves[rand.Intn(len(bestMoves))]
	}

	return bestMoves[0]
}

func (cc *ChessController) scoreMove(move *chess.Move, game *chess.Game, strategy models.AIStrategy) int {
	score := 0

	if strategy.PreferCaptures && move.HasTag(chess.Capture) {
		score += 50

		if strings.Contains(move.String(), "Q") {
			score += 90
		} else if strings.Contains(move.String(), "R") {
			score += 50
		} else if strings.Contains(move.String(), "B") || strings.Contains(move.String(), "N") {
			score += 30
		} else {
			score += 10
		}
	}

	if strategy.PreferCenter && cc.isCenterSquare(move.S2()) {
		score += 20
	}

	if move.HasTag(chess.Check) {
		score += 30
	}

	outcome := game.Outcome()
	if outcome == chess.WhiteWon || outcome == chess.BlackWon {
		score += 1000
	}

	if strategy.AvoidBlunders {
		piece := game.Position().Board().Piece(move.S1())
		if piece.Type() == chess.Queen && cc.isSquareUnderAttack(move.S2(), game) {
			score -= 80
		} else if piece.Type() == chess.Rook && cc.isSquareUnderAttack(move.S2(), game) {
			score -= 40
		}
	}

	if strategy.RandomFactor > 0 {
		randomBonus := int(float64(rand.Intn(20)) * strategy.RandomFactor)
		score += randomBonus
	}

	if move.HasTag(chess.KingSideCastle) || move.HasTag(chess.QueenSideCastle) {
		score += 25
	}

	if move.Promo() != chess.NoPieceType {
		score += 80
	}

	return score
}

func (cc *ChessController) isCenterSquare(square chess.Square) bool {
	centerSquares := []chess.Square{
		chess.D4, chess.D5, chess.E4, chess.E5,
	}

	extendedCenter := []chess.Square{
		chess.C3, chess.C4, chess.C5, chess.C6,
		chess.D3, chess.D6,
		chess.E3, chess.E6,
		chess.F3, chess.F4, chess.F5, chess.F6,
	}

	for _, center := range centerSquares {
		if square == center {
			return true
		}
	}

	for _, extended := range extendedCenter {
		if square == extended {
			return true
		}
	}

	return false
}

func (cc *ChessController) isSquareUnderAttack(square chess.Square, game *chess.Game) bool {
	board := game.Position().Board()
	turn := game.Position().Turn()

	for sq := chess.A1; sq <= chess.H8; sq++ {
		piece := board.Piece(sq)

		if piece != chess.NoPiece && piece.Color() != turn {
			moves := game.ValidMoves()
			for _, move := range moves {
				if move.S1() == sq && move.S2() == square {
					return true
				}
			}
		}
	}

	return false
}

func (cc *ChessController) findMoveFromAnalysis(analysis *models.GeminiMoveAnalysis, game *chess.Game) (*chess.Move, error) {
	validMoves := game.ValidMoves()

	if analysis.FromSquare != "" && analysis.ToSquare != "" {
		fromSquare := strings.ToLower(analysis.FromSquare)
		toSquare := strings.ToLower(analysis.ToSquare)

		for _, move := range validMoves {
			if strings.ToLower(move.S1().String()) == fromSquare && strings.ToLower(move.S2().String()) == toSquare {
				return move, nil
			}
		}
	}

	if analysis.MoveNotation != "" {
		notation := strings.ToLower(strings.TrimSpace(analysis.MoveNotation))
		for _, move := range validMoves {
			moveStr := strings.ToLower(move.String())
			if moveStr == notation || moveStr == strings.ReplaceAll(notation, "x", "") {
				return move, nil
			}
		}
	}

	if analysis.ToSquare != "" {
		toSquare := strings.ToLower(analysis.ToSquare)
		var possibleMoves []*chess.Move
		for _, move := range validMoves {
			if strings.ToLower(move.S2().String()) == toSquare {
				possibleMoves = append(possibleMoves, move)
			}
		}

		if len(possibleMoves) == 1 {
			return possibleMoves[0], nil
		} else if len(possibleMoves) > 1 {
			if analysis.FromSquare != "" {
				fromSquare := strings.ToLower(analysis.FromSquare)
				for _, move := range possibleMoves {
					if strings.ToLower(move.S1().String()) == fromSquare {
						return move, nil
					}
				}
			}
			return nil, fmt.Errorf("ambiguous move to %s, please specify which piece", analysis.ToSquare)
		}
	}

	return nil, fmt.Errorf("move not found or illegal")
}

func (cc *ChessController) getGameStatus(game *chess.Game) string {
	outcome := game.Outcome()
	method := game.Method()

	switch outcome {
	case chess.WhiteWon:
		if method == chess.Checkmate {
			return "checkmate"
		}
		return "white_wins"
	case chess.BlackWon:
		if method == chess.Checkmate {
			return "checkmate"
		}
		return "black_wins"
	case chess.Draw:
		return "draw"
	default:
		validMoves := game.ValidMoves()
		if len(validMoves) > 0 {
			board := game.Position().Board()
			kingSquare := chess.NoSquare
			turn := game.Position().Turn()

			for sq := chess.A1; sq <= chess.H8; sq++ {
				piece := board.Piece(sq)
				if piece.Type() == chess.King && piece.Color() == turn {
					kingSquare = sq
					break
				}
			}

			if kingSquare != chess.NoSquare && cc.isSquareUnderAttack(kingSquare, game) {
				return "check"
			}
		}
		return "ongoing"
	}
}

func (cc *ChessController) getWinner(game *chess.Game) string {
	switch game.Outcome() {
	case chess.WhiteWon:
		return "white"
	case chess.BlackWon:
		return "black"
	case chess.Draw:
		return "draw"
	default:
		return ""
	}
}

func (cc *ChessController) analyzeMoveTWithGemini(message, fen string) (*models.GeminiMoveAnalysis, error) {
	prompt := fmt.Sprintf(`
Analyze this chess move request and extract the move information.

Current board position (FEN): %s
Player message: "%s"

Task: Determine if this is a valid move request and extract move details.

Chess squares are labeled a1-h8 (files a-h, ranks 1-8).

IMPORTANT: Support both English and Indonesian language:

English examples:
- "pawn from e2 to e4" → from: "e2", to: "e4", notation: "e2e4"
- "knight to f3" → to: "f3", notation: "nf3"
- "queen takes d5" → to: "d5", notation: "qxd5"
- "castle kingside" → notation: "o-o"
- "move rook to a1" → to: "a1"

Indonesian examples:
- "b1 ke c3" → from: "b1", to: "c3", notation: "b1c3"
- "pawn dari e2 ke e4" → from: "e2", to: "e4", notation: "e2e4"  
- "kuda ke f3" → to: "f3", notation: "nf3"
- "benteng ke a1" → to: "a1"
- "ratu ambil d5" → to: "d5", notation: "qxd5"
- "rokade pendek" → notation: "o-o"
- "rokade panjang" → notation: "o-o-o"

Key Indonesian chess terms:
- "ke" = "to"
- "dari" = "from"
- "ambil" = "takes"
- "kuda" = "knight"
- "benteng" = "rook"
- "gajah" = "bishop"
- "ratu" = "queen"
- "raja" = "king"
- "pion/bidak" = "pawn"

Pattern recognition:
- "[square] ke [square]" = move from first square to second square
- "kuda ke [square]" = knight to square
- Any format like "a1 ke b2" should be interpreted as move from a1 to b2

Examples of invalid requests:
- "hello" → not a chess move
- "good game" → not a chess move
- "how are you" → not a chess move

Respond with ONLY a JSON object in this exact format:
{
  "is_valid_request": true/false,
  "from_square": "b1",
  "to_square": "c3", 
  "move_notation": "b1c3",
  "confidence": 9,
  "explanation": "User wants to move piece from b1 to c3"
}

If the message is clearly not about making a chess move, set is_valid_request to false and explain why in explanation.
`, fen, message)

	response, err := helper.PromptOllama(prompt)
	if err != nil {
		return nil, err
	}

	response = strings.TrimSpace(response)

	if strings.HasPrefix(response, "```json") {
		response = strings.TrimPrefix(response, "```json")
		response = strings.TrimSuffix(response, "```")
	} else if strings.HasPrefix(response, "```") {
		response = strings.TrimPrefix(response, "```")
		response = strings.TrimSuffix(response, "```")
	}

	response = strings.TrimSpace(response)

	var analysis models.GeminiMoveAnalysis
	if err := json.Unmarshal([]byte(response), &analysis); err != nil {
		return nil, fmt.Errorf("failed to parse Gemini response: %v", err)
	}

	return &analysis, nil
}
