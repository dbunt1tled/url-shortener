//go:build windows
// +build windows

package transport

import (
	"github.com/cloudwego/hertz/pkg/network"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/network/standard"
)

// NewServer fallback — standard (TCP only).
func NewServer(opts ...config.Option) *server.Hertz {
	fallbackTransport := func(opt *config.Options) network.Transporter {
		return standard.NewTransporter(opt)
	}

	opts = append(opts, server.WithTransport(fallbackTransport))
	return server.New(opts...)
}
