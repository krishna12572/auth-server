package schema

import (
	"context"

	"entgo.io/ent"
	"entgo.io/ent/privacy"
)

type ViewerContext struct {
	UserID  int
	IsAdmin bool
}

type viewerCtxKey struct{}

func NewViewerContext(ctx context.Context, v ViewerContext) context.Context {
	return context.WithValue(ctx, viewerCtxKey{}, v)
}

func ViewerFromContext(ctx context.Context) (ViewerContext, bool) {
	v, ok := ctx.Value(viewerCtxKey{}).(ViewerContext)
	return v, ok
}

func AllowIfAdmin() privacy.MutationRuleFunc {
	return privacy.MutationRuleFunc(func(ctx context.Context, m ent.Mutation) error {
		v, ok := ViewerFromContext(ctx)
		if ok && v.IsAdmin {
			return privacy.Allow
		}
		return privacy.Skip
	})
}

func AllowIfOwner() privacy.MutationRuleFunc {
	return privacy.MutationRuleFunc(func(ctx context.Context, m ent.Mutation) error {
		v, ok := ViewerFromContext(ctx)
		if !ok {
			return privacy.Deny
		}
		type idMutation interface {
			ID() (int, bool)
		}
		if idm, ok := m.(idMutation); ok {
			if id, exists := idm.ID(); exists && id == v.UserID {
				return privacy.Allow
			}
		}
		return privacy.Skip
	})
}

func DenyIfNoViewer() privacy.MutationRuleFunc {
	return privacy.MutationRuleFunc(func(ctx context.Context, m ent.Mutation) error {
		_, ok := ViewerFromContext(ctx)
		if !ok {
			return privacy.Deny
		}
		return privacy.Skip
	})
}
