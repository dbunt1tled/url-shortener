package e

import (
	"context"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"github.com/pkg/errors"
)

func ErrorHandler(ctx context.Context, c *app.RequestContext, err error) {
	e := c.Errors.Last()
	if e == nil {
		return
	}
	var er *DomainError
	switch {
	case errors.As(e.Err, &er):
		c.JSON(er.Code, map[string]any{
			"error":  er.Error(),
			"status": er.Status,
		})
	default:
		c.JSON(consts.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
}
