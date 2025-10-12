package logger

import (
	"fmt"
	"log/slog"

	"github.com/pkg/errors"
)

type StackTracer interface {
	StackTrace() errors.StackTrace
	Error() string
}

func Error(err error) []slog.Attr {
	var er StackTracer
	if errors.As(err, &er) {
		stack := ""
		for _, f := range er.StackTrace() {
			stack += fmt.Sprintf("%+s:%d\n", f, f)
		}
		return []slog.Attr{
			{Key: "stack", Value: slog.StringValue(stack)},
			{Key: "message", Value: slog.StringValue(err.Error())},
		}
	}
	return []slog.Attr{
		{Key: "message", Value: slog.StringValue(err.Error())},
	}
}
