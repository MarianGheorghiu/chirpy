# chirpy

A tiny Twitter-style API written in Go with PostgreSQL.

## Setup

### Requirements
- Go 1.22+
- PostgreSQL

### Environment
Create a `.env` file in the project root:
```env
JWT_SECRET="your-jwt-secret"
POLKA_KEY="f271c81ff7084ee5b99a5091b42d486e"
DB_URL="postgres://user:pass@localhost:5432/chirpy?sslmode=disable"
PLATFORM="dev"

Database

Run migrations (example with goose):

goose postgres "$DB_URL" up


If you change queries, regenerate sqlc types:

sqlc generate

Run
go run .
# listens on :8080

Minimal API
Users

POST /api/users — create user (returns id, created_at, updated_at, email, is_chirpy_red)

POST /api/login — login, returns JWT access (and refresh flow endpoints)

PUT /api/users — update own email & password
Header: Authorization: Bearer <token>

Chirps

POST /api/chirps — create chirp (auth; 140 chars; basic bad-word filter)

GET /api/chirps — list chirps
Query params:

author_id=<uuid> (optional)

sort=asc|desc (optional, default asc)

GET /api/chirps/{chirpID} — get one chirp

DELETE /api/chirps/{chirpID} — delete own chirp
Success: 204; not author: 403; not found: 404

Webhooks (Polka)

POST /api/polka/webhooks
Header: Authorization: ApiKey <POLKA_KEY>
Body:

{"event":"user.upgraded","data":{"user_id":"<uuid>"}}


Unknown events → 204

User upgraded → 204

User not found → 404

Notes

Passwords hashed with argon2id

Auth tokens are JWT (HS256)

SQL generated with sqlc, migrations with goose