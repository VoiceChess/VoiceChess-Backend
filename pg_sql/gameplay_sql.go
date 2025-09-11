package pg_sql

var (
	Move = `
	INSERT INTO public.moves (game_id, fen, move)
		VALUES ($1, $2, $3);
	`

	CreateGame = `
	INSERT INTO public.games (user_id) VALUES ($1) RETURNING id;
	`
)
