package middleware

import (
	"context"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

func LocaleMiddleware(bundle *i18n.Bundle) app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		lang := string(c.Request.Header.Peek("Accept-Language"))
		if lang == "" {
			lang = c.Param("lang")
			if lang == "" {
				lang = "en"
			}
		}
		l, err := language.Parse(lang)
		if err != nil {
			l = language.English
		}
		locale := i18n.NewLocalizer(bundle, l.String(), language.English.String())
		c.Set("localizer", locale)
		c.Next(ctx)
	}
}
