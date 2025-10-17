package middleware

import (
	"context"
	"errors"
	"net/http"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/dbunt1tled/url-shortener/internal/config"
	"github.com/dbunt1tled/url-shortener/internal/lib/http-server/e"
	"github.com/dbunt1tled/url-shortener/internal/lib/locale"
)

func ErrorHandler(logger *config.AppLogger) app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		c.Next(ctx)

		err := c.Errors.Last()
		if err == nil {
			return
		}

		var er *e.DomainError
		switch {
		case errors.As(err.Err, &er):
			logger.Error(er.Error(), er)
			c.JSON(er.Code, map[string]any{
				"error":  locale.LCtx(c, er.Msg, nil),
				"status": er.Status,
			})
		default:
			logger.Error(er.Error(), er)
			c.JSON(http.StatusInternalServerError, map[string]any{
				"error": locale.LCtx(c, "error.server_error", nil),
			})
		}
	}
}
