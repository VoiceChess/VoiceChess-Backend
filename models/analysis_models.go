package models

type Move struct {
	Move string `db:"move"`
	Fen  string `db:"fen"`
}

type MoveAnalysis struct {
	Move     string `json:"move"`
	Fen      string `json:"fen"`
	BestMove string `json:"best_move"`
}

type Game struct {
	GameID     string `db:"game_id"`
	Date       string `db:"date"`
	MoveAmount int    `db:"move_amount"`
}

type StockfishAnalysisResult struct {
	Fen      string
	BestMove string
}

const (
	BotLevelEasy   = 2
	BotLevelMedium = 5
	BotLevelhard   = 10
)

const GetFenFromPicturePrompt = `
	You are an OCR (Optical Character Recognition) expert. Your task is to extract the Forsyth-Edwards Notation (FEN)
	from the given picture. The FEN is a standard notation for describing 
	a particular board position of a chess game.
	I noticed that FEN contains other informations like active color, castling availability, en passant target square, halfmove clock, and fullmove number, etc.
	However, for this task, we are only interested in the piece placement part of the FEN.
	Other informations can be filled with anything that is valid for FEN.
	given the following picture, extract the FEN string that represents the piece placement on the chessboard.
	Respond with only the FEN string and nothing else.
	If the image is not an image of a chessboard or if you cannot determine the FEN, respond with "InvalidImage" and nothing else.
`
const InvalidImage = "InvalidImage"
