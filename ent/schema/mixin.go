package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/mixin"
	"github.com/google/uuid"
)

// BaseMixin holds common fields for all entities.
type BaseMixin struct {
	mixin.Schema
}

// Annotations of the BaseMixin.
func (BaseMixin) Annotations() []schema.Annotation {
	return []schema.Annotation{
		// Enable database comments for fields
		entsql.WithComments(true),
	}
}

// Fields of the BaseMixin.
func (BaseMixin) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).
			Default(func() uuid.UUID {
				u, _ := uuid.NewV7()
				return u
			}).
			Immutable().
			Comment("Primary key ID, auto-generated UUID"),
		field.Bool("enabled").
			Default(true).
			Comment("Whether the record is enabled, true means enabled, false means disabled"),
		field.Time("created_at").
			Default(time.Now).
			Immutable().
			Comment("Creation time, automatically set"),
		field.Time("updated_at").
			Default(time.Now).
			UpdateDefault(time.Now).
			Comment("Update time, automatically updated"),
	}
}
