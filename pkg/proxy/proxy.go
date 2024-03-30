package proxy

import (
	"fmt"
	"net/http"
)

type (
	Proxy interface {
		http.Handler
		http.RoundTripper
		CheckIsAPISupported(path string) bool
	}
	proxyConfig struct {
		remoteAddr string

		certFile string
		keyFile  string
	}
	Option func(cfg *proxyConfig)
)

func OpenAI(remoteAddr string, opts ...Option) (Proxy, error) {
	if remoteAddr == "" {
		return nil, fmt.Errorf("remote address must set")
	}
	cfg := &proxyConfig{
		remoteAddr: remoteAddr + ":443",
	}

	for _, o := range opts {
		o(cfg)
	}
	return newOpenAIProxy(cfg)
}

func WithTLS(certFile, keyFile string) Option {
	return func(cfg *proxyConfig) {
		cfg.certFile = certFile
		cfg.keyFile = keyFile
	}
}
