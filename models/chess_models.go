package models

type ChessRequest struct {
	Message string `json:"message"`
	Fen     string `json:"fen" binding:"required"`
	Type    string `json:"type" binding:"required,oneof=black white"`
	Mode    string `json:"mode" binding:"required,oneof=easy medium hard"`
}

type ChessResponse struct {
	Move      string `json:"move,omitempty"`
	NewFen    string `json:"new_fen,omitempty"`
	Status    string `json:"status,omitempty"`
	Error     string `json:"error,omitempty"`
	IsGameEnd bool   `json:"is_game_end,omitempty"`
	Winner    string `json:"winner,omitempty"`
}

type MoveIntentAnalysis struct {
	IsValidRequest bool   `json:"is_valid_request"`
	FromSquare     string `json:"from_square"`
	ToSquare       string `json:"to_square"`
	MoveNotation   string `json:"move_notation"`
	Confidence     int    `json:"confidence"`
	Explanation    string `json:"explanation"`
}

type AIStrategy struct {
	Depth          int
	RandomFactor   float64
	PreferCaptures bool
	PreferCenter   bool
	AvoidBlunders  bool
}

func GetAIStrategy(mode string) AIStrategy {
	switch mode {
	case "easy":
		return AIStrategy{
			Depth:          1,
			RandomFactor:   0.4,
			PreferCaptures: false,
			PreferCenter:   false,
			AvoidBlunders:  false,
		}
	case "medium":
		return AIStrategy{
			Depth:          2,
			RandomFactor:   0.2,
			PreferCaptures: true,
			PreferCenter:   true,
			AvoidBlunders:  true,
		}
	case "hard":
		return AIStrategy{
			Depth:          3,
			RandomFactor:   0.05,
			PreferCaptures: true,
			PreferCenter:   true,
			AvoidBlunders:  true,
		}
	default:
		return GetAIStrategy("medium")
	}
}
