ALTER TABLE public.games
    ADD COLUMN move_amount integer,
    ADD COLUMN end_type varchar(50),
    ADD COLUMN result varchar(50);

ALTER TABLE public.users ADD COLUMN username varchar(255);

ALTER TABLE public.moves ADD COLUMN user_id uuid REFERENCES public.users(id);