package model

import validation "github.com/go-ozzo/ozzo-validation"

// Task ...
type Task struct {
	ID          *string `json:"id,omitempty"`
	ColumnID    *string `json:"column_id,omitempty"`
	Name        *string `json:"name,omitempty"`
	Description *string `json:"description,omitempty"`
	Position    *int    `json:"position,omitempty"`
}

// Validate ...
func (t *Task) Validate() error {
	return validation.ValidateStruct(
		t,
		validation.Field(
			&t.Name,
			validation.Required,
			validation.Length(1, 500),
		),
		validation.Field(
			&t.Description,
			validation.Required,
			validation.Length(0, 5000),
		),
		validation.Field(
			&t.ColumnID,
			validation.Required,
		),
	)
}

// ValidatePatch ...
func (t *Task) ValidatePatch() error {
	return validation.ValidateStruct(
		t,
		validation.Field(
			&t.Name,
			validation.Length(1, 500),
		),
		validation.Field(
			&t.Description,
			validation.Length(0, 5000),
		),
	)
}
