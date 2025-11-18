package models

type Move struct {
	Move string `db:"move"`
	Fen  string `db:"fen"`
}

type MoveAnalysis struct {
	Move             string          `json:"move"`
	Fen              string          `json:"fen"`
	BestMove         string          `json:"best_move"`
	EvalGraph        float64         `json:"eval_graph"`
	IsEvalMate       bool            `json:"is_eval_mate"`
	MoveGrade        string          `json:"move_grade"`
	OverviewSection  OverviewSection `json:"overview_section"`
	ThreatsSection   ThreatsSection  `json:"threats_section"`
	BestMoveSection  BestMoveSection `json:"best_move_section"`
	StrategySections StrategySection `json:"strategy_section"`
}

type OverviewSection struct {
	GradeExplanation     string `json:"grade_explanation"`
	Threats              string `json:"threats"`
	ColorWithAdvantage   string `json:"color_with_advantage"`
	AdvantageExplanation string `json:"advantage_explanation"`
	NotableMoves         string `json:"notable_moves"`
}

type ThreatsSection struct {
	ThreateningMove string `json:"threatening_move"`
	Explanation     string `json:"explanation"`
}

type BestMoveSection struct {
	BestMove    string `json:"best_move"`
	Explanation string `json:"explanation"`
}

type StrategySection struct {
	Title       string `json:"title"`
	Explanation string `json:"explanation"`
}

type Game struct {
	GameID     string `db:"game_id"`
	Date       string `db:"date"`
	MoveAmount int    `db:"move_amount"`
}

type StockfishAnalysisResult struct {
	Fen      string
	BestMove string
	Eval     int
	IsMate   bool
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

const GeminiModel = "gemini-2.0-flash"

type GetGradeExplanationData struct {
	StartFEN     string
	PlayerColor  string
	PlayerMove   string
	ResultingFEN string
	EvalGraph    float64
	IsEvalMate   bool
	MoveGrade    string
}

const GetGradeExplanationPrompt = `
	from this given FEN: {{.StartFEN}} ,
	someone made this move: {{.PlayerMove}} and now the FEN is: {{.ResultingFEN}}.
	the stockfish eval graph shows this score: {{.EvalGraph}}, and is mate: {{.IsEvalMate}}.
	provide me an explanation of why this move is graded as {{.MoveGrade}} for {{.PlayerColor}}. make it short and concise text only without any format to the point and dont return the fen
`

const GetThreatPrompt = `
	from this given FEN: %s, what are the threats for the next player? make it short and concise text only without any format to the point and dont return the fen
`

const GetAdvantageExplanationPrompt = `
	from the eval graph score of %s and this FEN %s explain why is this the case %s make it short and concise text only without any format to the point and dont return the fen
`

const GetNotableMovesPrompt = `
	from this given FEN: %s, what would be the notable move make it short and concise text only without any format to the point and dont return the fen
`

const GetThreateningMoveExplanationPrompt = `
	from this given FEN: %s, explain the threatening move: %s make it short and concise text only without any format to the point and dont return the fen
`

const GetBestMoveExplanationPrompt = `
	from this given FEN: %s, explain  why the best move is: %s . make it short and concise text only without any format to the point and dont return the fen
`

const GetStrategyTitlePrompt = `
	from this FEN %s, provide a short title for the best strategy they can use in this position.
`

const GetStrategyExplanationPrompt = `
	from this FEN %s and this title %s, provide a short and concise text only without any format explanation of the strategy.
`
