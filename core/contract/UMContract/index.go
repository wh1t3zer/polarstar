package UMContract

// U本位合约
import (
	"context"
	"fmt"
	"github.com/adshao/go-binance/v2/futures"
	"polarstar/util"
)

var (
	apiKey        string
	secretKey     string
	apiKeyTest    string
	secretKeyTest string
)

func init() {
	apiKey = util.Config.Binance.MainNet.ApiKey
	secretKey = util.Config.Binance.MainNet.SecretKey
	apiKeyTest = util.Config.Binance.TestNet.ApiKeyTest
	secretKeyTest = util.Config.Binance.TestNet.SecretKeyTest
}

// 合约K线
func GetUMKline(symbol string, Itl string) {
	futureDoneC := make(chan struct{})
	futureStopC := make(chan struct{})
	fsHandler := futures.WsKlineHandler(func(event *futures.WsKlineEvent) {
		fmt.Println(event)
	})
	fsErrHandler := futures.ErrHandler(func(err error) {
		if err != nil {
			fmt.Println("fsErr : ", err)
			return
		}
	})
	futureDoneC, futureStopC, err := futures.WsKlineServe(symbol, Itl, fsHandler, fsErrHandler)
	if err != nil {
		fmt.Println("fsKline err : ", err)
		return
	}
	go func() {
		for {
			select {
			case <-futureStopC:
				fmt.Println("stop futures stream")
				close(futureStopC)
			}
		}
	}()
	<-futureDoneC
}

// 调整杠杆
func ChangeLever(symbol string, lever int) {
	fmt.Println(apiKey, secretKey)
	ctx := context.Background()
	client := futures.NewClient(apiKey, secretKey)
	resp, err := client.NewChangeLeverageService().Symbol(symbol).Leverage(lever).Do(ctx)
	if err != nil {
		fmt.Println("change lever fail: ", err)
		return
	}
	fmt.Println(&resp)

}

// 下单买入
func ContractOrderBuy() {

}

// 下单买入测试
func ContractOrderBuyTest() {

}

// 下单卖出
func ContractOrderSell() {

}

// 下单卖出测试
func ContractOrderSellTest() {

}

// 撤单
func CancelOrder() {

}

// 撤单测试
func CancelOrderTest() {

}

// 获取所有在挂单的信息
func GetCAOrderList() {

}

// 获取挂单信息
func GetSymbolOrder() {
}
