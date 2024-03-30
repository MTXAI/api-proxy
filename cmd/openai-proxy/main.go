package main

import (
	"flag"
	"fmt"
	"log/slog"
	"os"

	"github.com/mtxai/api-proxy/pkg/config"
	"github.com/mtxai/api-proxy/pkg/proxy"
	"github.com/mtxai/api-proxy/pkg/server"
	"github.com/mtxai/api-proxy/pkg/utils"
)

var (
	// common config
	logDirFlag = flag.String("logs-dir", "./logs", "output logs dir, --logs-dir=./logs")
	debugFlag  = flag.Bool("debug", false, "debug logs level, --debug (default false)")
	stdoutFlag = flag.Bool("stdout", false, "output logs stdout, --stdout (default false)")
	addrFlag   = flag.String("addr", "0.0.0.0", "listen addr, --addr=0.0.0.0")
	portFlag   = flag.Int("port", 6789, "listen port, --port=6789")

	// proxy config
	remoteAddrFlag = flag.String("remote-addr", "api.openai.com", "openai api addr, --openai-addr=api.openai.com")
	clientCertFile = flag.String("client-cert-file", "", "client cert file, --client-cert-file=./client/ca-cert.pem")
	clientKeyFile  = flag.String("client-key-file", "", "client key file, --client-key-file=./client/ca-key.pem")

	// http server config
	serverCertFile = flag.String("server-cert-file", "", "server cert file, --server-cert-file=./server/ca-cert.pem")
	serverKeyFile  = flag.String("server-key-file", "", "server key file, --server-key-file=./server/ca-key.pem")
)

func main() {
	flag.Parse()
	if !flag.Parsed() {
		flag.Usage()
		os.Exit(0)
	}

	cfg := &config.Config{
		Addr: *addrFlag,
		Port: *portFlag,
	}

	fmt.Printf("init proxy...\n")
	initLogOrDie()

	verifyConfigOrDie(cfg)

	stopCh := utils.SetupSignalHandler()
	startProxyServerOrDie(cfg, stopCh)
	<-stopCh
	slog.Info("exit!")
}

func initLogOrDie() {
	err := utils.InitLog(*logDirFlag, *debugFlag, *stdoutFlag)
	if err != nil {
		panic(err)
	}

	err = utils.InitGinLog(*logDirFlag, *debugFlag, *stdoutFlag)
	if err != nil {
		panic(err)
	}

	slog.Info("logger init", "path", *logDirFlag)
}

func verifyConfigOrDie(cfg *config.Config) {
	if cfg.Addr == "" {
		cfg.Addr = "0.0.0.0"
		slog.Warn("proxy addr not set, will listen 0.0.0.0")
	}
	if cfg.Port == 0 {
		panic("proxy port must set")
	}

	slog.Info("API Proxy config verified", "cfg", cfg)
}

func startProxyServerOrDie(cfg *config.Config, stopCh chan struct{}) {
	p, err := proxy.OpenAI(*remoteAddrFlag, proxy.WithTLS(*clientCertFile, *clientKeyFile))
	if err != nil {
		panic(fmt.Sprintf("create openai proxy failed: %s", err.Error()))
	}

	srv, err := server.NewServer(p)
	go func() {
		defer func() {
			utils.PanicRecover("ProxyServer", recover(), true)
			// stop and exit
			close(stopCh)
		}()

		slog.Info("start openai proxy server")
		if *serverCertFile != "" && *serverKeyFile != "" {
			err = srv.RunTLS(fmt.Sprintf("%s:%d", cfg.Addr, cfg.Port), *serverCertFile, *serverKeyFile)
		} else {
			err = srv.Run(fmt.Sprintf("%s:%d", cfg.Addr, cfg.Port))
		}
		if err != nil {
			panic(fmt.Errorf("start openai proxy server failed, err: %v", err))
		}
	}()
}
