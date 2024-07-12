package UMContract

// U本位合约
import (
	"context"
	"fmt"
	"github.com/adshao/go-binance/v2/futures"
	"github.com/gorilla/websocket"
	"polarstar/util"
	"time"
)

var (
	logger = util.NewCustomLogger()
)

// GetUMKline 合约K线
/*
	params
	symbol 交易对
	Itl	  时间段
	conn  websocket参数
	stopC 停止信号
*/
func GetUMKline(symbol string, Itl string, conn *websocket.Conn, stopC chan struct{}) {
	futureDoneC := make(chan struct{})
	futureStopC := make(chan struct{})
	fsHandler := futures.WsContinuousKlineHandler(func(event *futures.WsContinuousKlineEvent) {
		err := conn.WriteJSON(event.Kline)
		if err != nil {
			logger.Error("发送Kline数据到前端失败: %v", err)
			err := conn.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("Error: %v", err)))
			if err != nil {
				return
			}
		}
	})
	fsErrHandler := futures.ErrHandler(func(err error) {
		if err != nil {
			logger.Error("U本位合约Error: %v", err)
			err := conn.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("Error: %v", err)))
			if err != nil {
				return
			}
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
		err := conn.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("Error: %v", err)))
		if err != nil {
			return
		}
		return
	}
	go func() {
		<-stopC
		close(futureStopC) // 通知 Binance WebSocket 处理停止
		logger.Info("关闭U本位通信流进程")
	}()
	<-futureDoneC
}

// ChangeUMLever 调整杠杆
/*
	params
	symbol 交易对
	lever  杠杆数
*/
func ChangeUMLever(symbol string, lever int) *futures.SymbolLeverage {
	ctx := context.Background()
	client := futures.NewClient(util.BC.Binance.MainNet.ApiKey, util.BC.Binance.MainNet.SecretKey)
	resp, err := client.NewChangeLeverageService().Symbol(symbol).Leverage(lever).Do(ctx)
	if err != nil {
		logger.Error("调整UM杠杆失败: %v", err)
		return nil
	}
	logger.Info("更改UM杠杆成功\n当前交易对: %s，杠杆数: %d，最大可开仓资金数: %s", resp.Symbol, resp.Leverage, resp.MaxNotionalValue)
	return resp
}

// ContractOrderBuy 下单买入
/*
	tips: orderType
	LIMIT 限价单
	MARKET 市价单
	STOP 止损限价单
	STOP_MARKET 止损市价单
	TAKE_PROFIT 止盈限价单
	TAKE_PROFIT_MARKET 止盈市价单
	TRAILING_STOP_MARKET 跟踪止损单
	params
	symbol 交易对
	quantity 下单数量
	price  价格
	orderT	订单类型
	positionSide 持仓模式	单/双
	side   买入/做多 卖出/做空
*/
func ContractOrderBuy(symbol string, quantity string, price string, orderT string, positionSide futures.PositionSideType, side futures.SideType) (res *futures.CreateOrderResponse, err error) {
	var orderType futures.OrderType
	ctx := context.Background()
	//client := futures.NewClient(util.BC.Binance.MainNet.ApiKey, util.BC.Binance.MainNet.SecretKey)
	client := futures.NewClient(util.BC.Binance.TestNet.ApiKeyTest1, util.BC.Binance.TestNet.SecretKeyTest1)
	client.BaseURL = "https://testnet.binancefuture.com"
	switch orderT {
	case "limit":
		orderType = futures.OrderTypeLimit
	case "market":
		orderType = futures.OrderTypeMarket
	case "stop":
		orderType = futures.OrderTypeStop
	}
	res, err = client.NewCreateOrderService().
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
		return nil, err
	}
	updateTime := time.Unix(0, res.UpdateTime*int64(time.Millisecond))
	formattedTime := updateTime.Format("2006-01-02 15:04:05")
	logger.Info("下单成功:\n 订单ID: %d,订单方向: %s,交易对: %s,下单时间: %s", res.OrderID, res.Side, res.Symbol, formattedTime)
	return res, nil
}

// 平仓
func ContractOrderSell() {
	//ctx := context.Background()

}

// CancelOrder 撤单
/*
	params
	symbol 交易对
	orderId 订单ID
*/
func CancelOrder(symbol string, orderId int64) (resp *futures.CancelOrderResponse, err error) {
	ctx := context.Background()
	//client := futures.NewClient(util.BC.Binance.MainNet.ApiKey, util.BC.Binance.MainNet.SecretKey)
	client := futures.NewClient(util.ApiKeyTest1, util.SecretKeyTest1)
	client.BaseURL = "https://testnet.binancefuture.com"
	resp, err = client.NewCancelOrderService().
		Symbol(symbol).
		OrderID(orderId).
		Do(ctx)
	if err != nil {
		logger.Error("撤单失败: %v", err)
		return nil, err
	}
	logger.Info("U本位撤单成功: %v", resp)
	return resp, nil
}

// CancelHoldOrderList 批量撤单
/*
	params
	symbol 交易对
*/
func CancelHoldOrderList(symbol string) (err error) {
	ctx := context.Background()
	client := futures.NewClient(util.BC.Binance.TestNet.ApiKeyTest1, util.BC.Binance.TestNet.SecretKeyTest1)
	client.BaseURL = "https://testnet.binancefuture.com"
	//client := futures.NewClient(util.BC.Binance.MainNet.ApiKey, util.BC.Binance.MainNet.SecretKey)
	err = client.NewCancelAllOpenOrdersService().
		Symbol(symbol).
		Do(ctx)
	if err != nil {
		logger.Error("批量撤单失败:\n %v", err)
		return err
	}
	logger.Info("U本位批量撤单成功")
	return nil
}

// GetHoldOrderList 获取当前所有挂单的信息
/*
	params
	symbol 交易对
*/
func GetHoldOrderList(symbol string) (resp []*futures.Order, err error) {
	ctx := context.Background()
	client := futures.NewClient(util.BC.Binance.TestNet.ApiKeyTest1, util.BC.Binance.TestNet.SecretKeyTest1)
	client.BaseURL = "https://testnet.binancefuture.com"
	//client := futures.NewClient(util.BC.Binance.MainNet.ApiKey, util.BC.Binance.MainNet.SecretKey)
	resp, err = client.NewListOpenOrdersService().
		Symbol(symbol).
		Do(ctx)
	if err != nil {
		logger.Error("获取全部挂单信息失败: %v", err)
		return nil, err
	}
	logger.Info("获得全部挂单信息成功")
	return resp, nil
}

// GetHoldOrder 获取当前所选目标挂单信息
/*
	params
	symbol 交易对
	oderId 订单号
*/
func GetHoldOrder(symbol string, orderId int64) (resp *futures.Order, err error) {
	ctx := context.Background()
	client := futures.NewClient(util.BC.Binance.TestNet.ApiKeyTest1, util.BC.Binance.TestNet.SecretKeyTest1)
	client.BaseURL = "https://testnet.binancefuture.com"
	//client := futures.NewClient(util.BC.Binance.MainNet.ApiKey, util.BC.Binance.MainNet.SecretKey)
	resp, err = client.NewGetOpenOrderService().
		Symbol(symbol).
		OrderID(orderId).
		Do(ctx)
	if err != nil {
		logger.Error("获取单个挂单信息失败: %v", err)
		return nil, err
	}
	logger.Info("获得单个挂单信息成功: %v", resp)
	return resp, nil
}

// GetUMOrder 获得成交历史订单
// 默认近7天、全部
/*
	params
	symbol 交易对
*/
func GetUMOrder(symbol string) (resp []*futures.Order, err error) {
	ctx := context.Background()
	//client := futures.NewClient(util.BC.Binance.MainNet.ApiKey, util.BC.Binance.MainNet.SecretKey)
	client := futures.NewClient(util.BC.Binance.TestNet.ApiKeyTest1, util.BC.Binance.TestNet.SecretKeyTest1)
	client.BaseURL = "https://testnet.binancefuture.com"
	resp, err = client.NewListOrdersService().
		Symbol(symbol).
		Do(ctx)
	if err != nil {
		logger.Error("获得所有U本位订单失败: %v", err)
		return nil, err
	}
	logger.Info("获得U本位订单成功")
	return resp, nil
}

// GetUMOrderDetail 获得某个历史订单
/*
	params
	symbol 交易对
	orderID 订单号
*/
func GetUMOrderDetail(symbol string, orderID int64) (resp []*futures.Order, err error) {
	ctx := context.Background()
	//client := futures.NewClient(util.BC.Binance.MainNet.ApiKey, util.BC.Binance.MainNet.SecretKey)
	client := futures.NewClient(util.BC.Binance.TestNet.ApiKeyTest1, util.BC.Binance.TestNet.SecretKeyTest1)
	client.BaseURL = "https://testnet.binancefuture.com"
	resp, err = client.NewListOrdersService().
		Symbol(symbol).
		OrderID(orderID).
		Do(ctx)
	if err != nil {
		logger.Error("获得所有U本位订单失败: %v", err)
		return nil, err
	}
	logger.Info("获得U本位订单成功")
	return resp, nil
}

func GetHistoryKline(symbol string, Itl string, limit int) (resList []futures.Kline, total int) {
	ctx := context.Background()
	total = 0
	client := futures.NewClient(util.BC.Binance.TestNet.ApiKeyTest1, util.BC.Binance.TestNet.SecretKeyTest1)
	resp, err := client.NewKlinesService().
		Symbol(symbol).
		Interval(Itl).
		Limit(limit).
		Do(ctx)
	if err != nil {
		logger.Error("创建K线失败:\n %v", err)
		return nil, total
	}
	for _, res := range resp {
		resList = append(resList, *res)
		total++
	}
	return resList, total
}
