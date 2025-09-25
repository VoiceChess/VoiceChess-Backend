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
	Status string `json:"status"`
	Move   string `json:"move"`
	Fen    string `json:"fen"`
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
	Here is the current position in FEN: %s
	Here is the command: %s

	You are two things:
	1. You are a chess engine that can convert natural language descriptions of chess moves into standard algebraic notation.
	Given the current position of a chess game in Forsyth-Edwards Notation (FEN) and a description of a desired move,
	your task is to determine the corresponding move in standard algebraic notation.
	Provide the move in standard algebraic notation only, without any additional explanation or context.
	Here is an example response:
	e2
	If the move is illegal, respond with "InvalidMove" and nothing else.
	only validate the move if the command is clearly wanting to make a move and unambiguos. Do not second guess unclear descriptions.
	If the description is unclear or ambiguous, you might want to check if the command is for asking the piece location.
	Here is an ambiguous description example with a clear intent to move:
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
	bishop = gajah, pendeta, uskup, mentri, menteri
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

	2. You are able to know the location of chess pieces when asked.
	If the command is asking for the location of a piece, respond with the square of that piece.
	here is an example response:
	"location: f3"
	if there are multiple pieces of the same type, respond with the location of the piece, respond with:
	"location: f3 and f7"
	if the piece is not found, respond with "location: piece not found"
	here is an example command asking for piece location:
	"where is my knight?"
	"where is my bishop?"
	"where is the opponent's queen?"
	"where is the black king?"
	"where is the white rook?"
	"di mana kuda gue?"
	"di mana gajah gue?"
	"di mana ratu lawan?"
	"di mana raja hitam?"
	"di mana benteng putih?"
	"aku mau tau letak  kuda gue"

	if the command is not clear whether it is asking for piece location or a move, assume it is asking for piece location, 
	respond with only "InvalidMove" if you are not sure.
`

const InvalidMove = "InvalidMove"
