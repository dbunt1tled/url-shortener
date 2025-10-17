package validator

import (
	"fmt"
	"strings"

	"errors"

	"github.com/bytedance/go-tagexpr/v2/validator"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server/binding"
	"github.com/dbunt1tled/url-shortener/internal/lib/locale"
)

const (
	delimeter = ": "
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
		return e.FailPath + delimeter + e.Msg
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

func LValidate(c *app.RequestContext, v interface{}) error {
	var (
		err     error
		errStr  strings.Builder
		e, eMsg []string
	)
	if err = c.BindAndValidate(v); err == nil {
		return nil
	}

	e = strings.Split(err.Error(), "\t")
	for _, ev := range e {
		eMsg = strings.Split(ev, delimeter)
		errStr.WriteString(fmt.Sprintf("%s\n", locale.LCtx(c, eMsg[1], locale.M{"Field": eMsg[0]})))
	}

	return errors.New(strings.TrimSuffix(errStr.String(), "\n"))
}
