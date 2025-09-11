package pg_sql

var (
	CreateUser = `
	INSERT INTO public.users (id) VALUES ($1);
	`
)
