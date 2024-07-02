package CMContract

import (
	"context"
	"fmt"
	"github.com/adshao/go-binance/v2/delivery"
	"polarstar/util"
)

var (
	logger = util.NewCustomLogger()
)

// 币本位合约

// 合约K线
func GetCMKline(symbol string, Itl string) {
	deliveryDoneC := make(chan struct{})
	deliveryStopC := make(chan struct{})
	deHandler := delivery.WsContinuousKlineHandler(func(event *delivery.WsContinuousKlineEvent) {
		fmt.Println(event.Kline)
	})
	deErrHandler := delivery.ErrHandler(func(err error) {
		if err != nil {
			logger.Error("币本位合约Err: %v", err)
			return
		}
	})
	deliveryDoneC, deliveryStopC, err := delivery.WsContinuousKlineServe(symbol, "perpetual", Itl, deHandler, deErrHandler)
	if err != nil {
		logger.Error("获取币本位合约K线失败: %v", err)
		return
	}
	go func() {
		for {
			select {
			case <-deliveryStopC:
				logger.Info("关闭币本位通信流进程")
				close(deliveryStopC)
				return
			}
		}
	}()
	<-deliveryDoneC
}

// 调整杠杆
func ChangeCMLever(symbol string, lever int) {
	ctx := context.Background()
	client := delivery.NewClient(util.ApiKey, util.SecretKey)
	resp, err := client.NewChangeLeverageService().
		Leverage(lever).
		Symbol(symbol).
		Do(ctx)
	if err != nil {
		logger.Error("调整CM杠杆失败:\n %v", err)
	}
	logger.Info("更改CM杠杆成功\n当前交易对: %s，杠杆数: %d，最大可开仓资金数: %s", resp.Symbol, resp.Leverage, resp.MaxQuantity)
}

// 下单买入
func ContractOrderBuy() {

}

// 平仓
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
