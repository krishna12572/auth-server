env "local" {
  url = "postgres://auth:auth@localhost:5432/authdb?sslmode=disable"
  dev = "postgres://auth:auth@localhost:5432/authdb_dev?sslmode=disable"
  migration {
    dir = "file://migrations"
  }
}