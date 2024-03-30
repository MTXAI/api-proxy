package server

import (
	"bytes"

	"github.com/gin-gonic/gin"
	"github.com/mtxai/api-proxy/pkg/proxy"
)

func APIProxyStatisticsMiddleware(p proxy.Proxy) gin.HandlerFunc {
	return func(c *gin.Context) {
		if !p.IsAPISupported(c.Request.URL.Path) {
			return
		}
		w := &CustomResponseWriter{
			body:           bytes.NewBufferString(""),
			ResponseWriter: c.Writer,
		}
		c.Writer = w
		c.Next()

		p.PerformStatistics(c.Request.URL.Path, w.body)
	}
}

// more middlewares
