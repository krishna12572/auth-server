# auth-server

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

```bash
git clone https://github.com/krishna12572/auth-server
cd auth-server
```

### 2. Configure environment

```bash
cp .env.example .env
```

Edit `.env` with your values (defaults work out of the box):
PORT=8082
DB_USER=auth
DB_PASSWORD=auth
DB_PORT=5432
DB_NAME=authdb

### 3. Start the database

```bash
docker compose up -d
```

### 4. Run the server

```bash
go run server.go
```

Server starts at http://localhost:8082. Migrations run automatically.

---

## Seed the database

The database starts empty. Insert a test user:

```bash
docker exec -it auth_postgres psql -U auth -d authdb
```

Then paste:

```sql
INSERT INTO users (email, password_hash, created_at)
VALUES ('admin@example.com', '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', NOW());
```

Type `\q` to exit. This creates:
- Email: `admin@example.com`
- Password: `password`

---

## Tests

```bash
go test -v ./...
```

Covers password hashing, JWT generation/validation/expiry, token tampering, login flow, and refresh token rotation.

---

## Automation

```bash
make run          # start the server
make build        # compile to ./bin/auth-server
make test         # run all tests with race detector
make docker-up    # start PostgreSQL
make docker-down  # stop PostgreSQL
make generate     # re-run Ent + gqlgen codegen
```

---

## API

Open http://localhost:8082 for the GraphQL playground.

### Login
```graphql
mutation {
  login(email: "admin@example.com", password: "password") {
    accessToken
    refreshToken
  }
}
```

### Me (add `Authorization: Bearer <token>` header)
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
│   │   ├── user.go          # Custom Email type, privacy policy, edges
│   │   ├── refreshtoken.go  # RefreshToken schema with policy
│   │   ├── privacy.go       # AllowIfAdmin, AllowIfOwner, DenyIfNoViewer rules
│   │   └── schema_test.go   # Email validation + policy tests
├── graph/
│   ├── helpers.go           # generateToken, validateToken, generateRefreshToken
│   ├── resolver.go          # Root resolver with PasswordChecker interface
│   └── schema.resolvers.go  # Login, Refresh, Logout, Me resolvers
├── .env.example             # Environment variable template
├── docker-compose.yml       # Uses .env variables
├── atlas.hcl                # Uses .env variables
├── Makefile                 # Automation targets
└── server.go                # Entry point, reads all config from .env