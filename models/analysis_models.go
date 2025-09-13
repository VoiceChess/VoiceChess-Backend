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
	example:
	if you don't know the active color, you can fill it with "w" or "b".'
	if you don't know the castling availability, you can fill it with "KQkq" or "-"
	if you don't know the en passant target square, you can fill it with "-" or "e3"
	if you don't know the halfmove clock, you can fill it with "0"
	if you don't know the fullmove number, you can fill it with "1"
	Your task is to analyze the image and extract the piece placement part of the FEN.
	Then, construct a valid FEN string by appending the other necessary parts with any valid values.
	Finally, respond with the complete FEN string.

	given the following picture, extract the FEN string that represents the piece placement on the chessboard.

	If the image is not an image of a chessboard or if you cannot determine the FEN, respond with "InvalidImage" and nothing else.
	Respond with only the FEN string and nothing else.

	Make sure the FEN is valid and correctly formatted.
	this is an example of a valid FEN: 
	rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1

	these are some examples of invalid FENs:
	rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0
	r2k1r2/7r/4q2b/p6p/8/2R5/4PPPP/ w - - 0 1
	r2k1r2/7r/4q2b/p6p/8/2R5/4PPPP/
	rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBN w KQkq - 0 1
	rnbqkbnrr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1
	rnbqkbnx/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1
	rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNRK w KQkq - 0 1
	rnbqkbnr/pppppppp/8/8/8/8/PPPPPPP1/P3K2R w KQkq - 0 1
	rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR x KQkq - 0 1

	if the FEN you generate is invalid, respond with "InvalidImage" and nothing else.
	
`
const InvalidImage = "InvalidImage"
