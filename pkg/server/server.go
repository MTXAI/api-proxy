package server

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mtxai/api-proxy/pkg/proxy"
)

type Server struct {
	*gin.Engine
	p proxy.Proxy
}

func NewServer(p proxy.Proxy) (*Server, error) {
	r := gin.New()
	srv := &Server{
		Engine: r,
		p:      p,
	}
	srv.Register()
	srv.RegisterMiddleware()
	return srv, nil
}

func (srv *Server) Register() {
	srv.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})
	srv.NoRoute(APIProxyHandler(srv.p))
}

func (srv *Server) RegisterMiddleware() {
	srv.Use(gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		return fmt.Sprintf("time=%s proto=%s method=%s path=%s code=%d clientIP=%s latenc=%s ua=\"%s\" msg=%s\n",
			param.TimeStamp.Format(time.RFC3339),
			param.Request.Proto,
			param.Method,
			param.Path,
			param.StatusCode,
			param.ClientIP,
			param.Latency,
			param.Request.UserAgent(),
			param.ErrorMessage,
		)
	}))
	srv.Use(gin.Recovery())
	srv.Use(APIProxyStatisticsMiddleware(srv.p))
}
