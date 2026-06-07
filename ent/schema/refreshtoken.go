package schema

import (
"entgo.io/ent"
"entgo.io/ent/schema/edge"
"entgo.io/ent/schema/field"
)

type RefreshToken struct {
ent.Schema
}

func (RefreshToken) Mixin() []ent.Mixin {
return []ent.Mixin{UserMixin{}}
}

func (RefreshToken) Fields() []ent.Field {
return []ent.Field{
field.String("token").Unique(),
field.Time("expires_at"),
}
}

func (RefreshToken) Edges() []ent.Edge {
return []ent.Edge{
edge.From("user", User.Type).Ref("refresh_tokens").Unique().Required(),
}
}
