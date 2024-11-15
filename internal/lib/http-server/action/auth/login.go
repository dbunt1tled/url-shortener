package auth

import (
	"encoding/json"
	"github.com/go-playground/validator/v10"
	"github.com/google/jsonapi"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"go_first/internal/lib/common"
	"go_first/internal/lib/common/model/user"
	"go_first/internal/lib/common/password"
	"go_first/internal/lib/service/security"
	"go_first/storage"
	"log/slog"
	"net/http"
	"strings"
)

type Request struct {
	Login    string `json:"login" validate:"required"`
	Password string `json:"password" validate:"required"`
}

func LoginAction(
	log *slog.Logger,
	storage storage.Storage,
	locale *i18n.Localizer,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var request Request
		var err error
		var u *user.User
		if err = json.NewDecoder(r.Body).Decode(&request); err != nil {
			http.Error(w, err.Error(), http.StatusUnprocessableEntity)
			return
		}
		if err = validator.New().Struct(request); err != nil {
			http.Error(w, "Error validate request"+err.Error()+"\n"+strings.Join(common.ValidationErrorString(err.(validator.ValidationErrors)), "\n"), 422)
			return
		}
		u, err = storage.GetUserIdentity(request.Login)

		if err != nil {
			http.Error(w, locale.MustLocalize(&i18n.LocalizeConfig{
				MessageID: "login_or_password_wrong",
			}), http.StatusUnprocessableEntity)
			return
		}
		if u.Status != user.StatusActive {
			http.Error(w, locale.MustLocalize(&i18n.LocalizeConfig{
				MessageID: "login_or_password_wrong",
			}), http.StatusUnprocessableEntity)
			return
		}
		ok, err := password.Compare(request.Password, u.Password)
		if !ok || err != nil {
			http.Error(w, locale.MustLocalize(&i18n.LocalizeConfig{
				MessageID: "login_or_password_wrong",
			}), http.StatusUnprocessableEntity)
			return
		}
		tokens, err := security.GetAuthTokens(*u)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if err := jsonapi.MarshalPayload(w, tokens); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}
