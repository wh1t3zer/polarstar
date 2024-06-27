package UMContract

// U本位合约
import (
	"context"
	"fmt"
	"github.com/adshao/go-binance/v2/futures"
	"polarstar/util"
)

var (
	logger = util.NewCustomLogger()
)

// GetUMKline 合约K线
/*
	params
	symbol 交易对
	Itl	  时间段
*/
func GetUMKline(symbol string, Itl string) {
	futureDoneC := make(chan struct{})
	futureStopC := make(chan struct{})
	fsHandler := futures.WsContinuousKlineHandler(func(event *futures.WsContinuousKlineEvent) {
		fmt.Println(event.Kline)
	})
	fsErrHandler := futures.ErrHandler(func(err error) {
		if err != nil {
			logger.Error("U本位合约Error: %v", err)
			return
		}
	})
	fsKlineSubcribeArgs := &futures.WsContinuousKlineSubcribeArgs{
		Interval:     Itl,
		Pair:         symbol,
		ContractType: "perpetual",
	}
	futureDoneC, futureStopC, err := futures.WsContinuousKlineServe(fsKlineSubcribeArgs, fsHandler, fsErrHandler)
	if err != nil {
		logger.Error("获取U本位合约K线失败: %v", err)
		return
	}
	go func() {
		for {
			select {
			case <-futureStopC:
				logger.Info("关闭U本位通信流进程")
				close(futureStopC)
			}
		}
	}()
	<-futureDoneC
}

// ChangeUMLever 调整杠杆
/*
	params
	symbol 交易对
	lever  杠杆数
*/
func ChangeUMLever(symbol string, lever int) {
	ctx := context.Background()
	client := futures.NewClient(util.ApiKey, util.SecretKey)
	resp, err := client.NewChangeLeverageService().Symbol(symbol).Leverage(lever).Do(ctx)
	if err != nil {
		logger.Error("调整UM杠杆失败: %v", err)
		return
	}
	logger.Info("更改UM杠杆成功\n当前交易对: %s，杠杆数: %d，最大可开仓资金数: %s", resp.Symbol, resp.Leverage, resp.MaxNotionalValue)
}

// ContractOrderBuy 下单买入
/*
	params: orderType
	LIMIT 限价单
	MARKET 市价单
	STOP 止损限价单
	STOP_MARKET 止损市价单
	TAKE_PROFIT 止盈限价单
	TAKE_PROFIT_MARKET 止盈市价单
	TRAILING_STOP_MARKET 跟踪止损单
*/
func ContractOrderBuy(symbol string, quantity string, price string, orderT string, positionSide futures.PositionSideType) {
	var orderType futures.OrderType
	side := futures.SideTypeBuy
	ctx := context.Background()
	client := futures.NewClient(util.ApiKey, util.SecretKey)
	switch orderT {
	case "limit":
		orderType = futures.OrderTypeLimit
	case "market":
		orderType = futures.OrderTypeMarket
	case "stop":
		orderType = futures.OrderTypeStop
	}
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
		logger.Error("下单失败: %v", err)
	}
	fmt.Println(res)
}

// 下单买入测试
func ContractOrderBuyTest(symbol string, quantity string, price string, orderT string, positionSide futures.PositionSideType) {
	var orderType futures.OrderType
	side := futures.SideTypeBuy
	ctx := context.Background()
	futures.UseTestnet = true
	client := futures.NewClient(util.ApiKeyTest, util.SecretKeyTest)
	switch orderT {
	case "limit":
		orderType = futures.OrderTypeLimit
	case "market":
		orderType = futures.OrderTypeMarket
	case "stop":
		orderType = futures.OrderTypeStop
	}
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
		logger.Info("下单失败: %v", err)
		return
	}
	fmt.Println(res)
}

// 下单卖出
func ContractOrderSell() {
	//ctx := context.Background()
}

//// 下单卖出测试
//func ContractOrderSellTest() {
//	ctx := context.Background()
//	client := futures.NewClient(util.ApiKey, util.SecretKeyTest)
//	client.NewCreateOrderService().Symbol().Side().Type().Quantity().Price().TimeInForce().Do()
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
		logger.Error("撤单失败: %v", err)
		return
	}
	logger.Info("U本位撤单成功: %v", resp)
}

// 撤单测试
func CancelOrderTest(symbol string, orderId int64) {
	ctx := context.Background()
	client := futures.NewClient(util.ApiKeyTest, util.SecretKeyTest)
	futures.UseTestnet = true
	resp, err := client.NewCancelOrderService().
		Symbol(symbol).
		Do(ctx)
	if err != nil {
		logger.Error("撤单失败: %v", err)
		return
	}
	logger.Info("测试撤单成功: %v", resp)
}

// 批量撤单
func CancelOrderList(symbol string) {
	ctx := context.Background()
	client := futures.NewClient(util.ApiKey, util.SecretKey)
	err := client.NewCancelAllOpenOrdersService().
		Symbol(symbol).
		Do(ctx)
	if err != nil {
		logger.Error("批量撤单失败: ", err)
		return
	}
	logger.Info("U本位批量撤单成功")
}

// 批量撤单测试
func CancelOrderListTest(symbol string) {
	ctx := context.Background()
	client := futures.NewClient(util.ApiKey, util.SecretKey)
	futures.UseTestnet = true
	err := client.NewCancelAllOpenOrdersService().
		Symbol(symbol).
		Do(ctx)
	if err != nil {
		logger.Error("撤单失败: %v", err)
		return
	}
	logger.Info("批量撤单成功")
}

// 获取当前所有挂单的信息
func GetSymbolOrderList(symbol string) {
	ctx := context.Background()
	client := futures.NewClient(util.ApiKey, util.SecretKey)
	resp, err := client.NewListOpenOrdersService().
		Symbol(symbol).
		Do(ctx)
	if err != nil {
		logger.Error("获取全部挂单信息失败: %v", err)
		return
	}
	logger.Info("获得全部挂单信息成功: %v", resp)
}

// 获取当前目标挂单信息
func GetSymbolOrder(symbol string, orderId int64) {
	ctx := context.Background()
	client := futures.NewClient(util.ApiKey, util.SecretKey)
	resp, err := client.NewGetOpenOrderService().
		Symbol(symbol).
		OrderID(orderId).
		Do(ctx)
	if err != nil {
		logger.Error("获取单个挂单信息失败: %v", err)
		return
	}
	logger.Info("获得单个挂单信息成功: %v", resp)
}

// 获得成交历史订单
// 默认近7天、全部
func GetUMOrder(symbol string, orderId int64) {
	ctx := context.Background()
	client := futures.NewClient(util.ApiKey, util.SecretKey)
	resp, err := client.NewListAccountTradeService().
		Symbol(symbol).
		OrderID(orderId).
		Do(ctx)
	if err != nil {
		logger.Error("获得所有U本位订单失败: %v", err)
		return
	}
	logger.Info("获得U本位订单成功: %v", resp)
}

// 获得成交历史订单测试
func GetUMOrderTest(symbol string, orderId int64) {
	ctx := context.Background()
	client := futures.NewClient(util.ApiKeyTest, util.SecretKeyTest)
	resp, err := client.NewListAccountTradeService().
		Symbol(symbol).
		OrderID(orderId).
		Do(ctx)
	if err != nil {
		logger.Error("获得测试所有U本位订单失败: %v", err)
		return
	}
	logger.Info("获得测试U本位订单成功: %v", resp)
}
