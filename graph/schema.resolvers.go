package graph

import (
	"auth-server/ent/refreshtoken"
	"auth-server/ent/user"
	"auth-server/graph/model"
	"context"
	"fmt"
	"time"

	"github.com/99designs/gqlgen/graphql"
	"golang.org/x/crypto/bcrypt"
)

type PasswordChecker interface {
	Check(hash, password string) error
}

type BcryptChecker struct{}

func (BcryptChecker) Check(hash, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}

var defaultChecker PasswordChecker = BcryptChecker{}

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }

func (r *Resolver) Mutation() MutationResolver { return &mutationResolver{r} }
func (r *Resolver) Query() QueryResolver       { return &queryResolver{r} }

func (r *Resolver) checker() PasswordChecker {
	if r.PasswordChecker != nil {
		return r.PasswordChecker
	}
	return defaultChecker
}

func (r *mutationResolver) Login(ctx context.Context, email string, password string) (*model.AuthPayload, error) {
	u, err := r.Client.User.Query().Where(user.EmailEQ(email)).Only(ctx)
	if err != nil {
		return nil, graphql.ErrorOnPath(ctx, fmt.Errorf("user not found"))
	}
	if err := r.checker().Check(u.PasswordHash, password); err != nil {
		return nil, graphql.ErrorOnPath(ctx, fmt.Errorf("invalid password"))
	}
	access, err := generateToken(u.ID)
	if err != nil {
		return nil, graphql.ErrorOnPath(ctx, err)
	}
	refresh := generateRefreshToken()
	_, err = r.Client.RefreshToken.Create().SetToken(refresh).SetUser(u).SetExpiresAt(time.Now().Add(24 * time.Hour)).Save(ctx)
	if err != nil {
		return nil, graphql.ErrorOnPath(ctx, err)
	}
	return &model.AuthPayload{AccessToken: access, RefreshToken: refresh}, nil
}

func (r *mutationResolver) Refresh(ctx context.Context, refreshToken string) (*model.AuthPayload, error) {
	rt, err := r.Client.RefreshToken.Query().Where(refreshtoken.TokenEQ(refreshToken)).WithUser().Only(ctx)
	if err != nil {
		return nil, graphql.ErrorOnPath(ctx, fmt.Errorf("invalid refresh token"))
	}
	if time.Now().After(rt.ExpiresAt) {
		return nil, graphql.ErrorOnPath(ctx, fmt.Errorf("refresh token expired"))
	}
	_ = r.Client.RefreshToken.DeleteOne(rt).Exec(ctx)
	newAccess, err := generateToken(rt.Edges.User.ID)
	if err != nil {
		return nil, graphql.ErrorOnPath(ctx, err)
	}
	newRefresh := generateRefreshToken()
	_, err = r.Client.RefreshToken.Create().SetToken(newRefresh).SetUser(rt.Edges.User).SetExpiresAt(time.Now().Add(24 * time.Hour)).Save(ctx)
	if err != nil {
		return nil, graphql.ErrorOnPath(ctx, err)
	}
	return &model.AuthPayload{AccessToken: newAccess, RefreshToken: newRefresh}, nil
}

func (r *mutationResolver) Logout(ctx context.Context, refreshToken string) (bool, error) {
	_, err := r.Client.RefreshToken.Delete().Where(refreshtoken.TokenEQ(refreshToken)).Exec(ctx)
	if err != nil {
		return false, graphql.ErrorOnPath(ctx, err)
	}
	return true, nil
}

func (r *queryResolver) Me(ctx context.Context) (*model.User, error) {
	tokenStr, ok := ctx.Value(ContextKey("token")).(string)
	if !ok || tokenStr == "" {
		return nil, graphql.ErrorOnPath(ctx, fmt.Errorf("unauthenticated"))
	}
	claims, err := validateToken(tokenStr)
	if err != nil {
		return nil, graphql.ErrorOnPath(ctx, fmt.Errorf("invalid token"))
	}
	userID := int(claims["user_id"].(float64))
	u, err := r.Client.User.Get(ctx, userID)
	if err != nil {
		return nil, graphql.ErrorOnPath(ctx, fmt.Errorf("user not found"))
	}
	return &model.User{ID: fmt.Sprintf("%d", u.ID), Email: u.Email}, nil
}