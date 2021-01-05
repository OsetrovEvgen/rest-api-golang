package model

import validation "github.com/go-ozzo/ozzo-validation"

// Project ...
type Project struct {
	ID          *string `json:"id,omitempty"`
	Name        *string `json:"name,omitempty"`
	Description *string `json:"description,omitempty"`
}

// Validate ...
func (p *Project) Validate() error {
	return validation.ValidateStruct(
		p,
		validation.Field(
			&p.Name,
			validation.Required,
			validation.Length(1, 500),
		),
		validation.Field(
			&p.Description,
			validation.Required,
			validation.Length(0, 1000),
		),
	)
}

// ValidatePatch ...
func (p *Project) ValidatePatch() error {
	return validation.ValidateStruct(
		p,
		validation.Field(
			&p.Name,
			validation.Length(1, 500),
		),
		validation.Field(
			&p.Description,
			validation.Length(0, 1000),
		),
	)
}
