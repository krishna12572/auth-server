package graph

import "auth-server/ent"

type ContextKey string

type Resolver struct {
Client          *ent.Client
PasswordChecker PasswordChecker
}
