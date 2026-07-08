ALTER TABLE IF EXISTS public.moves DROP CONSTRAINT IF EXISTS moves_user_id_fkey;
ALTER TABLE IF EXISTS public.moves DROP CONSTRAINT IF EXISTS moves_game_id_fkey;
ALTER TABLE IF EXISTS public.games DROP CONSTRAINT IF EXISTS games_user_id_fkey;

ALTER TABLE public.users ALTER COLUMN id TYPE TEXT USING id::TEXT;
ALTER TABLE public.games ALTER COLUMN user_id TYPE TEXT USING user_id::TEXT;

ALTER TABLE public.games DROP COLUMN IF EXISTS fen;
ALTER TABLE public.games DROP COLUMN IF EXISTS result;
ALTER TABLE public.games DROP COLUMN IF EXISTS end_type;
ALTER TABLE public.games DROP COLUMN IF EXISTS move_amount;

ALTER TABLE public.games
    ADD CONSTRAINT games_user_id_fkey
    FOREIGN KEY (user_id) REFERENCES public.users(id) ON DELETE CASCADE;

ALTER TABLE public.moves
    ADD CONSTRAINT moves_game_id_fkey
    FOREIGN KEY (game_id) REFERENCES public.games(id) ON DELETE CASCADE;
