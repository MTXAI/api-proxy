package proxy

import (
	"bytes"
	"fmt"
	"net/http"
)

type (
	Proxy interface {
		http.Handler
		http.RoundTripper
		IsAPISupported(path string) bool
		PerformStatistics(path string, body *bytes.Buffer)
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
		remoteAddr: remoteAddr,
	}

	for _, o := range opts {
		o(cfg)
	}
	return newOpenAIProxy(cfg)
}

func WithTLS(certFile, keyFile string) Option {
	return func(cfg *proxyConfig) {
		cfg.remoteAddr = cfg.remoteAddr + ":443"
		cfg.certFile = certFile
		cfg.keyFile = keyFile
	}
}
