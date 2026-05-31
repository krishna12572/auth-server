auth-server README with copyable code blocks
auth-server
A JWT-based authentication server built with Go and GraphQL. Handles login, logout, token refresh, and protected routes.

Go
PostgreSQL
Ent
gqlgen
Atlas
JWT + bcrypt
Requirements
Git · Go 1.21+ · Docker

Installing requirements (Windows)
1
Install Git

Download the installer from https://git-scm.com/download/win and run it. Click through all the defaults — they are fine for most users. Once installed, close and reopen your Command Prompt.

2
Install Go

Download the Windows installer from https://go.dev/dl/ and run it.

3
Install Docker Desktop

Download from https://www.docker.com/products/docker-desktop/. You do not need to sign in to a Docker account — just launch Docker Desktop and skip or close any sign-in prompt.

After installing, verify all three with:

git --version
go version
docker --version

Copy
Setup
1
Clone the repository

git clone https://github.com/krishna12572/auth-server
cd auth-server

Copy
2
Start the database

docker compose up -d

Copy
You should see: Container auth_postgres Started

3
Start the server

go run server.go

Copy
Migrations run automatically on startup. The port is read from .env (default 8082).

4
Seed the database (required)

The database starts empty — you must insert a test user manually before you can log in. Open a new terminal and run:

docker exec auth_postgres psql -U auth -d authdb -c "INSERT INTO users (email, password_hash, created_at) VALUES ('admin@example.com', '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', NOW());"

Copy
You should see INSERT 0 1 — this confirms the user was created.

Email
admin@example.com
Password
password
Atlas migrations
Migration files are in the migrations/ folder, generated with Atlas. Migrations run automatically when the server starts. To apply manually, download Atlas from https://atlasgo.io and run:

atlas migrate apply --env local

Copy
API
Open http://localhost:8082 in your browser for the GraphQL playground.

Login
mutation {
  login(email: "admin@example.com", password: "password") {
    accessToken
    refreshToken
  }
}

Copy
Get current user
Add to the Headers tab in GraphiQL:

{ "Authorization": "Bearer YOUR_ACCESS_TOKEN" }

Copy
Then run:

query {
  me {
    id
    email
  }
}

Copy
Refresh token
First login to get a refresh token, then immediately run:

mutation {
  refresh(refreshToken: "YOUR_REFRESH_TOKEN") {
    accessToken
    refreshToken
  }
}

Copy
Refresh tokens expire after 24 hours. You must login first to get a valid one.
Logout
mutation {
  logout(refreshToken: "YOUR_REFRESH_TOKEN")
}

Copy
Tests
go test -v

Copy
Covers password hashing, JWT generation/validation, token expiry, login success/failure, and refresh token rotation.

How it works
Login

Checks the password with bcrypt. If correct, returns a JWT access token (1 hour expiry) and a refresh token stored in the database.

me query

Reads the user ID from the token in the Authorization header.

Refresh

Deletes the old token and issues a new pair.

Logout

Deletes the refresh token from the database.

