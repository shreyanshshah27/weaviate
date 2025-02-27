//                           _       _
// __      _____  __ ___   ___  __ _| |_ ___
// \ \ /\ / / _ \/ _` \ \ / / |/ _` | __/ _ \
//  \ V  V /  __/ (_| |\ V /| | (_| | ||  __/
//   \_/\_/ \___|\__,_| \_/ |_|\__,_|\__\___|
//
//  Copyright © 2016 - 2023 Weaviate B.V. All rights reserved.
//
//  CONTACT: hello@weaviate.io
//

// Code generated by go-swagger; DO NOT EDIT.

package objects

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"fmt"
	"io"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"

	"github.com/weaviate/weaviate/entities/models"
)

// ObjectsReferencesDeleteReader is a Reader for the ObjectsReferencesDelete structure.
type ObjectsReferencesDeleteReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *ObjectsReferencesDeleteReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 204:
		result := NewObjectsReferencesDeleteNoContent()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	case 401:
		result := NewObjectsReferencesDeleteUnauthorized()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	case 403:
		result := NewObjectsReferencesDeleteForbidden()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	case 404:
		result := NewObjectsReferencesDeleteNotFound()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	case 500:
		result := NewObjectsReferencesDeleteInternalServerError()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result

	default:
		return nil, runtime.NewAPIError("unknown error", response, response.Code())
	}
}

// NewObjectsReferencesDeleteNoContent creates a ObjectsReferencesDeleteNoContent with default headers values
func NewObjectsReferencesDeleteNoContent() *ObjectsReferencesDeleteNoContent {
	return &ObjectsReferencesDeleteNoContent{}
}

/*
ObjectsReferencesDeleteNoContent handles this case with default header values.

Successfully deleted.
*/
type ObjectsReferencesDeleteNoContent struct {
}

func (o *ObjectsReferencesDeleteNoContent) Error() string {
	return fmt.Sprintf("[DELETE /objects/{id}/references/{propertyName}][%d] objectsReferencesDeleteNoContent ", 204)
}

func (o *ObjectsReferencesDeleteNoContent) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	return nil
}

// NewObjectsReferencesDeleteUnauthorized creates a ObjectsReferencesDeleteUnauthorized with default headers values
func NewObjectsReferencesDeleteUnauthorized() *ObjectsReferencesDeleteUnauthorized {
	return &ObjectsReferencesDeleteUnauthorized{}
}

/*
ObjectsReferencesDeleteUnauthorized handles this case with default header values.

Unauthorized or invalid credentials.
*/
type ObjectsReferencesDeleteUnauthorized struct {
}

func (o *ObjectsReferencesDeleteUnauthorized) Error() string {
	return fmt.Sprintf("[DELETE /objects/{id}/references/{propertyName}][%d] objectsReferencesDeleteUnauthorized ", 401)
}

func (o *ObjectsReferencesDeleteUnauthorized) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	return nil
}

// NewObjectsReferencesDeleteForbidden creates a ObjectsReferencesDeleteForbidden with default headers values
func NewObjectsReferencesDeleteForbidden() *ObjectsReferencesDeleteForbidden {
	return &ObjectsReferencesDeleteForbidden{}
}

/*
ObjectsReferencesDeleteForbidden handles this case with default header values.

Forbidden
*/
type ObjectsReferencesDeleteForbidden struct {
	Payload *models.ErrorResponse
}

func (o *ObjectsReferencesDeleteForbidden) Error() string {
	return fmt.Sprintf("[DELETE /objects/{id}/references/{propertyName}][%d] objectsReferencesDeleteForbidden  %+v", 403, o.Payload)
}

func (o *ObjectsReferencesDeleteForbidden) GetPayload() *models.ErrorResponse {
	return o.Payload
}

func (o *ObjectsReferencesDeleteForbidden) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.ErrorResponse)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewObjectsReferencesDeleteNotFound creates a ObjectsReferencesDeleteNotFound with default headers values
func NewObjectsReferencesDeleteNotFound() *ObjectsReferencesDeleteNotFound {
	return &ObjectsReferencesDeleteNotFound{}
}

/*
ObjectsReferencesDeleteNotFound handles this case with default header values.

Successful query result but no resource was found.
*/
type ObjectsReferencesDeleteNotFound struct {
	Payload *models.ErrorResponse
}

func (o *ObjectsReferencesDeleteNotFound) Error() string {
	return fmt.Sprintf("[DELETE /objects/{id}/references/{propertyName}][%d] objectsReferencesDeleteNotFound  %+v", 404, o.Payload)
}

func (o *ObjectsReferencesDeleteNotFound) GetPayload() *models.ErrorResponse {
	return o.Payload
}

func (o *ObjectsReferencesDeleteNotFound) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.ErrorResponse)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewObjectsReferencesDeleteInternalServerError creates a ObjectsReferencesDeleteInternalServerError with default headers values
func NewObjectsReferencesDeleteInternalServerError() *ObjectsReferencesDeleteInternalServerError {
	return &ObjectsReferencesDeleteInternalServerError{}
}

/*
ObjectsReferencesDeleteInternalServerError handles this case with default header values.

An error has occurred while trying to fulfill the request. Most likely the ErrorResponse will contain more information about the error.
*/
type ObjectsReferencesDeleteInternalServerError struct {
	Payload *models.ErrorResponse
}

func (o *ObjectsReferencesDeleteInternalServerError) Error() string {
	return fmt.Sprintf("[DELETE /objects/{id}/references/{propertyName}][%d] objectsReferencesDeleteInternalServerError  %+v", 500, o.Payload)
}

func (o *ObjectsReferencesDeleteInternalServerError) GetPayload() *models.ErrorResponse {
	return o.Payload
}

func (o *ObjectsReferencesDeleteInternalServerError) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.ErrorResponse)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}
