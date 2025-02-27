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

package backups

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"fmt"
	"io"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"

	"github.com/weaviate/weaviate/entities/models"
)

// BackupsRestoreStatusReader is a Reader for the BackupsRestoreStatus structure.
type BackupsRestoreStatusReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *BackupsRestoreStatusReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := NewBackupsRestoreStatusOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	case 401:
		result := NewBackupsRestoreStatusUnauthorized()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	case 403:
		result := NewBackupsRestoreStatusForbidden()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	case 404:
		result := NewBackupsRestoreStatusNotFound()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	case 500:
		result := NewBackupsRestoreStatusInternalServerError()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result

	default:
		return nil, runtime.NewAPIError("unknown error", response, response.Code())
	}
}

// NewBackupsRestoreStatusOK creates a BackupsRestoreStatusOK with default headers values
func NewBackupsRestoreStatusOK() *BackupsRestoreStatusOK {
	return &BackupsRestoreStatusOK{}
}

/*
BackupsRestoreStatusOK handles this case with default header values.

Backup restoration status successfully returned
*/
type BackupsRestoreStatusOK struct {
	Payload *models.BackupRestoreStatusResponse
}

func (o *BackupsRestoreStatusOK) Error() string {
	return fmt.Sprintf("[GET /backups/{backend}/{id}/restore][%d] backupsRestoreStatusOK  %+v", 200, o.Payload)
}

func (o *BackupsRestoreStatusOK) GetPayload() *models.BackupRestoreStatusResponse {
	return o.Payload
}

func (o *BackupsRestoreStatusOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.BackupRestoreStatusResponse)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewBackupsRestoreStatusUnauthorized creates a BackupsRestoreStatusUnauthorized with default headers values
func NewBackupsRestoreStatusUnauthorized() *BackupsRestoreStatusUnauthorized {
	return &BackupsRestoreStatusUnauthorized{}
}

/*
BackupsRestoreStatusUnauthorized handles this case with default header values.

Unauthorized or invalid credentials.
*/
type BackupsRestoreStatusUnauthorized struct {
}

func (o *BackupsRestoreStatusUnauthorized) Error() string {
	return fmt.Sprintf("[GET /backups/{backend}/{id}/restore][%d] backupsRestoreStatusUnauthorized ", 401)
}

func (o *BackupsRestoreStatusUnauthorized) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	return nil
}

// NewBackupsRestoreStatusForbidden creates a BackupsRestoreStatusForbidden with default headers values
func NewBackupsRestoreStatusForbidden() *BackupsRestoreStatusForbidden {
	return &BackupsRestoreStatusForbidden{}
}

/*
BackupsRestoreStatusForbidden handles this case with default header values.

Forbidden
*/
type BackupsRestoreStatusForbidden struct {
	Payload *models.ErrorResponse
}

func (o *BackupsRestoreStatusForbidden) Error() string {
	return fmt.Sprintf("[GET /backups/{backend}/{id}/restore][%d] backupsRestoreStatusForbidden  %+v", 403, o.Payload)
}

func (o *BackupsRestoreStatusForbidden) GetPayload() *models.ErrorResponse {
	return o.Payload
}

func (o *BackupsRestoreStatusForbidden) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.ErrorResponse)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewBackupsRestoreStatusNotFound creates a BackupsRestoreStatusNotFound with default headers values
func NewBackupsRestoreStatusNotFound() *BackupsRestoreStatusNotFound {
	return &BackupsRestoreStatusNotFound{}
}

/*
BackupsRestoreStatusNotFound handles this case with default header values.

Not Found - Backup does not exist
*/
type BackupsRestoreStatusNotFound struct {
	Payload *models.ErrorResponse
}

func (o *BackupsRestoreStatusNotFound) Error() string {
	return fmt.Sprintf("[GET /backups/{backend}/{id}/restore][%d] backupsRestoreStatusNotFound  %+v", 404, o.Payload)
}

func (o *BackupsRestoreStatusNotFound) GetPayload() *models.ErrorResponse {
	return o.Payload
}

func (o *BackupsRestoreStatusNotFound) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.ErrorResponse)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewBackupsRestoreStatusInternalServerError creates a BackupsRestoreStatusInternalServerError with default headers values
func NewBackupsRestoreStatusInternalServerError() *BackupsRestoreStatusInternalServerError {
	return &BackupsRestoreStatusInternalServerError{}
}

/*
BackupsRestoreStatusInternalServerError handles this case with default header values.

An error has occurred while trying to fulfill the request. Most likely the ErrorResponse will contain more information about the error.
*/
type BackupsRestoreStatusInternalServerError struct {
	Payload *models.ErrorResponse
}

func (o *BackupsRestoreStatusInternalServerError) Error() string {
	return fmt.Sprintf("[GET /backups/{backend}/{id}/restore][%d] backupsRestoreStatusInternalServerError  %+v", 500, o.Payload)
}

func (o *BackupsRestoreStatusInternalServerError) GetPayload() *models.ErrorResponse {
	return o.Payload
}

func (o *BackupsRestoreStatusInternalServerError) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.ErrorResponse)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}
