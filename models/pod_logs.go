// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"

	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
)

// PodLogs pod logs
//
// swagger:model PodLogs
type PodLogs struct {

	// name.
	Name string `json:"name,omitempty"`

	// output.
	Output string `json:"output,omitempty"`
}

// Validate validates this pod logs
func (m *PodLogs) Validate(formats strfmt.Registry) error {
	return nil
}

// ContextValidate validates this pod logs based on context it is used
func (m *PodLogs) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (m *PodLogs) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *PodLogs) UnmarshalBinary(b []byte) error {
	var res PodLogs
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
