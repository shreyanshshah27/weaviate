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

// ObjectsClassReferencesPutReader is a Reader for the ObjectsClassReferencesPut structure.
type ObjectsClassReferencesPutReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *ObjectsClassReferencesPutReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := NewObjectsClassReferencesPutOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	case 400:
		result := NewObjectsClassReferencesPutBadRequest()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	case 401:
		result := NewObjectsClassReferencesPutUnauthorized()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	case 403:
		result := NewObjectsClassReferencesPutForbidden()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	case 404:
		result := NewObjectsClassReferencesPutNotFound()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	case 422:
		result := NewObjectsClassReferencesPutUnprocessableEntity()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	case 500:
		result := NewObjectsClassReferencesPutInternalServerError()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result

	default:
		return nil, runtime.NewAPIError("unknown error", response, response.Code())
	}
}

// NewObjectsClassReferencesPutOK creates a ObjectsClassReferencesPutOK with default headers values
func NewObjectsClassReferencesPutOK() *ObjectsClassReferencesPutOK {
	return &ObjectsClassReferencesPutOK{}
}

/*
ObjectsClassReferencesPutOK handles this case with default header values.

Successfully replaced all the references.
*/
type ObjectsClassReferencesPutOK struct {
}

func (o *ObjectsClassReferencesPutOK) Error() string {
	return fmt.Sprintf("[PUT /objects/{className}/{id}/references/{propertyName}][%d] objectsClassReferencesPutOK ", 200)
}

func (o *ObjectsClassReferencesPutOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	return nil
}

// NewObjectsClassReferencesPutBadRequest creates a ObjectsClassReferencesPutBadRequest with default headers values
func NewObjectsClassReferencesPutBadRequest() *ObjectsClassReferencesPutBadRequest {
	return &ObjectsClassReferencesPutBadRequest{}
}

/*
ObjectsClassReferencesPutBadRequest handles this case with default header values.

Malformed request.
*/
type ObjectsClassReferencesPutBadRequest struct {
	Payload *models.ErrorResponse
}

func (o *ObjectsClassReferencesPutBadRequest) Error() string {
	return fmt.Sprintf("[PUT /objects/{className}/{id}/references/{propertyName}][%d] objectsClassReferencesPutBadRequest  %+v", 400, o.Payload)
}

func (o *ObjectsClassReferencesPutBadRequest) GetPayload() *models.ErrorResponse {
	return o.Payload
}

func (o *ObjectsClassReferencesPutBadRequest) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.ErrorResponse)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewObjectsClassReferencesPutUnauthorized creates a ObjectsClassReferencesPutUnauthorized with default headers values
func NewObjectsClassReferencesPutUnauthorized() *ObjectsClassReferencesPutUnauthorized {
	return &ObjectsClassReferencesPutUnauthorized{}
}

/*
ObjectsClassReferencesPutUnauthorized handles this case with default header values.

Unauthorized or invalid credentials.
*/
type ObjectsClassReferencesPutUnauthorized struct {
}

func (o *ObjectsClassReferencesPutUnauthorized) Error() string {
	return fmt.Sprintf("[PUT /objects/{className}/{id}/references/{propertyName}][%d] objectsClassReferencesPutUnauthorized ", 401)
}

func (o *ObjectsClassReferencesPutUnauthorized) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	return nil
}

// NewObjectsClassReferencesPutForbidden creates a ObjectsClassReferencesPutForbidden with default headers values
func NewObjectsClassReferencesPutForbidden() *ObjectsClassReferencesPutForbidden {
	return &ObjectsClassReferencesPutForbidden{}
}

/*
ObjectsClassReferencesPutForbidden handles this case with default header values.

Forbidden
*/
type ObjectsClassReferencesPutForbidden struct {
	Payload *models.ErrorResponse
}

func (o *ObjectsClassReferencesPutForbidden) Error() string {
	return fmt.Sprintf("[PUT /objects/{className}/{id}/references/{propertyName}][%d] objectsClassReferencesPutForbidden  %+v", 403, o.Payload)
}

func (o *ObjectsClassReferencesPutForbidden) GetPayload() *models.ErrorResponse {
	return o.Payload
}

func (o *ObjectsClassReferencesPutForbidden) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.ErrorResponse)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewObjectsClassReferencesPutNotFound creates a ObjectsClassReferencesPutNotFound with default headers values
func NewObjectsClassReferencesPutNotFound() *ObjectsClassReferencesPutNotFound {
	return &ObjectsClassReferencesPutNotFound{}
}

/*
ObjectsClassReferencesPutNotFound handles this case with default header values.

Source object doesn't exist.
*/
type ObjectsClassReferencesPutNotFound struct {
}

func (o *ObjectsClassReferencesPutNotFound) Error() string {
	return fmt.Sprintf("[PUT /objects/{className}/{id}/references/{propertyName}][%d] objectsClassReferencesPutNotFound ", 404)
}

func (o *ObjectsClassReferencesPutNotFound) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	return nil
}

// NewObjectsClassReferencesPutUnprocessableEntity creates a ObjectsClassReferencesPutUnprocessableEntity with default headers values
func NewObjectsClassReferencesPutUnprocessableEntity() *ObjectsClassReferencesPutUnprocessableEntity {
	return &ObjectsClassReferencesPutUnprocessableEntity{}
}

/*
ObjectsClassReferencesPutUnprocessableEntity handles this case with default header values.

Request body is well-formed (i.e., syntactically correct), but semantically erroneous. Are you sure the property exists or that it is a class?
*/
type ObjectsClassReferencesPutUnprocessableEntity struct {
	Payload *models.ErrorResponse
}

func (o *ObjectsClassReferencesPutUnprocessableEntity) Error() string {
	return fmt.Sprintf("[PUT /objects/{className}/{id}/references/{propertyName}][%d] objectsClassReferencesPutUnprocessableEntity  %+v", 422, o.Payload)
}

func (o *ObjectsClassReferencesPutUnprocessableEntity) GetPayload() *models.ErrorResponse {
	return o.Payload
}

func (o *ObjectsClassReferencesPutUnprocessableEntity) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.ErrorResponse)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewObjectsClassReferencesPutInternalServerError creates a ObjectsClassReferencesPutInternalServerError with default headers values
func NewObjectsClassReferencesPutInternalServerError() *ObjectsClassReferencesPutInternalServerError {
	return &ObjectsClassReferencesPutInternalServerError{}
}

/*
ObjectsClassReferencesPutInternalServerError handles this case with default header values.

An error has occurred while trying to fulfill the request. Most likely the ErrorResponse will contain more information about the error.
*/
type ObjectsClassReferencesPutInternalServerError struct {
	Payload *models.ErrorResponse
}

func (o *ObjectsClassReferencesPutInternalServerError) Error() string {
	return fmt.Sprintf("[PUT /objects/{className}/{id}/references/{propertyName}][%d] objectsClassReferencesPutInternalServerError  %+v", 500, o.Payload)
}

func (o *ObjectsClassReferencesPutInternalServerError) GetPayload() *models.ErrorResponse {
	return o.Payload
}

func (o *ObjectsClassReferencesPutInternalServerError) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.ErrorResponse)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}
