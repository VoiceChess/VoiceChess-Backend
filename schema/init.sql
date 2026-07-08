CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS public.users (
    id TEXT PRIMARY KEY,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS public.games (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id TEXT NOT NULL REFERENCES public.users(id) ON DELETE CASCADE,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS public.moves (
    id UUID DEFAULT uuid_generate_v4() PRIMARY KEY NOT NULL,
    game_id UUID NOT NULL REFERENCES public.games(id) ON DELETE CASCADE,
    move_order INT,
    move VARCHAR(10) NOT NULL,
    fen VARCHAR(500) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

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
