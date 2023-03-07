// Code generated by go-swagger; DO NOT EDIT.

package k8s_monitoring_tool

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	"k8s-monitoring-tool/models"
)

// GetV1alpha1PodsLogOKCode is the HTTP code returned for type GetV1alpha1PodsLogOK
const GetV1alpha1PodsLogOKCode int = 200

/*GetV1alpha1PodsLogOK Fetching of all pod status is successful.

swagger:response getV1alpha1PodsLogOK
*/
type GetV1alpha1PodsLogOK struct {

	/*
	  In: Body
	*/
	Payload string `json:"body,omitempty"`
}

// NewGetV1alpha1PodsLogOK creates GetV1alpha1PodsLogOK with default headers values
func NewGetV1alpha1PodsLogOK() *GetV1alpha1PodsLogOK {

	return &GetV1alpha1PodsLogOK{}
}

// WithPayload adds the payload to the get v1alpha1 pods log o k response
func (o *GetV1alpha1PodsLogOK) WithPayload(payload string) *GetV1alpha1PodsLogOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get v1alpha1 pods log o k response
func (o *GetV1alpha1PodsLogOK) SetPayload(payload string) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetV1alpha1PodsLogOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(200)
	payload := o.Payload
	if err := producer.Produce(rw, payload); err != nil {
		panic(err) // let the recovery middleware deal with this
	}
}

// GetV1alpha1PodsLogInternalServerErrorCode is the HTTP code returned for type GetV1alpha1PodsLogInternalServerError
const GetV1alpha1PodsLogInternalServerErrorCode int = 500

/*GetV1alpha1PodsLogInternalServerError Internal server error

swagger:response getV1alpha1PodsLogInternalServerError
*/
type GetV1alpha1PodsLogInternalServerError struct {

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewGetV1alpha1PodsLogInternalServerError creates GetV1alpha1PodsLogInternalServerError with default headers values
func NewGetV1alpha1PodsLogInternalServerError() *GetV1alpha1PodsLogInternalServerError {

	return &GetV1alpha1PodsLogInternalServerError{}
}

// WithPayload adds the payload to the get v1alpha1 pods log internal server error response
func (o *GetV1alpha1PodsLogInternalServerError) WithPayload(payload *models.Error) *GetV1alpha1PodsLogInternalServerError {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get v1alpha1 pods log internal server error response
func (o *GetV1alpha1PodsLogInternalServerError) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetV1alpha1PodsLogInternalServerError) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(500)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}