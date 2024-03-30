package proxy

import (
	"log/slog"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"

	"github.com/mtxai/api-proxy/pkg/utils"
)

type (
	// OpenAI support Server-Sent Events
	openAI struct {
		cfg          *proxyConfig
		reverseProxy *httputil.ReverseProxy
		remoteServer *url.URL
		transport    http.RoundTripper
	}
)

func newOpenAIProxy(cfg *proxyConfig) (Proxy, error) {
	remoteServer := &url.URL{
		Scheme: "https",
		Host:   cfg.remoteAddr,
	}
	reverseProxy := httputil.NewSingleHostReverseProxy(remoteServer)

	var err error
	transport := http.DefaultTransport
	if cfg.certFile != "" && cfg.keyFile != "" {
		transport, err = utils.HTTPSTransport(cfg.certFile, cfg.keyFile)
		if err != nil {
			return nil, err
		}
	}

	proxy := &openAI{
		cfg:          cfg,
		reverseProxy: reverseProxy,
		remoteServer: remoteServer,
		transport:    transport,
	}

	reverseProxy.Rewrite = proxy.modifyRequest(reverseProxy.Director)
	reverseProxy.Director = nil
	reverseProxy.ModifyResponse = proxy.modifyResponse()
	reverseProxy.ErrorHandler = proxy.errorHandler()
	reverseProxy.Transport = proxy.transport
	// ignore FlushInterval
	reverseProxy.FlushInterval = -1

	return proxy, nil
}

func (proxy *openAI) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	proxy.reverseProxy.ServeHTTP(rw, req)
}

func (proxy *openAI) RoundTrip(req *http.Request) (*http.Response, error) {
	return proxy.transport.RoundTrip(req)
}

// modifyRequest modify request send to remote addr
func (proxy *openAI) modifyRequest(director func(req *http.Request)) func(*httputil.ProxyRequest) {
	defer utils.PanicRecover("modifyRequest", recover(), true)

	return func(proxyReq *httputil.ProxyRequest) {
		// rewrite url, replace addr with remote addr
		inReq := proxyReq.In
		outReq := proxyReq.Out
		director(outReq)

		if !proxy.CheckIsAPISupported(outReq.URL.Path) {
			return
		}

		// modify request header
		outReq.Header.Set("X-Forwarded-For", inReq.Header.Get("X-Forwarded-For"))
		outReq.Header.Set("X-Real-IP", inReq.Header.Get("X-Real-IP"))
		outReq.Header.Set("Host", proxy.remoteServer.Host)
		outReq.Host = proxy.remoteServer.Host
		outReq.RemoteAddr = inReq.RemoteAddr
	}
}

// modifyResponse modify response return to client
func (proxy *openAI) modifyResponse() func(resp *http.Response) (err error) {
	defer utils.PanicRecover("modifyResponse", recover(), true)

	return func(resp *http.Response) error {
		if resp == nil {
			slog.Debug("receive a nil response, ignore")
			return nil
		}

		if resp.Body == nil {
			slog.Debug("receive a nil response body, ignore")
			return nil
		}

		if resp.StatusCode != http.StatusOK {
			slog.Debug("receive a not ok response code, ignore", "code", resp.StatusCode)
			return nil
		}

		req := resp.Request
		if req == nil {
			slog.Debug("receive a nil request, ignore")
			return nil
		}

		if !proxy.CheckIsAPISupported(req.URL.Path) {
			return nil
		}

		// modify response header
		resp.Header.Add("Access-Control-Allow-Origin", "*")
		return nil
	}
}

func (proxy *openAI) errorHandler() func(rw http.ResponseWriter, req *http.Request, err error) {
	return func(rw http.ResponseWriter, req *http.Request, err error) {
		slog.Error("remote proxy got error", "url", utils.ReqString(req), "err", err)
	}
}

func (proxy *openAI) CheckIsAPISupported(path string) bool {
	if strings.HasPrefix(path, "/v1") {
		return true
	}
	return false
}
