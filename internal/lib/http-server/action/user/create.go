package user

import (
	"encoding/json"
	"github.com/go-playground/validator/v10"
	"github.com/google/jsonapi"
	"go_first/internal/lib/common"
	"go_first/internal/lib/common/model/user"
	"go_first/storage"
	"log/slog"
	"net/http"
	"strings"
)

type Request struct {
	FirstName   string `json:"firstName" validate:"required" validateMsg:"A name is required"`
	SecondName  string `json:"secondName" validate:"required"`
	Email       string `json:"email" validate:"required,email"`
	PhoneNumber string `json:"phoneNumber" validate:"required"`
	Password    string `json:"password" validate:"required"`
}

func CreateUserAction(log *slog.Logger, storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var request Request
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			http.Error(w, err.Error(), http.StatusUnprocessableEntity)
			return
		}
		if err := validator.New().Struct(request); err != nil {
			http.Error(w, "Error validate request"+err.Error()+"\n"+strings.Join(common.ValidationErrorString(err.(validator.ValidationErrors)), "\n"), 422)
			return
		}
		u, err := storage.CreateUser(
			request.FirstName,
			request.SecondName,
			request.Email,
			request.PhoneNumber,
			request.Password,
			user.StatusActive,
		)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if err := jsonapi.MarshalPayload(w, u); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}
