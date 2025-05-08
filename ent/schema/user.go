package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"github.com/google/uuid"
)

// User holds the schema definition for the User entity.
type User struct {
	ent.Schema
}

// Annotations of the User.
func (User) Annotations() []schema.Annotation {
	return []schema.Annotation{
		// Enable database comments for table and fields
		entsql.WithComments(true),
		// Add table comment that appears in both generated code and database
		schema.Comment("User table, stores information for all users in the system"),
	}
}

// Fields of the User.
func (User) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("tenant_id", uuid.UUID{}).
			Comment("Associated tenant ID"),
		field.String("username").
			NotEmpty().
			MaxLen(100).
			Comment("Username for login"),
		field.String("given_name").
			Optional().
			MaxLen(100).
			Comment("Display name of the user"),
		field.Bytes("avatar").
			Nillable().Optional().
			Comment("User avatar, stored in binary format"),
		field.String("email").
			NotEmpty().
			MaxLen(255).
			Comment("User email address"),
		field.Bool("email_verified").
			Default(false).
			Comment("Whether the email has been verified"),
		field.String("password_bcrypt").
			NotEmpty().
			MaxLen(255).
			Comment("Encrypted password in bcrypt format"),
		field.String("totp_secret").
			Nillable().Optional().
			MaxLen(255).
			Comment("Two-factor authentication secret key"),
		field.Bool("totp_enabled").
			Default(false).
			Comment("Whether two-factor authentication is enabled"),
		field.Time("last_login_at").
			Nillable().Optional().
			Comment("Last login time"),
		field.Time("last_logout_at").
			Nillable().Optional().
			Comment("Last logout time"),
		field.String("last_login_ip").
			Nillable().Optional().
			MaxLen(45).
			Comment("Last login IP address"),
		field.Bool("is_superuser").
			Default(false).
			Comment("Whether the user is a superadmin"),
		field.Bool("is_locked").
			Default(false).
			Comment("Whether the user is locked"),
		field.String("lock_reason").
			Nillable().Optional().
			MaxLen(255).
			Comment("Reason for user account being locked"),
		field.JSON("custattr", map[string]any{}).
			Optional().
			Comment("Custom attributes stored in JSON format"),
	}
}

// Edges of the User.
func (User) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("tenant", Tenant.Type).
			Ref("users").
			Field("tenant_id").
			Unique().
			Required().
			Comment("Tenant this user belongs to"), // Many-to-one: User -> Tenant
	}
}

// Annotations for unique constraints
func (User) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("tenant_id", "username").Unique(), // Unique per tenant: username
		index.Fields("tenant_id", "email").Unique(),    // Unique per tenant: email
	}
}

func (User) Mixin() []ent.Mixin {
	return []ent.Mixin{
		BaseMixin{},
	}
}
