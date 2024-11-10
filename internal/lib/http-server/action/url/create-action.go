package url

import (
	"encoding/json"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"go_first/internal/lib/common"
	response "go_first/internal/lib/common/http"
	"go_first/internal/lib/logger"
	"go_first/storage"
	"log/slog"
	"net/http"
)

type Request struct {
	URL string `json:"url" validate:"required,url"`
}

func createUrlAction(log *slog.Logger, storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var request Request
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			log.Error("Error decode request", logger.Error(err))
			render.JSON(w, r, response.Error(err.Error()))
			return
		}
		if err := validator.New().Struct(request); err != nil {
			log.Error("Error validate request", logger.Error(err))
			render.JSON(w, r, response.ValidationErrorResponse(err.(validator.ValidationErrors)))
			return
		}
		alias, err := common.RandStringBytes(32)
		if err != nil {
			log.Error("Error generate alias", logger.Error(err))
			render.JSON(w, r, response.Error(err.Error()))
			return
		}
		id, err := storage.CreateURL(request.URL, alias)

		if err != nil {
			log.Error("Error create url", logger.Error(err))
			render.JSON(w, r, response.Error(err.Error()))
			return
		}
		render.JSON(w, r, response.Ok(map[string]any{
			"id":    id,
			"alias": alias,
		}))
	}
}
