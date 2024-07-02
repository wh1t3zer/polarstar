package middleware

import (
	"github.com/gin-gonic/gin"
	"polarstar/util"
)

func IPAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		isMatched := false
		for _, host := range util.CC.Http.AllowIP {
			if c.ClientIP() == host {
				isMatched = true
			}
		}
		if !isMatched {
			c.Abort()
			return
		}
		c.Next()
	}
}
