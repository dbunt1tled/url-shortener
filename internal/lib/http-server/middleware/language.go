package middleware

import (
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
	"net/http"
)

func Language(locale *i18n.Localizer, bundle *i18n.Bundle) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			var l language.Tag
			var err error
			l = language.English
			lang := r.FormValue("lang")
			if lang != "" {
				l, err = language.Parse(lang)
				if err != nil {
					lang = r.Header.Get("Accept-Language")
					l, err = language.Parse(lang)
					if err != nil {
						l = language.English
					}
				}
			}
			locale = i18n.NewLocalizer(bundle, l.String())
			next.ServeHTTP(w, r)
		}

		return http.HandlerFunc(fn)
	}

}
