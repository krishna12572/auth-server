# auth-server

A JWT-based authentication server built with Go and GraphQL. Handles login, logout, token refresh, and protected routes.

## Stack

- Go
- PostgreSQL (Docker)
- Ent
- gqlgen
- Atlas
- JWT + bcrypt

## Requirements

- Go 1.21+
- Docker

---

## Setup

### 1. Start the database

```bash
docker compose up -d
```

### 2. Start the server

```bash
go run server.go
```

Migrations run automatically on startup. The port is read from `.env` (default 8082).

### 3. Seed a user

```bash
go run main.go
```

Creates a test user. If you get a duplicate email error, change the email in `main.go` and run again.

---

## Atlas migrations

Migration files are in the `migrations/` folder, generated with Atlas.

To apply migrations manually:

```bash
.\atlas.exe migrate apply --env local
```

---

## API

Open `http://localhost:8082` in your browser for the GraphQL playground.

**Login**
```graphql
mutation {
  login(email: "test5@example.com", password: "password123") {
    accessToken
    refreshToken
  }
}
```

**Get current user**

Add to Headers tab:
```json
{ "Authorization": "Bearer YOUR_ACCESS_TOKEN" }
```

```graphql
query {
  me {
    id
    email
  }
}
```

**Refresh**
```graphql
mutation {
  refresh(refreshToken: "YOUR_REFRESH_TOKEN") {
    accessToken
    refreshToken
  }
}
```

**Logout**
```graphql
mutation {
  logout(refreshToken: "YOUR_REFRESH_TOKEN")
}
```

---

## Tests

```bash
go test -v
```

Covers password hashing, JWT generation/validation, token expiry, login success/failure, and refresh token rotation.

---

## How it works

Login checks the password with bcrypt. If correct, it returns a JWT access token (1 hour expiry) and a refresh token stored in the database. The `me` query reads the user ID from the token in the Authorization header. Refreshing deletes the old token and issues a new pair. Logout just deletes the refresh token from the database.
