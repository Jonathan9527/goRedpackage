# learnGO

A minimal Go Gin web service.

## Run

```sh
docker compose up -d
go run main.go
```

The server listens on `http://localhost:8080`.

## Hot Reload

Install Air once:

```sh
go install github.com/air-verse/air@latest
```

Start local dependencies and run the development server:

```sh
docker compose up -d
air
```

Air watches Go files and automatically rebuilds and restarts the service.

## Routes

- `GET /` returns a greeting.
- `GET /health` returns service health.
- `GET /users?limit=20&offset=0` returns users.
- `GET /users/:account` returns one user by account.

## Project Structure

- `main.go`: application entrypoint.
- `internal/config`: environment-based application configuration.
- `internal/database`: database connection setup.
- `internal/handler`: HTTP request handlers.
- `internal/middleware`: Gin middleware.
- `internal/repository`: database queries.
- `internal/router`: route registration and route-level middleware binding.
- `internal/service`: business logic.
- `internal/model`: response and data models.

## Local Services

Start PostgreSQL, Redis, and RabbitMQ with Docker:

```sh
docker compose up -d
```

Default connection settings:

- PostgreSQL
- Host: `localhost`
- Port: `5432`
- User: `learngo`
- Password: `learngo_password`
- Database: `learngo`

- Redis
- Host: `localhost`
- Port: `6379`

- RabbitMQ
- Host: `localhost`
- AMQP Port: `5672`
- Management UI: `http://localhost:15672`
- Username: `learngo`
- Password: `learngo_password`

You can override them with environment variables from `.env.example`.

The application creates and seeds the `users` table on startup:

- `account`: unique user account, for example `user0001`.
- `username`: Chinese username, for example `用户0001`.
- `balance`: user balance stored as `NUMERIC(12, 2)`.
- Initial seed data: `1000` users with `balance = 100.00`.

## Rate Limiting

The user detail route uses a token bucket middleware:

- Bucket capacity: `10` requests.
- Refill rate: `5` requests per second.
- Rejected requests return `429 Too Many Requests`.
- Limited route: `GET /users/:account`.
