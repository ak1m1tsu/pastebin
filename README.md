# Pastebin

## How to run?

Rename `.env.example` to `.env` and run next command:

```shell
make docker/up
```

Then in `postgres` container shell run SQL commands from `schema/schema.sql`

```shell
docker compose exec postgres psql -U pastebin -d pastebin
```

Next run migrations:

```shell
 migrate -path=./migrations -database 'postgres://pastebin:pastebin@localhost:5432/pastebin?sslmode=disable&search_path=public' up
```

Finally go to `localhost:8085/api/v1/swagger/index.html` to see swagger API docs.
