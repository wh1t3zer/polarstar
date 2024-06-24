package main

import (
	"os"
	"os/signal"
	"polarstar/core/contract/UMContract"
	"polarstar/util"
	"syscall"
)

func main() {
	util.Banner()
	//UMContract.GetUMKline("btcusdt", "1m")
	util.InitConfig("./conf/", "base")
	UMContract.ChangeLever("btcusdt", 100)
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGKILL, syscall.SIGQUIT, syscall.SIGINT, syscall.SIGTERM)
	<-quit
}
