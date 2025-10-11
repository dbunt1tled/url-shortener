package e

import "github.com/cloudwego/hertz/pkg/protocol/consts"

type DomainError struct {
	Msg    string `json:"message"`
	Code   int    `json:"code"`
	Status int    `json:"status"`
}

type HTTPError interface {
	Error() string
}

func NewDomainError(msg string, code int, status int) *DomainError {
	return &DomainError{
		Msg:    msg,
		Code:   code,
		Status: status,
	}
}

func NewNotFoundError(msg string, status int) HTTPError {
	return NewDomainError(msg, consts.StatusNotFound, status)
}

func NewBadRequestError(msg string, status int) HTTPError {
	return NewDomainError(msg, consts.StatusBadRequest, status)
}

func NewInternalError(msg string, status int) HTTPError {
	return NewDomainError(msg, consts.StatusInternalServerError, status)
}

func NewForbiddenError(msg string, status int) HTTPError {
	return NewDomainError(msg, consts.StatusForbidden, status)
}

func NewUnauthorizedError(msg string, status int) HTTPError {
	return NewDomainError(msg, consts.StatusUnauthorized, status)
}

func NewValidationError(msg string, status int) HTTPError {
	return NewDomainError(msg, consts.StatusUnprocessableEntity, status)
}

func NewUnprocessableEntityError(msg string, code int) HTTPError {
	return NewDomainError(msg, code, consts.StatusUnprocessableEntity)
}

func (e DomainError) Error() string {
	return e.Msg
}
