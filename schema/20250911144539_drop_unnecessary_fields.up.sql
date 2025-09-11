ALTER TABLE public.games
    DROP COLUMN move_amount,
    DROP COLUMN end_type,
    DROP COLUMN result;

ALTER TABLE public.users DROP COLUMN username;

ALTER TABLE public.moves DROP COLUMN user_id;