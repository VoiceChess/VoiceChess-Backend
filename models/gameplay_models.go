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

type HintRequest struct {
	Fen string `json:"fen" binding:"required"`
}

type PlayerMoveByTranscriptionRequest struct {
	Fen           string `json:"fen" binding:"required"`
	Transcription string `json:"transcription" binding:"required"`
}

type PlayerMoveByTranscription struct {
	Move string `json:"move"`
	Fen  string `json:"fen"`
}

const HintPrompt = `
	You are a chess tutor. Your task is to provide helpful hints to the player based on the current
	position of the chess game. Analyze the position and suggest a move that would improve the player's chances of winning.
	We provide the current position in Forsyth-Edwards Notation (FEN).
	Here is the current position in FEN: %s
    Do not blantantly give away the best move, but guide the player towards it.
	Make sure your hint is clear and concise, focusing on strategic elements of the game.
	Do not mention the FEN or any specific moves in your hint.
	Do not re-explain what the FEN means. Straight to the what the hint is. 
	Make a medium length response, around 1-3 sentences.
	Here is an example response:

	There are pieces that can be developed to control the center of the board
	Consider moving your knight to a position where it can attack the opponent's queen.
`

const MoveFromDescriptionPrompt = `
	You are a chess engine that can convert natural language descriptions of chess moves into standard algebraic notation.
	Given the current position of a chess game in Forsyth-Edwards Notation (FEN) and a description of a desired move,
	your task is to determine the corresponding move in standard algebraic notation.
	Here is the current position in FEN: %s
	Here is the description of the desired move: %s
	Provide the move in standard algebraic notation only, without any additional explanation or context.
	Here is an example response:
	e2
	If you cannot determine the move, respond with "InvalidMove" and nothing else.
	If the move is illegal, respond with "InvalidMove" and nothing else.
	only validate the move if the description is clear and unambiguos. Do not second guess unclear descriptions.
	If the description is unclear or ambiguous, respond with "InvalidMove" and nothing else.
	Here is an ambiguous description example:
	"move the piece in front of the king"
	"pawn go forward"

	Here is a clear description example:
	"move the pawn in front of the king to e4"
	"move the knight to f3"
	"move the bishop to c4"
	"castle kingside"
	"castle queenside"
	"move the queen to h5"
	"i think i will move my pawn to d5"
	"pawn form e2 to e4 looks promising"
	
	input might be in indonesian, so be aware of that.
	king = raja
	queen = ratu, permaisuri, putri
	bishop = gajah, pendeta, uskup
	knight = kuda, kesatria, kavaleri
	rook = benteng, menara, kapal
	pawn = pion, bidak, prajurit, serdadu, tentara
	
	Here is an indonesian description example:
	"pion di depan raja ke e4"
	"pindahkan kuda ke f3"
	"pindahkan gajah ke c4"
	"rokade ke sisi raja"
	"rokade ke sisi ratu"
	"pindahkan ratu ke h5"
	"saya rasa saya akan memindahkan pion saya ke d5"
	"pion dari e2 ke e4 terlihat menjanjikan"
	"gw mau majuin pion ke e4"
	"pion gue ke e4 aja"
	"pion di e2 ke e4"
	"gw mau pion gw ke e4"
`

const InvalidMove = "InvalidMove"
