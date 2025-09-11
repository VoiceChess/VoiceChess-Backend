package models

type BotMove struct {
	Move string `json:"bot_move"`
	Fen  string `json:"fen"`
}

type PlayerMoveRequest struct {
	Move     string `json:"move" binding:"required"`
	Fen      string `json:"fen" binding:"required"`
	BotLevel string `json:"bot_level" binding:"required"`
}
