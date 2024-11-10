package middleware

import (
	"github.com/go-chi/chi/v5/middleware"
	"log/slog"
	"net/http"
	"time"
)

func Logger(log *slog.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {

		log = log.With(slog.String("component", "middleware/logger"))
		log.Info("http logger enabled")

		fn := func(w http.ResponseWriter, r *http.Request) {
			entry := log.With(
				slog.String("method", r.Method),
				slog.String("url", r.URL.String()),
				slog.String("path", r.URL.Path),
				slog.String("remoteAddr", r.RemoteAddr),
				slog.String("userAgent", r.UserAgent()),
				slog.String("request_id", middleware.GetReqID(r.Context())),
			)
			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
			t1 := time.Now()
			defer func() {
				entry.Info("http request complete",
					slog.Int("status", ww.Status()),
					slog.Duration("duration", time.Since(t1)),
					slog.Int("bytes", ww.BytesWritten()),
					slog.String("content-type", ww.Header().Get("Content-Type")),
				)
			}()
			next.ServeHTTP(ww, r)
		}

		return http.HandlerFunc(fn)
	}

}
