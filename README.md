# auth-server

A JWT-based authentication server built with Go and GraphQL. Handles login, logout, token refresh, and protected routes.

## Stack

- Go
- PostgreSQL (Docker)
- Ent
- gqlgen
- Atlas
- JWT + bcrypt

---

## Requirements

- Git
- Go 1.21+
- Docker

---

## Installing Requirements (Windows)

Before setting up the project, make sure you have the following tools installed.

### 1. Install Git

Download the installer from https://git-scm.com/download/win and run it. Click through all the defaults — they are fine for most users. Once installed, close and reopen your Command Prompt.

### 2. Install Go

Download the Windows installer from https://go.dev/dl/ and run it.

### 3. Install Docker Desktop

Download Docker Desktop from https://www.docker.com/products/docker-desktop/ and install it. You do **not** need to sign in to a Docker account to use it locally — just launch Docker Desktop and skip or close any sign-in prompt.

### Verify your installations

After installing, open a new Command Prompt and run:

```bash
git --version
go version
docker --version
```

All three should print a version number. If any fail, reinstall that tool and reopen your terminal.

---

## Setup

### 1. Clone the repository

```bash
git clone https://github.com/krishna12572/auth-server
cd auth-server
```

### 2. Start the database

```bash
docker compose up -d
```

You should see output like `Container auth_postgres Started`.

### 3. Start the server

```bash
go run server.go
```

Migrations run automatically on startup. The port is read from `.env` (default 8082).

### 4. Seed the database (Required)

The database starts empty — you **must** insert a test user manually before you can log in. Open a new terminal and run:

```bash
docker exec auth_postgres psql -U auth -d authdb -c "INSERT INTO users (email, password_hash, created_at) VALUES ('admin@example.com', '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', NOW());"
```

You should see `INSERT 0 1` — this confirms the user was created successfully.

This creates a user with:
- Email: `admin@example.com`
- Password: `password`

---

## Atlas Migrations

Migration files are in the `migrations/` folder, generated with Atlas.

Migrations run automatically when the server starts. To apply manually, download Atlas from https://atlasgo.io and run:

```bash
atlas migrate apply --env local
```

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

### Get current user

Add to the Headers tab in GraphiQL:

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
