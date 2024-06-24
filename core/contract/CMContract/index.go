package CMContract

import (
	"fmt"
	"github.com/adshao/go-binance/v2/delivery"
	"github.com/e421083458/golang_common/lib"
)

// 币本位合约
var (
	apiKey        = lib.GetStringConf("base.Binance.Main.apiKey")
	secretKey     = lib.GetStringConf("base.Binance.Main.secretKey")
	apiKeyTest    = lib.GetStringConf("base.Binance.Test.apiKeyTest")
	secretKeyTest = lib.GetStringConf("base.Binance.Test.secretKeyTest")
)

// 合约K线
func GetCMKline(symbol string, Itl string) {
	deHandler := delivery.WsKlineHandler(func(event *delivery.WsKlineEvent) {
		fmt.Println(event)
	})
	deErrHandler := delivery.ErrHandler(func(err error) {
		if err != nil {
			fmt.Println("dsErr : ", err)
			return
		}
	})
	deliveryDoneC, deliveryStopC, err := delivery.WsKlineServe(symbol, Itl, deHandler, deErrHandler)
	if err != nil {
		fmt.Println("dsKline err : ", err)
		return
	}
	go func() {
		for {
			select {
			case <-deliveryStopC:
				fmt.Println("stop delivery stream")
				close(deliveryStopC)
				return
			}
		}
	}()
	<-deliveryDoneC
}

// 调整杠杆
func ChangeLever() {

}

// 下单买入
func ContractOrderBuy() {

}

// 下单卖出
func ContractOrderSell() {

}

// 撤单
func CancelOrder() {

}

// 获取所有在挂单的信息
func GetCAOrderList() {

}

// 获取挂单信息
func GetSymbolOrder() {
}
