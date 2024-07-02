package main

import (
	"github.com/gin-gonic/gin"
	"os"
	"os/signal"
	"polarstar/route"
	"polarstar/util"
	"syscall"
)

func main() {
	util.Init()
	util.InitGinLogger()
	r := gin.Default()
	route.InitRoutes(r)
	r.Run(util.CC.Http.Addr)
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGKILL, syscall.SIGQUIT, syscall.SIGINT, syscall.SIGTERM)
	<-quit
}
