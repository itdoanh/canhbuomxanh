# Part 8 (Part H) - Hardening and Release Delivery

## Scope Completed
1. Security hardening for authentication and runtime middleware.
2. Observability endpoints for operational checks.
3. Load and abuse validation scripts for repeatable checks.
4. Release packaging baseline with containerized deployment flow.

## Security Hardening Delivered
- Password hashing migrated to bcrypt in backend/internal/platform/security/password.go.
- Auth validation now enforces:
  - normalized email format
  - minimum password length (8)
- Global middleware stack in backend/internal/app/server.go:
  - request id
  - panic recovery
  - request logger
  - security headers
  - CORS with configurable allowed origins
  - rate limiter with configurable window/max
- New config fields in backend/internal/config/config.go:
  - ALLOWED_ORIGINS
  - RATE_LIMIT_MAX
  - RATE_LIMIT_WINDOW_SECONDS

## Observability Delivered
- GET /api/v1/health
  - returns service status, db status, timestamp.
- GET /api/v1/metrics
  - returns basic counters:
    - users
    - posts
    - votes
    - queuedFlags
  - returns service name and timestamp.

## Validation Tooling Delivered
- backend/scripts/load-check.sh
  - concurrent health endpoint load probe.
- backend/scripts/abuse-check.sh
  - burst probe to confirm limiter behavior via HTTP 429.

## Release Artifacts Delivered
- backend/Dockerfile for API image build.
- docker-compose.yml at repository root with:
  - mysql service
  - api service
  - schema auto-bootstrap for local release simulation.

## Operational Notes
- Default limiter key is ip + path.
- Hardening values are env-driven and can be tuned per environment.
- For production, JWT secret and DB credentials must be replaced with secure values.

## Exit Criteria Check
- Security baseline: completed.
- Metrics baseline: completed.
- Abuse/load reproducibility: completed.
- Release runbook baseline: completed.
