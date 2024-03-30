package server

import (
	"bytes"

	"github.com/gin-gonic/gin"
	"github.com/mtxai/api-proxy/pkg/proxy"
)

func PerformStatisticsMiddleware(p proxy.Proxy) gin.HandlerFunc {
	return func(c *gin.Context) {
		w := &CustomResponseWriter{
			body:           bytes.NewBufferString(""),
			ResponseWriter: c.Writer,
		}
		c.Writer = w
		c.Next()

		if p.IsAPISupported(c.Request.URL.Path) {
			p.PerformStatistics(c.Request.URL.Path, w.body)
		}
		w.body.Reset()
	}
}

// more middlewares
