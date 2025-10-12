package routes

import (
	"context"
	"strings"
	"time"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/middlewares/server/recovery"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/dbunt1tled/url-shortener/assets"
	"github.com/dbunt1tled/url-shortener/internal/app/shorturl"
	"github.com/dbunt1tled/url-shortener/internal/config"
	"github.com/dbunt1tled/url-shortener/internal/lib/hasher"
	"github.com/dbunt1tled/url-shortener/internal/lib/http-server/middleware"
	"github.com/hertz-contrib/cors"
)

type Router struct {
	srv        *server.Hertz
	cfg        *config.Config
	logger     *config.AppLogger
	hasher     *hasher.Hasher
	urlHandler *shorturl.URLHandler
}

func NewRouter(
	srv *server.Hertz,
	cfg *config.Config,
	logger *config.AppLogger,
	hasher *hasher.Hasher,
	urlHandler *shorturl.URLHandler,
) *Router {
	return &Router{
		srv:        srv,
		cfg:        cfg,
		logger:     logger,
		hasher:     hasher,
		urlHandler: urlHandler,
	}
}

func (r *Router) Register() {
	r.srv.GET("/favicon.ico", func(c context.Context, ctx *app.RequestContext) {
		ctx.Header("Content-Type", "image/x-icon")
		data, _ := assets.Files.ReadFile("favicon.ico")
		_, _ = ctx.Write(data)
	})
	r.srv.Use(
		recovery.Recovery(),
		middleware.LoggerMiddleware(r.logger, r.cfg.LogLevelStatus),
		middleware.ErrorHandler(r.logger),
	)
	url := r.srv.Group("/url",
		cors.New(cors.Config{
			AllowOrigins:     strings.Split(r.cfg.AccessControlAllowOrigin, ","),
			AllowMethods:     strings.Split(r.cfg.AccessControlAllowMethods, ","),
			AllowHeaders:     strings.Split(r.cfg.AccessControlAllowHeaders, ","),
			ExposeHeaders:    strings.Split(r.cfg.AccessControlExposeHeaders, ","),
			AllowCredentials: true,
			MaxAge:           30 * 24 * time.Hour}),
		middleware.AuthBearerMiddleware(r.hasher),
	)
	url.POST("", r.urlHandler.NewUrl)
	r.srv.GET("/:alias", r.urlHandler.Redirect)
}
