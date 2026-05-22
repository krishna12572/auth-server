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

The database already contains a test user. If it doesn't, insert one manually:

```bash
docker exec auth_postgres psql -U auth -d authdb -c "INSERT INTO users (email, password_hash, created_at) VALUES ('admin@example.com', '\$2a\$10\$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', NOW());"
```

This creates a user with:
- Email: `admin@example.com`
- Password: `password`

---

## Atlas migrations

Migration files are in the `migrations/` folder, generated with Atlas.

Migrations run automatically when the server starts. To apply manually, download Atlas from https://atlasgo.io and run:

```bash
atlas migrate apply --env local
```

---

## API

Open `http://localhost:8082` in your browser for the GraphQL playground.

### Login

```graphql
mutation {
  login(email: "test5@example.com", password: "password123") {
    accessToken
    refreshToken
  }
}
```

### Get current user

Add to the **Headers** tab in GraphiQL:

```json
{ "Authorization": "Bearer YOUR_ACCESS_TOKEN" }
```

Then run:

```graphql
query {
  me {
    id
    email
  }
}
```

### Refresh token

First login to get a refresh token, then immediately run:

```graphql
mutation {
  refresh(refreshToken: "YOUR_REFRESH_TOKEN") {
    accessToken
    refreshToken
  }
}
```

> Note: refresh tokens expire after 24 hours. You must login first to get a valid one.

### Logout

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
