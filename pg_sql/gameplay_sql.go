package pg_sql

var (
	Move = `
	INSERT INTO public.moves (game_id, fen, move)
		VALUES ($1, $2, $3);
	`

	CreateGame = `
	INSERT INTO public.games (user_id) VALUES ($1) RETURNING id;
	`

	GameBelongsToUser = `
	SELECT EXISTS(SELECT 1 FROM public.games WHERE id = $1 AND user_id = $2);
	`

	DeleteLatestMove = `
	WITH last_move AS (
		SELECT id
		FROM public.moves
		WHERE game_id = $1
		ORDER BY created_at DESC
		LIMIT 1
	)
	DELETE FROM moves
	WHERE id = (SELECT id FROM last_move);
	`
)
