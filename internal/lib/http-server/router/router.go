package router

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"go_first/internal/config/env"
	auth "go_first/internal/lib/http-server/action/auth"
	"go_first/internal/lib/http-server/action/urlshort"
	"go_first/internal/lib/http-server/action/user"
	httpMiddlewares "go_first/internal/lib/http-server/middleware"
	"go_first/storage"
	"log/slog"
	"strings"
	"time"
)

func NewRouter(
	storage storage.Storage,
	cfg *env.Config,
	locale *i18n.Localizer,
	bundle *i18n.Bundle,
	log *slog.Logger,
) *chi.Mux {
	router := chi.NewRouter()
	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Compress(5))
	router.Use(middleware.Timeout(60 * time.Second))
	router.Use(httpMiddlewares.Logger(log))
	router.Use(httpMiddlewares.Language(locale, bundle))
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)
	router.Use(middleware.StripSlashes)
	router.Use(httpMiddlewares.Base())
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins: strings.Split(cfg.CORS.AccessControlAllowOrigin, ","),
		// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowedMethods:   strings.Split(cfg.CORS.AccessControlAllowMethods, ","),
		AllowedHeaders:   strings.Split(cfg.CORS.AccessControlAllowHeaders, ","),
		ExposedHeaders:   strings.Split(cfg.CORS.AccessControlExposeHeaders, ","),
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))
	router.Route("/url", func(r chi.Router) {
		r.Use(httpMiddlewares.AuthBearer(storage))
		r.Post("/", urlshort.CreateUrlAction(log, storage))
	})
	router.Route("/users", func(r chi.Router) {
		r.Post("/", user.CreateUserAction(log, storage))
	})
	router.Route("/auth", func(r chi.Router) {
		r.Post("/login", auth.LoginAction(log, storage, locale))
	})
	router.Get("/{alias}", urlshort.GetUrlAction(storage))

	return router
}
