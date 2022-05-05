package schema

import (
	"context"

	"entgo.io/bug/ent/hook"
	"entgo.io/ent"
	"entgo.io/ent/schema/field"

	gen "entgo.io/bug/ent"
)

// User holds the schema definition for the User entity.
type User struct {
	ent.Schema
}

// Fields of the User.
func (User) Fields() []ent.Field {
	return []ent.Field{
		field.String("hash_id").Optional(),
		field.Int("age"),
		field.String("name"),
	}
}

// Edges of the User.
func (User) Edges() []ent.Edge {
	return nil
}

// Hooks of the User.
func (User) Hooks() []ent.Hook {
	return []ent.Hook{
		hook.On(func(next ent.Mutator) ent.Mutator {
			return hook.UserFunc(func(ctx context.Context, m *gen.UserMutation) (ent.Value, error) {
				return next.Mutate(ctx, m)
			})
		}, ent.OpCreate),
	}
}
