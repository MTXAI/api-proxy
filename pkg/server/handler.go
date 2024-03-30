package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mtxai/api-proxy/pkg/proxy"
)

func APIProxyHandler(p proxy.Proxy) func(c *gin.Context) {
	return func(c *gin.Context) {
		if p.CheckIsAPISupported(c.Request.URL.Path) {
			p.ServeHTTP(c.Writer, c.Request)
		} else {
			c.JSON(http.StatusNotFound, gin.H{"message": "404 Not Found"})
		}
	}
}
