package validator

import (
	"github.com/bytedance/go-tagexpr/v2/validator"
	"github.com/cloudwego/hertz/pkg/app/server/binding"
)

type CustomValidator struct {
	validator *validator.Validator
}

func New() binding.StructValidator {
	v := &CustomValidator{
		validator: validator.New("vd"),
	}
	v.validator.SetErrorFactory(defaultErrorFactory)
	return v
}

func (v *CustomValidator) Validate(i interface{}) error {
	return v.validator.Validate(i, true)
}

func (v *CustomValidator) ValidateStruct(i interface{}) error {
	return v.validator.Validate(i, true)
}

func (v *CustomValidator) Engine() interface{} {
	return v.validator
}

func (v *CustomValidator) ValidateTag() string {
	return "vd"
}

type Error struct {
	FailPath, Msg string
}

func (e *Error) Error() string {
	if e.Msg != "" {
		return e.FailPath + ": " + e.Msg
	}
	return e.FailPath + ": invalid"
}

//go:nosplit
func defaultErrorFactory(failPath, msg string) error {
	return &Error{
		FailPath: failPath,
		Msg:      msg,
	}
}
