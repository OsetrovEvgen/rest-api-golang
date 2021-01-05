package model

import validation "github.com/go-ozzo/ozzo-validation"

// Comment ...
type Comment struct {
	ID     *string `json:"id,omitempty"`
	TaskID *string `json:"task_id,omitempty"`
	Text   *string `json:"text,omitempty"`
	Date   *string `json:"date,omitempty"`
}

// Validate ...
func (c *Comment) Validate() error {
	return validation.ValidateStruct(
		c,
		validation.Field(
			&c.Text,
			validation.Required,
			validation.Length(1, 500),
		),
	)
}

// ValidatePatch ...
func (c *Comment) ValidatePatch() error {
	return validation.ValidateStruct(
		c,
		validation.Field(
			&c.Text,
			validation.Length(1, 5000),
		),
	)
}
