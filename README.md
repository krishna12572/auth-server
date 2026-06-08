markdown# auth-server

A GraphQL authentication server built with Go, Ent, and gqlgen. Provides JWT-based login, token refresh, and logout.

## Stack

- Go 1.22+
- PostgreSQL 16 (Docker)
- Ent (ORM with code generation)
- gqlgen (GraphQL)
- Atlas (migrations)
- JWT + bcrypt

---

## Requirements

- Go 1.22+ — https://go.dev/dl
- Docker Desktop — https://www.docker.com/products/docker-desktop

---

## Setup

### 1. Clone the repo
git clone https://github.com/krishna12572/auth-server
cd auth-server

### 2. Configure environment

Windows:
copy .env.example .env

Mac/Linux:
cp .env.example .env

The default values work out of the box — no editing needed unless you want custom DB credentials.

### 3. Start the database
docker compose up -d

You should see: `Container auth_postgres Started`

### 4. Run the server
go run server.go

Server starts at http://localhost:8082. Database migrations run automatically on startup.

### 5. Seed the database (required before first login)

Connect to the database:
docker exec -it auth_postgres psql -U auth -d authdb

Paste this SQL and press Enter:
INSERT INTO users (email, password_hash, created_at) VALUES ('admin@example.com', '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', NOW());

Type `\q` and press Enter to exit psql.

This creates a user with:
- Email: `admin@example.com`
- Password: `password`

---

## Tests
go test -v ./...

Tests cover: password hashing, bcrypt salting, JWT generation, token expiry, wrong secret detection, tampered token detection, login flow, and refresh token uniqueness.

---

## Automation (Makefile)
make run          # start the server
make build        # compile to ./bin/auth-server
make test         # run all tests with race detector
make docker-up    # start PostgreSQL container
make docker-down  # stop PostgreSQL container
make generate     # re-run Ent + gqlgen code generation

---

## API

Open http://localhost:8082 in your browser for the GraphQL playground.

### Login
```graphql
mutation {
  login(email: "admin@example.com", password: "password") {
    accessToken
    refreshToken
  }
}
```

### Me
Add header: `Authorization: Bearer <accessToken>`
```graphql
query {
  me {
    id
    email
  }
}
```

### Refresh
```graphql
mutation {
  refresh(refreshToken: "YOUR_REFRESH_TOKEN") {
    accessToken
    refreshToken
  }
}
```

### Logout
```graphql
mutation {
  logout(refreshToken: "YOUR_REFRESH_TOKEN")
}
```

---

## Project Structure
auth-server/
├── ent/
│   ├── schema/
│   │   ├── user.go          # Custom Email type with validation, privacy policy, edges
│   │   ├── refreshtoken.go  # RefreshToken schema with privacy policy
│   │   ├── privacy.go       # AllowIfAdmin, AllowIfOwner, DenyIfNoViewer rules
│   │   └── schema_test.go   # Email validation tests + policy context tests
├── graph/
│   ├── helpers.go           # generateToken, validateToken, generateRefreshToken
│   ├── resolver.go          # Root resolver with injectable PasswordChecker interface
│   └── schema.resolvers.go  # Login, Refresh, Logout, Me — all use graphql.ErrorOnPath
├── migrations/              # Atlas SQL migration files
├── .env.example             # Environment variable template (safe to commit)
├── .env                     # Your local config (git-ignored)
├── docker-compose.yml       # Reads DB config from .env
├── atlas.hcl                # Reads DB config from .env
├── Makefile                 # Automation targets
└── server.go                # Entry point — reads all config from .env via godotenv
