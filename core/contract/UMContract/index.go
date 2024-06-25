package UMContract

// U本位合约
import (
	"context"
	"fmt"
	"github.com/adshao/go-binance/v2"
	"github.com/adshao/go-binance/v2/futures"
	"polarstar/util"
)

var (
	logger = util.NewCustomLogger()
)

// 合约K线
func GetUMKline(symbol string, Itl string) {
	futureDoneC := make(chan struct{})
	futureStopC := make(chan struct{})
	fsHandler := futures.WsKlineHandler(func(event *futures.WsKlineEvent) {
		fmt.Println(event.Kline)
	})
	fsErrHandler := futures.ErrHandler(func(err error) {
		if err != nil {
			logger.Error("币本位合约Error: %v", err)
			return
		}
	})
	futureDoneC, futureStopC, err := futures.WsKlineServe(symbol, Itl, fsHandler, fsErrHandler)
	if err != nil {
		logger.Error("获取币本位合约K线失败: %v", err)
		return
	}
	go func() {
		for {
			select {
			case <-futureStopC:
				logger.Info("关闭币本位通信流进程")
				close(futureStopC)
			}
		}
	}()
	<-futureDoneC
}

// 调整杠杆
func ChangeLever(symbol string, lever int) {
	ctx := context.Background()
	client := futures.NewClient(util.ApiKey, util.SecretKey)
	resp, err := client.NewChangeLeverageService().Symbol(symbol).Leverage(lever).Do(ctx)
	if err != nil {
		logger.Error("更改杠杆失败: %v", err)
		return
	}
	if resp != nil {
		logger.Info("更改杠杆成功\n当前交易对: %s，杠杆数: %d，最大可用资金数: %s", resp.Symbol, resp.Leverage, resp.MaxNotionalValue)
	}
}

// 下单买入
func ContractOrderBuy() {
	futures.NewClient(util.ApiKey, util.SecretKey)

}

// 下单买入测试
func ContractOrderBuyTest(symbol string, quantity string, price string, orderT string, positionSide futures.PositionSideType) {
	var orderType futures.OrderType
	side := futures.SideTypeBuy
	ctx := context.Background()
	futures.UseTestnet = true
	client := futures.NewClient(util.ApiKeyTest, util.SecretKeyTest)

	res, err := client.NewCreateOrderService().
		Symbol(symbol).
		TimeInForce(futures.TimeInForceTypeGTC).
		Price(price).
		Type(orderType).
		PositionSide(positionSide).
		Quantity(quantity).
		Side(side).
		Do(ctx)
	if err != nil {
		fmt.Println("下单失败: ", err)
		return
	}
	fmt.Println(&res)
}

// 下单卖出
func ContractOrderSell() {
	//ctx := context.Background()
}

//// 下单卖出测试
//func ContractOrderSellTest() {
//	ctx := context.Background()
//	client := futures.NewClient(apiKeyTest, secretKeyTest)
//
//}

// 撤单
func CancelOrder(symbol string, orderId int64) {
	ctx := context.Background()
	client := futures.NewClient(util.ApiKey, util.SecretKey)
	resp, err := client.NewCancelOrderService().
		Symbol(symbol).
		OrderID(orderId).
		Do(ctx)
	if err != nil {
		fmt.Println("撤单失败: ", err)
		return
	}
	fmt.Println(resp)

}

// 撤单测试
func CancelOrderTest(symbol string, orderId int64) {
	ctx := context.Background()
	client := futures.NewClient(util.ApiKeyTest, util.SecretKeyTest)
	client.BaseURL = binance.BaseAPITestnetURL
	resp, err := client.NewCancelOrderService().
		Symbol(symbol).
		Do(ctx)
	if err != nil {
		fmt.Println("撤单失败: ", err)
		return
	}
	fmt.Println(resp)
}

// 获取所有在挂单的信息
func GetSymbolOrderList(symbol string) {
	ctx := context.Background()
	client := futures.NewClient(util.ApiKey, util.SecretKey)
	resp, err := client.NewListOpenOrdersService().
		Symbol(symbol).
		Do(ctx)
	if err != nil {
		fmt.Println("获取单个挂单信息失败: ", err)
		return
	}
	fmt.Println(resp)
}

// 获取挂单信息
func GetSymbolOrder(symbol string, orderId int64) {
	ctx := context.Background()
	client := futures.NewClient(util.ApiKey, util.SecretKey)
	resp, err := client.NewGetOpenOrderService().
		Symbol(symbol).
		OrderID(orderId).
		Do(ctx)
	if err != nil {
		fmt.Println("获取挂单信息失败: ", err)
		return
	}
	fmt.Println(resp)
}

// 获取挂单信息测试
func GetSymbolOrderTest(symbol string, orderId int64) {
	ctx := context.Background()
	futures.UseTestnet = true
	client := futures.NewClient(util.ApiKeyTest, util.SecretKeyTest)
	resp, err := client.NewGetOpenOrderService().
		Symbol(symbol).
		OrderID(orderId).
		Do(ctx)
	if err != nil {
		fmt.Println("获取挂单信息失败: ", err)
		return
	}
	fmt.Println(resp)
}

// 获得所有历史订单
func GetUMOrderList() {

}

// 获得历史订单列表
func GetUMOrder() {

}

// 获得历史订单测试
func GetUMOrderTest() {

}
