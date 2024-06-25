package spot

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/adshao/go-binance/v2"
	"log"
	"polarstar/util"
)

var (
	logger = util.NewCustomLogger()
)

// 现货买入
// params: orderT: limit(限价单)、market(市价单)、stop(止损单)、stop-limit(限价止损单)
func SpotOrderBuy(symbol string, quantity string, price string, orderT string) {
	var orderType binance.OrderType
	client := binance.NewClient(util.ApiKey, util.SecretKey)
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
		logger.Error("买入失败: %v", err)
		return
	}
	// 输出订单详情
	logger.Info("订单详情: %v\n", order)
}

// 现货卖出
// params: orderT: limit(限价单)、market(市价单)、stop(止损单)、stop-limit(限价止损单)
func SpotOrderSell(symbol string, quantity string, price string, orderT string) {
	var orderType binance.OrderType
	client := binance.NewClient(util.ApiKey, util.SecretKey)
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
	order, err := client.NewCreateOrderService().
		Symbol(symbol).
		Side(side).
		Type(orderType).
		QuoteOrderQty(quantity).
		Price(price).
		Do(ctx)
	if err != nil {
		logger.Error("卖出失败: %v", err)
	}
	// 输出订单详情
	logger.Info("订单详情: %v", order)
}

// 现货买入测试
// params: orderT: limit(限价单)、market(市价单)、stop(止损单)、stop-limit(限价止损单)
func SpotOrderBuyTest(symbol string, quantity string, price string, orderT string) {
	// 测试
	//设置为测试网络（沙盒环境）模拟交易
	var orderType binance.OrderType
	client := binance.NewClient(util.ApiKeyTest, util.SecretKeyTest)
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
		logger.Error("创建订单失败: %v", err)
	}
	logger.Info("创建订单成功\n当前买入订单号: %d", resp.OrderID)

}

// 现货卖出测试
// params: orderT: limit(限价单)、market(市价单)、stop(止损单)、stop-limit(限价止损单)
func SpotOrderSellTest(symbol string, quantity string, price string, orderT string) {
	var orderType binance.OrderType
	client := binance.NewClient(util.ApiKeyTest, util.SecretKeyTest)
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
	client := binance.NewClient(util.ApiKey, util.SecretKey)
	cOrder, err := client.NewCancelOrderService().
		OrderID(orderId).
		Do(ctx)
	if err != nil {
		logger.Error("现货撤单失败: %v", err)
		return
	}
	logger.Info("现货撤单成功\n单号为: %d", cOrder.OrderID)
}

// 现货撤单测试
func SpotCancelOrderTest(orderId int64) {
	ctx := context.Background()
	client := binance.NewClient(util.ApiKeyTest, util.SecretKeyTest)
	client.BaseURL = binance.BaseAPITestnetURL
	cOrder, err := client.NewCancelOrderService().
		OrderID(orderId).
		Do(ctx)
	if err != nil {
		logger.Error("现货撤单失败: %v", err)
		return
	}
	logger.Info("现货撤单成功\n单号为: %d", cOrder.OrderID)
	fmt.Println(&cOrder)
}

// 获取交易历史
//
//	func GetTradeHistory() {
//		client := binance.NewClient(apiKey, util.SecretKey)
//	}
//
// 获取历史测试
func GetSpotHistoryTest() {
	client := binance.NewClient(util.ApiKeyTest, util.SecretKeyTest)
	client.BaseURL = binance.BaseAPITestnetURL
	ctx := context.Background()
	history, err := client.NewListOrdersService().Do(ctx)
	if err != nil {
		logger.Error("获取交易历史失败: ", err)
		return
	}
	logger.Info("获取交易历史成功: %v", history)
}

// 获得K线
func GetSpotKline(symbol string, Itl string) {
	spotDoneC := make(chan struct{})
	spotStopC := make(chan struct{})
	spotHandler := binance.WsKlineHandler(func(event *binance.WsKlineEvent) {
		fmt.Println(event.Kline)
	})
	spotErrHandler := binance.ErrHandler(func(err error) {
		if err != nil {
			logger.Error("k线加载器错误: %v", err)
			return
		}
	})
	spotDoneC, spotStopC, err := binance.WsKlineServe(symbol, Itl, spotHandler, spotErrHandler)
	if err != nil {
		logger.Error("获取现货K线失败: %v", err)
		return
	}
	go func() {
		for {
			select {
			case <-spotStopC:
				logger.Info("关闭现货数据流进程")
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
	client := binance.NewClient(util.ApiKey, util.SecretKey)
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

// 获得交易订单列表测试
func GetOrderListTest(symbol string) {
	// 测试
	ctx := context.Background()
	client := binance.NewClient(util.ApiKeyTest, util.SecretKeyTest)
	client.BaseURL = binance.BaseAPITestnetURL
	res, err := client.NewListOrdersService().
		Symbol(symbol).
		Do(ctx)
	if err != nil {
		logger.Error("获取订单列表失败: %v", err)
		return
	}
	for _, list := range res {
		resp, err := json.Marshal(list)
		if err != nil {
			fmt.Println("json编码错误")
		}
		fmt.Println(string(resp))
	}
}
