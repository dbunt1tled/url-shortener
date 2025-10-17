package middleware

import (
	"context"
	"strings"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"github.com/dbunt1tled/url-shortener/internal/lib/hasher"
	"github.com/dbunt1tled/url-shortener/internal/lib/http-server/e"
	"github.com/dbunt1tled/url-shortener/internal/lib/locale"
)

const (
	BearerSchema = "Bearer"
	MsgErr       = "error.unauthorized"
)

func AuthBearerMiddleware(hasher *hasher.Hasher) app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		var authToken string
		var isEmpty bool
		authToken, isEmpty = fromAuthHeader(c)
		if isEmpty {
			authToken, isEmpty = fromQueryParam(c)
			if isEmpty {
				c.AbortWithStatusJSON(
					consts.StatusUnauthorized,
					map[string]any{
						"error":  locale.LCtx(c, MsgErr, nil),
						"status": e.Err401AuthEmptyTokenError,
					})
				return
			}
		}
		claims, err := hasher.Decode(authToken, false)
		if err != nil {
			c.AbortWithStatusJSON(
				consts.StatusUnauthorized,
				map[string]any{
					"error":  locale.LCtx(c, MsgErr, nil),
					"status": e.Err401TokenError,
				})
			return
		}
		c.Set("user", claims)
		c.Next(ctx)
	}
}

func fromAuthHeader(c *app.RequestContext) (string, bool) {
	authHeader := c.Request.Header.Get("Authorization")
	if authHeader == "" {
		return "", true
	}
	authToken := strings.TrimSpace(strings.Split(authHeader, BearerSchema)[1])
	if authToken == "" {
		return "", true
	}
	return authToken, false
}

func fromQueryParam(c *app.RequestContext) (string, bool) {
	authToken := c.Query("token")
	if authToken == "" {
		return "", true
	}
	return authToken, false
}
