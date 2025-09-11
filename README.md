services folder is used to store business logics
controllers folder is used for http handling
models folder is used to store entity definitions
routes folder is used to define api routes
config folder is used to store configuration files
repo folder is used to interact with database

Database setup:
Install database migration tool: https://github.com/golang-migrate/migrate/tree/master/cmd/migrate

Create the database:
```
user=# CREATE DATABASE chessmate;
CREATE DATABASE

```

Export the PostgreSQL url to env.
```
export POSTGRESQL_URL='postgres://postgres@localhost:5432/chessmate?sslmode=disable'
```

Run this command to migrate database (files inside schema), adjust db connection with our setup:
```
migrate -verbose -path 'schema/' -database ${POSTGRESQL_URL} up
```

To create a new migration, run:
```
migrate create -ext sql -dir schema/ -format 20060102150405 init_mg
```


