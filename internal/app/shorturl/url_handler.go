package shorturl

import (
	"context"
	"log/slog"
	"time"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"github.com/dbunt1tled/url-shortener/internal/config"
	"github.com/dbunt1tled/url-shortener/internal/domain/enum"
	"github.com/dbunt1tled/url-shortener/internal/domain/repository"
	"github.com/dbunt1tled/url-shortener/internal/lib/http-server/e"
)

type URLHandler struct {
	urlService *URLService
	logger     *config.AppLogger
}

func NewURLHandler(service *URLService, logger *config.AppLogger) *URLHandler {
	return &URLHandler{urlService: service, logger: logger}
}

func (h *URLHandler) NewUrl(c context.Context, ctx *app.RequestContext) {
	var (
		url *URLBase
		req URLCreate
		err error
		uid int64
	)

	if err = ctx.BindAndValidate(&req); err != nil {
		ctx.Error(e.NewValidationError(err.Error(), e.Err422URLValidateCreateError))
		return
	}
	user, ok := ctx.Get("user")
	if ok {
		uid = int64(user.(map[string]interface{})["iss"].(float64))
		req.UserID = &uid
	}

	url, err = h.urlService.ShortenURL(c, &req)
	if err != nil {
		ctx.Error(e.NewUnprocessableEntityError(err.Error(), e.Err422URLSaveError))
		return
	}

	ctx.JSON(consts.StatusOK, map[string]interface{}{
		"data": url,
	})
}

func (h *URLHandler) Redirect(c context.Context, ctx *app.RequestContext) {
	var (
		url *URL
		req URLRedirect
		err error
	)

	if err = ctx.BindAndValidate(&req); err != nil {
		ctx.Error(e.NewValidationError(err.Error(), e.Err422URLValidateCreateError))
		return
	}

	url, err = h.urlService.One(c, []repository.Filter{
		{Field: "code", Operator: "=", Value: req.Alias},
		{Field: "status", Operator: "=", Value: enum.Active},
	}, nil)
	if err != nil {
		ctx.Error(e.NewUnprocessableEntityError(err.Error(), e.Err422URLFindError))
		return
	}
	if url == nil {
		h.logger.Warn("url not found", req)
		ctx.Error(e.NewNotFoundError("url not found", e.Err404URLNotFound))
		return
	}

	if url.ExpiredAt.Before(time.Now()) {
		h.logger.Warn("url not found", err, slog.Any("request", req))
		ctx.Error(e.NewNotFoundError("url not found", e.Err404URLExpired))
		_, err = h.urlService.Delete(c, url.ID)
		if err != nil {
			h.logger.Error("url delete error", err, slog.Any("request", req))
		}
		return
	}

	url, err = h.urlService.Update(c, url.ID, map[string]any{
		"count":           url.Count + 1,
		"last_visited_at": time.Now(),
	})
	if err != nil {
		h.logger.Error("url update error", err, slog.Any("request", req))
	}

	ctx.Redirect(consts.StatusMovedPermanently, []byte(url.URL))
}
