package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func Hello(c *gin.Context) {
	// 构造包含粗体黑字的HTML响应
	htmlResponse := "<html><body><b>后端服务已启动。</b></body></html>"
	// 返回HTML响应
	c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(htmlResponse))
}
