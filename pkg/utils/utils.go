package utils

import (
	"crypto/tls"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"
)

func ReqString(req *http.Request) string {
	if req == nil {
		return ""
	}
	return fmt.Sprintf("%s of %s", req.URL.String(), req.Form)
}

func PanicRecover(fcn string, r interface{}, printAllStack bool) {
	if r == nil {
		r = recover()
	}
	if r != nil {
		err, ok := r.(error)
		if !ok {
			err = fmt.Errorf("%v", r)
		}
		slog.Error("panic recover", "function", fcn, "err", err)
		if printAllStack {
			const size = 64 << 10
			buf := make([]byte, size)
			buf = buf[:runtime.Stack(buf, true)]
			slog.Error("print final stack", "stack", string(buf))
		}
	}
}

var shutdownSignals = []os.Signal{syscall.SIGINT, syscall.SIGTERM}
var onlyOneSignalHandler = make(chan struct{})

func SetupSignalHandler() (stopCh chan struct{}) {
	close(onlyOneSignalHandler)

	stop := make(chan struct{})
	c := make(chan os.Signal, 2)
	signal.Notify(c, shutdownSignals...)
	go func() {
		sig := <-c
		slog.Info("Received first signal", "signal", sig)
		close(stop)
		sig = <-c
		slog.Info("Received second signal", "signal", sig)
		os.Exit(1)
	}()

	return stop
}

func HTTPSTransport(certFile, keyFile string) (*http.Transport, error) {
	if certFile == "" || keyFile == "" {
		return nil, errors.New("x509 key pair not found")
	}
	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		return nil, err
	}

	tlsConf := &tls.Config{
		Certificates:       []tls.Certificate{cert},
		InsecureSkipVerify: true,
	}
	dialer := &net.Dialer{
		Timeout:   30 * time.Second,
		KeepAlive: 30 * time.Second,
	}
	return &http.Transport{
		TLSClientConfig:       tlsConf,
		Proxy:                 http.ProxyFromEnvironment,
		DialContext:           dialer.DialContext,
		ForceAttemptHTTP2:     true,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}, nil
}
