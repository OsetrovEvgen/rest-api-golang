package model

import validation "github.com/go-ozzo/ozzo-validation"

// Column ...
type Column struct {
	ID        *string `json:"id,omitempty"`
	ProjectID *string `json:"project_id,omitempty"`
	Name      *string `json:"name,omitempty"`
	Position  *int    `json:"position,omitempty"`
}

// Validate ...
func (c *Column) Validate() error {
	return validation.ValidateStruct(
		c,
		validation.Field(
			&c.Name,
			validation.Required,
			validation.Length(1, 255),
		),
		validation.Field(
			&c.ProjectID,
			validation.Required,
		),
	)
}

// ValidatePatch ...
func (c *Column) ValidatePatch() error {
	return validation.ValidateStruct(
		c,
		validation.Field(
			&c.Name,
			validation.Length(1, 255),
		),
	)
}
