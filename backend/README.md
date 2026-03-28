# Canh Buom Xanh Backend (Part H)

## Stack
- Go
- Fiber
- MySQL
- JWT (HS256)
- Docker (release packaging)

## Structure
- cmd/api: entrypoint
- internal/app: server bootstrap
- internal/config: env config loader
- internal/router: API route registry
- internal/modules/auth: auth module (register, login, me)
- internal/middleware: HTTP middleware (JWT guard)
- internal/platform/database: db client
- internal/platform/security: jwt/password helpers
- db/schema.sql: initial schema
- scripts/load-check.sh: load validation helper
- scripts/abuse-check.sh: rate-limit abuse validation helper

## Run Locally
1. Copy env file:
   - `cp .env.example .env`
2. Apply database schema:
   - `mysql -u root -p < db/schema.sql`
3. Start API:
   - `go run ./cmd/api`

## Core Endpoints
- GET /api/v1/health
- GET /api/v1/metrics
- POST /api/v1/auth/register
- POST /api/v1/auth/login
- GET /api/v1/auth/me (requires Bearer token)

## Hardening Defaults
- Password hashing uses bcrypt.
- Request middleware stack includes request-id, panic recovery, structured request logs, security headers, CORS, and per-path rate limit.
- Auth validation enforces email format and minimum password length.

## Hardening Environment Variables
- `ALLOWED_ORIGINS`: comma-separated allowed CORS origins.
- `RATE_LIMIT_MAX`: max requests per key per window.
- `RATE_LIMIT_WINDOW_SECONDS`: limiter window in seconds.

## Validation Scripts
- Load check:
  - `./scripts/load-check.sh http://127.0.0.1:8080 200 20`
- Abuse check (expect partial HTTP 429):
  - `./scripts/abuse-check.sh http://127.0.0.1:8080 180`

## Release with Docker
1. From repository root:
   - `docker compose up --build`
2. API is available at:
   - `http://127.0.0.1:8080`
