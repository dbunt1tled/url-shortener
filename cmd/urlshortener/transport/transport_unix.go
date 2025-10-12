//go:build linux || darwin
// +build linux darwin

package transport

import (
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/common/config"
	"github.com/cloudwego/hertz/pkg/network/netpoll"
)

// NewServer создает Hertz с netpoll (Linux/macOS).
func NewServer(opts ...config.Option) *server.Hertz {
	opts = append(opts, server.WithTransport(netpoll.NewTransporter))
	return server.New(opts...)
}
