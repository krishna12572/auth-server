env "local" {
  url = "postgres://${DB_USER}:${DB_PASSWORD}@localhost:${DB_PORT}/${DB_NAME}?sslmode=disable"
  dev = "postgres://${DB_USER}:${DB_PASSWORD}@localhost:${DB_PORT}/${DB_NAME}_dev?sslmode=disable"
  migration {
    dir = "file://migrations"
  }
}
