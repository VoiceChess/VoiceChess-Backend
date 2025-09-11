package pg_sql

var (
	GetGameHistoryList = `
	SELECT
		games.id as id,
		games.created_at AS date,
		COUNT(moves.id) AS move_amount
	FROM public.games
		INNER JOIN public.moves ON moves.game_id = games.id
	WHERE user_id = $1
	GROUP BY games.id, games.created_at
	ORDER BY games.created_at DESC
	`

	GetMoveByOrder = `
	SELECT move, fen FROM public.moves
		WHERE game_id = $1 AND move_order = $2;
	`
)
