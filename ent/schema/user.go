package schema

import (
	"fmt"
	"regexp"
	"time"

	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/privacy"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/mixin"
)

type Email string

func (e Email) Validate() error {
	re := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	if !re.MatchString(string(e)) {
		return fmt.Errorf("invalid email address: %q", e)
	}
	return nil
}

func (e *Email) UnmarshalText(text []byte) error {
	*e = Email(text)
	return e.Validate()
}

func (e Email) MarshalText() ([]byte, error) {
	if err := e.Validate(); err != nil {
		return nil, err
	}
	return []byte(e), nil
}

type UserMixin struct {
	mixin.Schema
}

func (UserMixin) Policy() ent.Policy {
	return privacy.Policy{
		Mutation: privacy.MutationPolicy{
			AllowIfAdmin(),
			AllowIfOwner(),
			DenyIfNoViewer(),
		},
	}
}

type User struct {
	ent.Schema
}

func (User) Mixin() []ent.Mixin {
	return []ent.Mixin{UserMixin{}}
}

func (User) Fields() []ent.Field {
	return []ent.Field{
		field.String("email").GoType(Email("")).Unique().Annotations(entsql.Annotation{Size: 254}),
		field.String("password_hash"),
		field.Time("created_at").Default(time.Now).Immutable(),
	}
}

func (User) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("refresh_tokens", RefreshToken.Type),
	}
}