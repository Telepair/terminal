package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// Tenant holds the schema definition for the Tenant entity.
type Tenant struct {
	ent.Schema
}

// Annotations of the Tenant.
func (Tenant) Annotations() []schema.Annotation {
	return []schema.Annotation{
		// Enable database comments for table and fields
		entsql.WithComments(true),
		// Add table comment that appears in both generated code and database
		schema.Comment("Tenant table, stores information for all tenants in the system"),
	}
}

// Fields of the Tenant.
func (Tenant) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").
			NotEmpty().
			Unique().
			Comment("Unique identifier name for the tenant"),
		field.String("given_name").
			Optional().
			MaxLen(100).
			Comment("Display name of the tenant"),
		field.Bytes("avatar").
			Optional().
			Comment("Tenant avatar, stored in binary format"),
		field.Text("public_key").
			Optional().
			Comment("Tenant public key for security verification"),
		field.Bool("allow_registration").
			Default(true).
			Comment("Whether new user registration is allowed"),
		field.JSON("allowed_email_domains", []string{}).
			Optional().
			Comment("List of allowed email domains for user registration (e.g. ['company.com', 'org.edu'])"),
		field.String("admin_email").
			NotEmpty().
			MaxLen(255).
			Comment("Admin email address for tenant management"),
		field.JSON("custattr", map[string]any{}).
			Optional().
			Comment("Custom attributes stored in JSON format"),
		field.Text("description").
			Optional().
			Comment("Tenant description"),
	}
}

// Edges of the Tenant.
func (Tenant) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("users", User.Type).
			Comment("List of users associated with this tenant"),
	}
}

func (Tenant) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("name").Unique(),
	}
}

func (Tenant) Mixin() []ent.Mixin {
	return []ent.Mixin{
		BaseMixin{},
	}
}
