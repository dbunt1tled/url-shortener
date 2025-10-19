package e

import "github.com/cloudwego/hertz/pkg/protocol/consts"

type ErrNo struct {
	Msg    string `json:"message"`
	Code   int    `json:"code"`
	Status int    `json:"status"`
}

type HTTPError interface {
	Error() string
}

func NewErrNo(msg string, code int, status int) *ErrNo {
	return &ErrNo{
		Msg:    msg,
		Code:   code,
		Status: status,
	}
}

func NewNotFoundError(msg string, status int) HTTPError {
	return NewErrNo(msg, consts.StatusNotFound, status)
}

func NewBadRequestError(msg string, status int) HTTPError {
	return NewErrNo(msg, consts.StatusBadRequest, status)
}

func NewInternalError(msg string, status int) HTTPError {
	return NewErrNo(msg, consts.StatusInternalServerError, status)
}

func NewForbiddenError(msg string, status int) HTTPError {
	return NewErrNo(msg, consts.StatusForbidden, status)
}

func NewUnauthorizedError(msg string, status int) HTTPError {
	return NewErrNo(msg, consts.StatusUnauthorized, status)
}

func NewValidationError(msg string, status int) HTTPError {
	return NewErrNo(msg, consts.StatusUnprocessableEntity, status)
}

func NewUnprocessableEntityError(msg string, code int) HTTPError {
	return NewErrNo(msg, code, consts.StatusUnprocessableEntity)
}

func (e ErrNo) Error() string {
	return e.Msg
}
