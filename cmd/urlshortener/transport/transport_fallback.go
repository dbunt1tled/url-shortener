//go:build windows
// +build windows

package transport

import (
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/network/standard"
)

// NewServer fallback — standard (TCP only).
func NewServer(opts ...config.Option) *server.Hertz {
	opts = append(opts, server.WithTransport(standard.NewTransporter()))
	return server.New(opts...)
}
