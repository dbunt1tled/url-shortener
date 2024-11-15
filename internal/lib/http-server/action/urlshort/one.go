package urlshort

import (
	"github.com/go-chi/chi/v5"
	"go_first/storage"
	"go_first/storage/mysql"
	"net/http"
)

func GetUrlAction(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		alias := chi.URLParam(r, "alias")
		url, err := storage.GetURL(mysql.URLFilter{Alias: alias})
		if err != nil {
			http.Error(w, "", http.StatusNotFound)
			return
		}
		http.Redirect(w, r, url.URL, http.StatusFound)
	}
}
