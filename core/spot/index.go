package spot

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/adshao/go-binance/v2"
	"github.com/e421083458/golang_common/lib"
	"log"
	"polarstar/util"
)

var (
	apiKey        = lib.GetStringConf("base.Binance.Main.apiKey")
	secretKey     = lib.GetStringConf("base.Binance.Main.secretKey")
	apiKeyTest    = lib.GetStringConf("base.Binance.Test.apiKeyTest")
	secretKeyTest = lib.GetStringConf("base.Binance.Test.secretKeyTest")
)

// 现货买入
// params: orderT: limit(限价单)、market(市价单)、stop(止损单)、stop-limit(限价止损单)
func SpotOrderBuy(symbol string, quantity string, price string, orderT string) {
	var orderType binance.OrderType
	client := binance.NewClient(apiKey, secretKey)
	ctx := context.Background()
	side := binance.SideTypeBuy
	switch orderT {
	case "limit":
		orderType = binance.OrderTypeLimit
	case "market":
		orderType = binance.OrderTypeMarket
	case "stop":
		orderType = binance.OrderTypeStopLoss
	case "stop-limit":
		orderType = binance.OrderTypeStopLossLimit
	}
	order, err := client.NewCreateOrderService().
		Symbol(symbol).
		Side(side).
		Type(orderType).
		TimeInForce(binance.TimeInForceTypeGTC).
		Quantity(quantity).
		Price(price).
		Do(ctx)
	if err != nil {
		log.Fatalf("买入失败: %v", err)
	}
	// 输出订单详情
	fmt.Printf("订单详情: %+v\n", order)
}

// 现货卖出
// params: orderT: limit(限价单)、market(市价单)、stop(止损单)、stop-limit(限价止损单)
func SpotOrderSell(symbol string, quantity string, price string, orderT string) {
	var orderType binance.OrderType
	client := binance.NewClient(apiKey, secretKey)
	ctx := context.Background()
	side := binance.SideTypeBuy
	switch orderT {
	case "limit":
		orderType = binance.OrderTypeLimit
	case "market":
		orderType = binance.OrderTypeMarket
	case "stop":
		orderType = binance.OrderTypeStopLoss
	case "stop-limit":
		orderType = binance.OrderTypeStopLossLimit
	}
	order, err := client.NewCreateOrderService().
		Symbol(symbol).
		Side(side).
		Type(orderType).
		QuoteOrderQty(quantity).
		Price(price).
		Do(ctx)
	if err != nil {
		log.Fatalf("卖出失败: %v", err)
	}
	// 输出订单详情
	fmt.Printf("订单详情: %+v\n", order)
}

// 现货买入测试
// params: orderT: limit(限价单)、market(市价单)、stop(止损单)、stop-limit(限价止损单)
func SpotOrderBuyTest(symbol string, quantity string, price string, orderT string) {
	// 测试
	//设置为测试网络（沙盒环境）模拟交易
	var orderType binance.OrderType
	client := binance.NewClient(apiKeyTest, secretKeyTest)
	client.BaseURL = binance.BaseAPITestnetURL
	ctx := context.Background()
	side := binance.SideTypeBuy
	switch orderT {
	case "limit":
		orderType = binance.OrderTypeLimit
	case "market":
		orderType = binance.OrderTypeMarket
	case "stop":
		orderType = binance.OrderTypeStopLoss
	case "stop-limit":
		orderType = binance.OrderTypeStopLossLimit
	}
	// 创建订单
	resp, err := client.NewCreateOrderService().
		Symbol(symbol).
		Side(side).
		Type(orderType).
		Quantity(quantity).
		TimeInForce("GTC").
		Price(price).
		Do(ctx)

	if err != nil {
		log.Fatalf("创建订单失败: %v", err)
	} else {
		fmt.Println("订单创建成功", resp.OrderID)
	}
}

// 现货卖出测试
// params: orderT: limit(限价单)、market(市价单)、stop(止损单)、stop-limit(限价止损单)
func SpotOrderSellTest(symbol string, quantity string, price string, orderT string) {
	var orderType binance.OrderType
	client := binance.NewClient(apiKeyTest, secretKeyTest)
	ctx := context.Background()
	side := binance.SideTypeSell
	switch orderT {
	case "limit":
		orderType = binance.OrderTypeLimit
	case "market":
		orderType = binance.OrderTypeMarket
	case "stop":
		orderType = binance.OrderTypeStopLoss
	case "stop-limit":
		orderType = binance.OrderTypeStopLossLimit
	}
	client.BaseURL = binance.BaseAPITestnetURL
	order, err := client.NewCreateOrderService().
		Symbol(symbol).
		Side(side).Type(orderType).
		Quantity(quantity).
		TimeInForce("GTC").
		Price(price).
		Do(ctx)
	if err != nil {
		log.Fatalf("卖出失败: %v", err)
	}
	// 输出订单详情
	fmt.Printf("订单详情: %+v\n", order)
}

// 现货撤单
func SpotCancelOrder(orderId int64) {
	ctx := context.Background()
	client := binance.NewClient(apiKey, secretKey)
	cOrder, err := client.NewCancelOrderService().
		OrderID(orderId).
		Do(ctx)
	if err != nil {
		log.Fatal("撤单失败: ", err)
	}
	fmt.Println(&cOrder)
}

// 现货撤单测试
func SpotCancelOrderTest(orderId int64) {
	ctx := context.Background()
	client := binance.NewClient(apiKeyTest, secretKeyTest)
	client.BaseURL = binance.BaseAPITestnetURL
	cOrder, err := client.NewCancelOrderService().
		OrderID(orderId).
		Do(ctx)
	if err != nil {
		log.Fatal("撤单失败:", err)
	}
	fmt.Println(&cOrder)
}

// 获取交易历史
//
//	func GetTradeHistory() {
//		client := binance.NewClient(apiKey, secretKey)
//	}
//
// 获取历史测试
func GetSpotHistoryTest() {
	client := binance.NewClient(apiKeyTest, secretKeyTest)
	client.BaseURL = binance.BaseAPITestnetURL
	ctx := context.Background()
	history, err := client.NewListOrdersService().Do(ctx)
	if err != nil {
		fmt.Println("获取交易历史失败:", err)
		return
	}
	fmt.Println(history)
}

// 获得K线
func GetSpotKline(symbol string, Itl string) {
	spotDoneC := make(chan struct{})
	spotStopC := make(chan struct{})
	spotHandler := binance.WsKlineHandler(func(event *binance.WsKlineEvent) {
		fmt.Println(event)
	})
	spotErrHandler := binance.ErrHandler(func(err error) {
		if err != nil {
			fmt.Println("spotErr: ", err)
			return
		}
	})
	spotDoneC, spotStopC, err := binance.WsKlineServe(symbol, Itl, spotHandler, spotErrHandler)
	if err != nil {
		fmt.Println("spot Kline err : ", err)
		return
	}
	go func() {
		for {
			select {
			case <-spotStopC:
				fmt.Println("stop spot stream")
				close(spotStopC)
			}
		}
	}()
	<-spotDoneC

}

// 获得所有挂单列表
func GetOrderList(symbol string) {
	var info string
	ctx := context.Background()
	client := binance.NewClient(apiKey, secretKey)
	res, err := client.NewListOrdersService().
		Symbol(symbol).Do(ctx)
	if err != nil {
		fmt.Println(err)
		return
	}
	for _, list := range res {
		resp, err := json.Marshal(list)
		if err != nil {
			fmt.Println(err)
			return
		}
		info = string(resp)
	}
	util.PushWeChatBusiness(info)
}

// 获得订单列表测试
func GetOrderListTest(symbol string) {
	// 测试
	ctx := context.Background()
	client := binance.NewClient(apiKeyTest, secretKeyTest)
	client.BaseURL = binance.BaseAPITestnetURL
	res, err := client.NewListOrdersService().
		Symbol(symbol).
		Do(ctx)
	if err != nil {
		fmt.Println(err)
		return
	}
	for _, list := range res {
		resp, err := json.Marshal(list)
		if err != nil {
			return
		}
		fmt.Println(string(resp))
	}
}
