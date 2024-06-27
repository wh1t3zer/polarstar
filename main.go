package main

import (
	"os"
	"os/signal"
	"polarstar/core/contract/UMContract"
	"polarstar/util"
	"syscall"
)

func main() {
	util.Init()
	UMContract.GetUMKline("BNBUSDT", "1m")
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGKILL, syscall.SIGQUIT, syscall.SIGINT, syscall.SIGTERM)
	<-quit
}
