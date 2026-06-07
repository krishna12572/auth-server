package main

import (
"auth-server/ent"
"auth-server/graph"
"context"
"fmt"
"log"
"net/http"
"os"
"strings"

"github.com/99designs/gqlgen/graphql/handler"
"github.com/99designs/gqlgen/graphql/handler/transport"
"github.com/99designs/gqlgen/graphql/playground"
"github.com/joho/godotenv"
_ "github.com/lib/pq"
)

const defaultPort = "8082"

func main() {
_ = godotenv.Load()

port := os.Getenv("PORT")
if port == "" {
port = defaultPort
}

dsn := fmt.Sprintf(
"host=localhost port=%s user=%s password=%s dbname=%s sslmode=disable",
os.Getenv("DB_PORT"),
os.Getenv("DB_USER"),
os.Getenv("DB_PASSWORD"),
os.Getenv("DB_NAME"),
)

client, err := ent.Open("postgres", dsn)
if err != nil {
log.Fatalf("failed to connect to database: %v", err)
}
defer client.Close()

ctx := context.Background()
if err := client.Schema.Create(ctx); err != nil {
log.Fatalf("failed to run schema migration: %v", err)
}

srv := handler.New(graph.NewExecutableSchema(graph.Config{
Resolvers: &graph.Resolver{Client: client},
}))
srv.AddTransport(transport.Options{})
srv.AddTransport(transport.GET{})
srv.AddTransport(transport.POST{})

http.Handle("/", playground.Handler("GraphQL playground", "/query"))
http.Handle("/query", authMiddleware(srv))

log.Printf("server running at http://localhost:%s/", port)
log.Fatal(http.ListenAndServe(":"+port, nil))
}

func authMiddleware(next http.Handler) http.Handler {
return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
header := r.Header.Get("Authorization")
if strings.HasPrefix(header, "Bearer ") {
token := strings.TrimPrefix(header, "Bearer ")
ctx := context.WithValue(r.Context(), graph.ContextKey("token"), token)
r = r.WithContext(ctx)
}
next.ServeHTTP(w, r)
})
}
