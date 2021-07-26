// Code generated by go-swagger; DO NOT EDIT.

package k8s_monitoring_tool

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	"k8s-monitoring-tool/models"
)

// GetV1alpha1PodsOKCode is the HTTP code returned for type GetV1alpha1PodsOK
const GetV1alpha1PodsOKCode int = 200

/*GetV1alpha1PodsOK Fetching of all pod status is successful.

swagger:response getV1alpha1PodsOK
*/
type GetV1alpha1PodsOK struct {

	/*
	  In: Body
	*/
	Payload models.Pods `json:"body,omitempty"`
}

// NewGetV1alpha1PodsOK creates GetV1alpha1PodsOK with default headers values
func NewGetV1alpha1PodsOK() *GetV1alpha1PodsOK {

	return &GetV1alpha1PodsOK{}
}

// WithPayload adds the payload to the get v1alpha1 pods o k response
func (o *GetV1alpha1PodsOK) WithPayload(payload models.Pods) *GetV1alpha1PodsOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get v1alpha1 pods o k response
func (o *GetV1alpha1PodsOK) SetPayload(payload models.Pods) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetV1alpha1PodsOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(200)
	payload := o.Payload
	if payload == nil {
		// return empty array
		payload = models.Pods{}
	}

	if err := producer.Produce(rw, payload); err != nil {
		panic(err) // let the recovery middleware deal with this
	}
}

// GetV1alpha1PodsInternalServerErrorCode is the HTTP code returned for type GetV1alpha1PodsInternalServerError
const GetV1alpha1PodsInternalServerErrorCode int = 500

/*GetV1alpha1PodsInternalServerError Internal server error

swagger:response getV1alpha1PodsInternalServerError
*/
type GetV1alpha1PodsInternalServerError struct {

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewGetV1alpha1PodsInternalServerError creates GetV1alpha1PodsInternalServerError with default headers values
func NewGetV1alpha1PodsInternalServerError() *GetV1alpha1PodsInternalServerError {

	return &GetV1alpha1PodsInternalServerError{}
}

// WithPayload adds the payload to the get v1alpha1 pods internal server error response
func (o *GetV1alpha1PodsInternalServerError) WithPayload(payload *models.Error) *GetV1alpha1PodsInternalServerError {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get v1alpha1 pods internal server error response
func (o *GetV1alpha1PodsInternalServerError) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetV1alpha1PodsInternalServerError) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(500)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}
