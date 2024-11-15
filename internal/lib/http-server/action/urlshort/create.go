package urlshort

import (
	"encoding/json"
	"github.com/go-playground/validator/v10"
	"github.com/google/jsonapi"
	"go_first/internal/lib/common"
	"go_first/internal/lib/common/rand"
	"go_first/storage"
	"log/slog"
	"net/http"
	"strings"
)

type Request struct {
	URL string `json:"url" validate:"required,url"`
}

func CreateUrlAction(log *slog.Logger, storage storage.Storage) http.HandlerFunc {
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
		alias, err := rand.String(20)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		sURL, err := storage.CreateURL(request.URL, alias)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if err := jsonapi.MarshalPayload(w, sURL); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}
