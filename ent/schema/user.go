package schema

import (
	"context"

	"entgo.io/bug/ent/hook"
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"github.com/speps/go-hashids/v2"

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
			// This hook hashes the new insert ID of the user and sets that as the "public ID" column.
			return hook.UserFunc(func(ctx context.Context, m *gen.UserMutation) (ent.Value, error) {
				// Because we rely on the database to set the ID, we cannot must allow the mutations to continue first
				// so that the sql insert happens.  After that, we can then hash the ID and update the newly created
				// row.
				v, err := next.Mutate(ctx, m)
				if err != nil {
					return nil, err
				}

				if id, ok := m.ID(); ok {
					hashID, err := HashUserID(id)
					if err != nil {
						return nil, err
					}

					newUser, err := m.Client().User.UpdateOneID(id).SetHashID(hashID).Save(ctx)
					if err != nil {
						return nil, err
					}

					// We've updated the user and should return the new value of that user so that the call-site of
					// client.User.Create()....Save(ctx) receives this updated value.  The bug is that since we're doing a
					// downstream mutation (e.g. mutating after the call to next.Mutate instead of before), this updated
					// User value is not propagated to the call-site.
					return newUser, nil
				}

				return v, nil
			})
		}, ent.OpCreate),
	}
}

func HashUserID(id int) (string, error) {
	return hashID(id, "salt is so salty")
}

func hashID(id int, salt string) (string, error) {
	hd := hashids.NewData()
	hd.Salt = salt
	hd.MinLength = 30
	h, err := hashids.NewWithData(hd)
	if err != nil {
		return "", err
	}

	return h.Encode([]int{id})
}
