package db

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
)

func New() (*sql.DB, error) {
	databaseUrl := os.Getenv("POSTGRESQL_URL")
	if databaseUrl == "" {
		return nil, fmt.Errorf("POSTGRESQL_URL environment variable not set")
	}

	db, err := sql.Open("postgres", databaseUrl)
	if err != nil {
		return nil, fmt.Errorf("failed to open DB: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping DB: %w", err)
	}

	if os.Getenv("AUTO_MIGRATE") != "false" {
		if err := applyMigrations(db); err != nil {
			return nil, fmt.Errorf("failed to apply database schema: %w", err)
		}
	}

	return db, nil
}

func applyMigrations(db *sql.DB) error {
	_, err := db.Exec(schemaSQL)
	return err
}

const schemaSQL = `
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

ALTER TABLE IF EXISTS public.moves DROP CONSTRAINT IF EXISTS moves_game_id_fkey;
ALTER TABLE IF EXISTS public.moves DROP CONSTRAINT IF EXISTS moves_user_id_fkey;
ALTER TABLE IF EXISTS public.games DROP CONSTRAINT IF EXISTS games_user_id_fkey;

CREATE TABLE IF NOT EXISTS public.users (
	id TEXT PRIMARY KEY,
	created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
	updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

ALTER TABLE public.users ALTER COLUMN id TYPE TEXT USING id::TEXT;
ALTER TABLE public.users ADD COLUMN IF NOT EXISTS created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP;
ALTER TABLE public.users ADD COLUMN IF NOT EXISTS updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP;
DO $$
BEGIN
	IF EXISTS (SELECT 1 FROM information_schema.columns WHERE table_schema = 'public' AND table_name = 'users' AND column_name = 'username') THEN
		ALTER TABLE public.users ALTER COLUMN username DROP NOT NULL;
	END IF;
	IF EXISTS (SELECT 1 FROM information_schema.columns WHERE table_schema = 'public' AND table_name = 'users' AND column_name = 'email') THEN
		ALTER TABLE public.users ALTER COLUMN email DROP NOT NULL;
	END IF;
	IF EXISTS (SELECT 1 FROM information_schema.columns WHERE table_schema = 'public' AND table_name = 'users' AND column_name = 'password_hash') THEN
		ALTER TABLE public.users ALTER COLUMN password_hash DROP NOT NULL;
	END IF;
END $$;

CREATE TABLE IF NOT EXISTS public.games (
	id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
	user_id TEXT,
	created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
	updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

ALTER TABLE public.games ADD COLUMN IF NOT EXISTS user_id TEXT;
ALTER TABLE public.games ALTER COLUMN user_id TYPE TEXT USING user_id::TEXT;
ALTER TABLE public.games ADD COLUMN IF NOT EXISTS created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP;
ALTER TABLE public.games ADD COLUMN IF NOT EXISTS updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP;
ALTER TABLE public.games DROP COLUMN IF EXISTS fen;
ALTER TABLE public.games DROP COLUMN IF EXISTS result;
ALTER TABLE public.games DROP COLUMN IF EXISTS end_type;
ALTER TABLE public.games DROP COLUMN IF EXISTS move_amount;

CREATE TABLE IF NOT EXISTS public.moves (
	id UUID DEFAULT uuid_generate_v4() PRIMARY KEY NOT NULL,
	game_id UUID,
	move_order INT,
	move VARCHAR(10),
	fen VARCHAR(500),
	created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
	updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

ALTER TABLE public.moves ADD COLUMN IF NOT EXISTS game_id UUID;
ALTER TABLE public.moves ADD COLUMN IF NOT EXISTS move_order INT;
ALTER TABLE public.moves ADD COLUMN IF NOT EXISTS move VARCHAR(10);
ALTER TABLE public.moves ADD COLUMN IF NOT EXISTS fen VARCHAR(500);
ALTER TABLE public.moves ADD COLUMN IF NOT EXISTS created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP;
ALTER TABLE public.moves ADD COLUMN IF NOT EXISTS updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP;
DO $$
BEGIN
	IF NOT EXISTS (
		SELECT 1 FROM pg_constraint WHERE conname = 'games_user_id_fkey'
	) THEN
		ALTER TABLE public.games
			ADD CONSTRAINT games_user_id_fkey
			FOREIGN KEY (user_id) REFERENCES public.users(id) ON DELETE CASCADE;
	END IF;

	IF NOT EXISTS (
		SELECT 1 FROM pg_constraint WHERE conname = 'moves_game_id_fkey'
	) THEN
		ALTER TABLE public.moves
			ADD CONSTRAINT moves_game_id_fkey
			FOREIGN KEY (game_id) REFERENCES public.games(id) ON DELETE CASCADE;
	END IF;
END $$;

CREATE OR REPLACE FUNCTION public.set_move_order()
RETURNS TRIGGER AS $$
BEGIN
	IF NEW.move_order IS NULL THEN
		SELECT COALESCE(MAX(move_order), 0) + 1
		INTO NEW.move_order
		FROM public.moves
		WHERE game_id = NEW.game_id;
	END IF;
	RETURN NEW;
END;
$$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS moves_set_order ON public.moves;
CREATE TRIGGER moves_set_order
	BEFORE INSERT ON public.moves
	FOR EACH ROW
	EXECUTE FUNCTION public.set_move_order();

CREATE INDEX IF NOT EXISTS idx_games_user_id ON public.games(user_id);
CREATE INDEX IF NOT EXISTS idx_moves_game_id ON public.moves(game_id);
CREATE INDEX IF NOT EXISTS idx_moves_move_order ON public.moves(game_id, move_order);
`
