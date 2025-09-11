ALTER TABLE public.games
    ADD COLUMN user_id uuid REFERENCES public.users(id);