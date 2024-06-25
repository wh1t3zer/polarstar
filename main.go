package main

import (
	"os"
	"os/signal"
	"polarstar/core/spot"
	"polarstar/util"
	"syscall"
)

func main() {
	util.Init()
	spot.GetOrderListTest("btcusdt")
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGKILL, syscall.SIGQUIT, syscall.SIGINT, syscall.SIGTERM)
	<-quit
}
