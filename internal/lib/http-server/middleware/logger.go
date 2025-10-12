package middleware

import "C"
import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"github.com/dbunt1tled/url-shortener/internal/config"
)

func LoggerMiddleware(logger *config.AppLogger, levelStatus int) app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		t := time.Now()
		c.Next(ctx)

		if c.Response.StatusCode() < levelStatus {
			return
		}
		headers := make(map[string]string)

		c.Request.Header.VisitAll(func(k, v []byte) {
			headers[string(k)] = string(v)
		})
		logger.Log(
			ctx,
			parseLevel(c.Response.StatusCode()),
			fmt.Sprintf("[HTTP] request (%s) %s", c.Method(), c.FullPath()),
			slog.String("url", c.URI().String()),
			slog.Any("headers", headers),
			slog.String("remoteAddr", c.ClientIP()),
			slog.String("reqBody", string(c.Request.Body())),
			slog.String("userAgent", string(c.UserAgent())),
			slog.Int("status", c.Response.StatusCode()),
			slog.String("duration", time.Since(t).Round(time.Millisecond).String()),
			slog.String("resBody", string(c.Response.Body())),
		)
	}
}

func parseLevel(status int) slog.Level {
	switch {
	case status < consts.StatusMultipleChoices:
		return slog.LevelDebug
	case status >= consts.StatusMultipleChoices && status < consts.StatusBadRequest:
		return slog.LevelInfo
	case status >= consts.StatusBadRequest && status < consts.StatusInternalServerError:
		return slog.LevelWarn
	case status >= consts.StatusInternalServerError:
		return slog.LevelError
	default:
		return slog.LevelDebug
	}
}
