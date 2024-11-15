package apierror

import (
	"fmt"
	"github.com/google/jsonapi"
	"net/http"
	"runtime"
)

type APIError struct {
	ID     int         `json:"id,omitempty" jsonapi:"primary,error"`
	Status int         `json:"status,omitempty" jsonapi:"attr,status"`
	Errors []ErrorData `json:"errors,omitempty" jsonapi:"attr,errors"`
	Source Source      `json:"source,omitempty" jsonapi:"attr,source"`
}

type Source struct {
	Pointer string `json:"pointer,omitempty"  jsonapi:"attr,pointer"`
}

type ErrorData struct {
	Title  string `json:"title,omitempty" jsonapi:"attr,title"`
	Detail string `json:"detail,omitempty" jsonapi:"attr,detail"`
}

func newError(err error, code int, status int) *APIError {
	return newErrorString(err.Error(), code, status)
}

func newErrorString(err string, code int, status int) *APIError {
	return newErrorMap(map[string]interface{}{"api": err}, code, status)
}

func newErrorMap(err map[string]interface{}, code int, status int) *APIError {
	pc, _, line, _ := runtime.Caller(2)
	details := runtime.FuncForPC(pc)
	source := Source{Pointer: fmt.Sprintf("%s#%d", details.Name(), line)}
	data := make([]ErrorData, 0, len(err))
	for key, value := range err {
		data = append(data, ErrorData{Title: key, Detail: fmt.Sprintf("%v", value)})
	}
	return &APIError{
		ID:     code,
		Status: status,
		Errors: data,
		Source: source,
	}
}

func ErrorNotFound(w http.ResponseWriter, err error, code int) {
	w.WriteHeader(http.StatusNotFound)
	e := jsonapi.MarshalPayload(w, newError(err, code, http.StatusNotFound))
	if e != nil {
		http.Error(w, e.Error(), http.StatusNotFound)
	}
}
func ErrorInternal(w http.ResponseWriter, err error, code int) {
	w.WriteHeader(http.StatusInternalServerError)
	e := jsonapi.MarshalPayload(w, newError(err, code, http.StatusInternalServerError))
	if e != nil {
		http.Error(w, e.Error(), http.StatusInternalServerError)
	}
}
func ErrorUnauthorized(w http.ResponseWriter, err error, code int) {
	w.WriteHeader(http.StatusUnauthorized)
	e := jsonapi.MarshalPayload(w, newError(err, code, http.StatusUnauthorized))
	if e != nil {
		http.Error(w, e.Error(), http.StatusInternalServerError)
	}
}
func ErrorUnprocessableEntity(w http.ResponseWriter, err error, code int) {
	w.WriteHeader(http.StatusUnprocessableEntity)
	e := jsonapi.MarshalPayload(w, newError(err, code, http.StatusUnprocessableEntity))
	if e != nil {
		http.Error(w, e.Error(), http.StatusInternalServerError)
	}
}

func ErrorNotFoundMap(w http.ResponseWriter, err map[string]interface{}, code int) {
	w.WriteHeader(http.StatusNotFound)
	e := jsonapi.MarshalPayload(w, newErrorMap(err, code, http.StatusNotFound))
	if e != nil {
		http.Error(w, e.Error(), http.StatusNotFound)
	}
}
func ErrorInternalMap(w http.ResponseWriter, err map[string]interface{}, code int) {
	w.WriteHeader(http.StatusInternalServerError)
	e := jsonapi.MarshalPayload(w, newErrorMap(err, code, http.StatusInternalServerError))
	if e != nil {
		http.Error(w, e.Error(), http.StatusInternalServerError)
	}
}
func ErrorUnauthorizedMap(w http.ResponseWriter, err map[string]interface{}, code int) {
	w.WriteHeader(http.StatusUnauthorized)
	e := jsonapi.MarshalPayload(w, newErrorMap(err, code, http.StatusUnauthorized))
	if e != nil {
		http.Error(w, e.Error(), http.StatusInternalServerError)
	}
}
func ErrorUnprocessableEntityMap(w http.ResponseWriter, err map[string]interface{}, code int) {
	w.WriteHeader(http.StatusUnprocessableEntity)
	e := jsonapi.MarshalPayload(w, newErrorMap(err, code, http.StatusUnprocessableEntity))
	if e != nil {
		http.Error(w, e.Error(), http.StatusInternalServerError)
	}
}

func ErrorNotFoundString(w http.ResponseWriter, err string, code int) {
	w.WriteHeader(http.StatusNotFound)
	e := jsonapi.MarshalPayload(w, newErrorString(err, code, http.StatusNotFound))
	if e != nil {
		http.Error(w, e.Error(), http.StatusNotFound)
	}
}
func ErrorInternalString(w http.ResponseWriter, err string, code int) {
	w.WriteHeader(http.StatusInternalServerError)
	e := jsonapi.MarshalPayload(w, newErrorString(err, code, http.StatusInternalServerError))
	if e != nil {
		http.Error(w, e.Error(), http.StatusInternalServerError)
	}
}
func ErrorUnauthorizedString(w http.ResponseWriter, err string, code int) {
	w.WriteHeader(http.StatusUnauthorized)
	e := jsonapi.MarshalPayload(w, newErrorString(err, code, http.StatusUnauthorized))
	if e != nil {
		http.Error(w, e.Error(), http.StatusInternalServerError)
	}
}
func ErrorUnprocessableEntityString(w http.ResponseWriter, err string, code int) {
	w.WriteHeader(http.StatusUnprocessableEntity)
	e := jsonapi.MarshalPayload(w, newErrorString(err, code, http.StatusUnprocessableEntity))
	if e != nil {
		http.Error(w, e.Error(), http.StatusInternalServerError)
	}
}
