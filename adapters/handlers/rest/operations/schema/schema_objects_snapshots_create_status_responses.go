//                           _       _
// __      _____  __ ___   ___  __ _| |_ ___
// \ \ /\ / / _ \/ _` \ \ / / |/ _` | __/ _ \
//  \ V  V /  __/ (_| |\ V /| | (_| | ||  __/
//   \_/\_/ \___|\__,_| \_/ |_|\__,_|\__\___|
//
//  Copyright © 2016 - 2022 SeMI Technologies B.V. All rights reserved.
//
//  CONTACT: hello@semi.technology
//

// Code generated by go-swagger; DO NOT EDIT.

package schema

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	"github.com/semi-technologies/weaviate/entities/models"
)

// SchemaObjectsSnapshotsCreateStatusOKCode is the HTTP code returned for type SchemaObjectsSnapshotsCreateStatusOK
const SchemaObjectsSnapshotsCreateStatusOKCode int = 200

/*SchemaObjectsSnapshotsCreateStatusOK Snapshot creation status successfully returned

swagger:response schemaObjectsSnapshotsCreateStatusOK
*/
type SchemaObjectsSnapshotsCreateStatusOK struct {

	/*
	  In: Body
	*/
	Payload *models.SnapshotMeta `json:"body,omitempty"`
}

// NewSchemaObjectsSnapshotsCreateStatusOK creates SchemaObjectsSnapshotsCreateStatusOK with default headers values
func NewSchemaObjectsSnapshotsCreateStatusOK() *SchemaObjectsSnapshotsCreateStatusOK {

	return &SchemaObjectsSnapshotsCreateStatusOK{}
}

// WithPayload adds the payload to the schema objects snapshots create status o k response
func (o *SchemaObjectsSnapshotsCreateStatusOK) WithPayload(payload *models.SnapshotMeta) *SchemaObjectsSnapshotsCreateStatusOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the schema objects snapshots create status o k response
func (o *SchemaObjectsSnapshotsCreateStatusOK) SetPayload(payload *models.SnapshotMeta) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *SchemaObjectsSnapshotsCreateStatusOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(200)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// SchemaObjectsSnapshotsCreateStatusUnauthorizedCode is the HTTP code returned for type SchemaObjectsSnapshotsCreateStatusUnauthorized
const SchemaObjectsSnapshotsCreateStatusUnauthorizedCode int = 401

/*SchemaObjectsSnapshotsCreateStatusUnauthorized Unauthorized or invalid credentials.

swagger:response schemaObjectsSnapshotsCreateStatusUnauthorized
*/
type SchemaObjectsSnapshotsCreateStatusUnauthorized struct {
}

// NewSchemaObjectsSnapshotsCreateStatusUnauthorized creates SchemaObjectsSnapshotsCreateStatusUnauthorized with default headers values
func NewSchemaObjectsSnapshotsCreateStatusUnauthorized() *SchemaObjectsSnapshotsCreateStatusUnauthorized {

	return &SchemaObjectsSnapshotsCreateStatusUnauthorized{}
}

// WriteResponse to the client
func (o *SchemaObjectsSnapshotsCreateStatusUnauthorized) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.Header().Del(runtime.HeaderContentType) //Remove Content-Type on empty responses

	rw.WriteHeader(401)
}

// SchemaObjectsSnapshotsCreateStatusForbiddenCode is the HTTP code returned for type SchemaObjectsSnapshotsCreateStatusForbidden
const SchemaObjectsSnapshotsCreateStatusForbiddenCode int = 403

/*SchemaObjectsSnapshotsCreateStatusForbidden Forbidden

swagger:response schemaObjectsSnapshotsCreateStatusForbidden
*/
type SchemaObjectsSnapshotsCreateStatusForbidden struct {

	/*
	  In: Body
	*/
	Payload *models.ErrorResponse `json:"body,omitempty"`
}

// NewSchemaObjectsSnapshotsCreateStatusForbidden creates SchemaObjectsSnapshotsCreateStatusForbidden with default headers values
func NewSchemaObjectsSnapshotsCreateStatusForbidden() *SchemaObjectsSnapshotsCreateStatusForbidden {

	return &SchemaObjectsSnapshotsCreateStatusForbidden{}
}

// WithPayload adds the payload to the schema objects snapshots create status forbidden response
func (o *SchemaObjectsSnapshotsCreateStatusForbidden) WithPayload(payload *models.ErrorResponse) *SchemaObjectsSnapshotsCreateStatusForbidden {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the schema objects snapshots create status forbidden response
func (o *SchemaObjectsSnapshotsCreateStatusForbidden) SetPayload(payload *models.ErrorResponse) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *SchemaObjectsSnapshotsCreateStatusForbidden) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(403)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// SchemaObjectsSnapshotsCreateStatusNotFoundCode is the HTTP code returned for type SchemaObjectsSnapshotsCreateStatusNotFound
const SchemaObjectsSnapshotsCreateStatusNotFoundCode int = 404

/*SchemaObjectsSnapshotsCreateStatusNotFound Not Found - Snapshot does not exist

swagger:response schemaObjectsSnapshotsCreateStatusNotFound
*/
type SchemaObjectsSnapshotsCreateStatusNotFound struct {
}

// NewSchemaObjectsSnapshotsCreateStatusNotFound creates SchemaObjectsSnapshotsCreateStatusNotFound with default headers values
func NewSchemaObjectsSnapshotsCreateStatusNotFound() *SchemaObjectsSnapshotsCreateStatusNotFound {

	return &SchemaObjectsSnapshotsCreateStatusNotFound{}
}

// WriteResponse to the client
func (o *SchemaObjectsSnapshotsCreateStatusNotFound) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.Header().Del(runtime.HeaderContentType) //Remove Content-Type on empty responses

	rw.WriteHeader(404)
}

// SchemaObjectsSnapshotsCreateStatusInternalServerErrorCode is the HTTP code returned for type SchemaObjectsSnapshotsCreateStatusInternalServerError
const SchemaObjectsSnapshotsCreateStatusInternalServerErrorCode int = 500

/*SchemaObjectsSnapshotsCreateStatusInternalServerError An error has occurred while trying to fulfill the request. Most likely the ErrorResponse will contain more information about the error.

swagger:response schemaObjectsSnapshotsCreateStatusInternalServerError
*/
type SchemaObjectsSnapshotsCreateStatusInternalServerError struct {

	/*
	  In: Body
	*/
	Payload *models.ErrorResponse `json:"body,omitempty"`
}

// NewSchemaObjectsSnapshotsCreateStatusInternalServerError creates SchemaObjectsSnapshotsCreateStatusInternalServerError with default headers values
func NewSchemaObjectsSnapshotsCreateStatusInternalServerError() *SchemaObjectsSnapshotsCreateStatusInternalServerError {

	return &SchemaObjectsSnapshotsCreateStatusInternalServerError{}
}

// WithPayload adds the payload to the schema objects snapshots create status internal server error response
func (o *SchemaObjectsSnapshotsCreateStatusInternalServerError) WithPayload(payload *models.ErrorResponse) *SchemaObjectsSnapshotsCreateStatusInternalServerError {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the schema objects snapshots create status internal server error response
func (o *SchemaObjectsSnapshotsCreateStatusInternalServerError) SetPayload(payload *models.ErrorResponse) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *SchemaObjectsSnapshotsCreateStatusInternalServerError) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(500)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}
