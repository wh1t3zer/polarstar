package futures

import (
	"fmt"
	"github.com/adshao/go-binance/v2/futures"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"net/http"
	"polarstar/core/contract/UMContract"
	"strconv"
)

var (
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
)

// 改变杠杆
func ChangeLever(c *gin.Context) {
	symbol := c.PostForm("symbol")
	lever, err := strconv.Atoi(c.PostForm("lever"))
	if err != nil {
		return
	}
	info := UMContract.ChangeUMLever(symbol, lever)
	c.JSON(200, gin.H{"data": info})
}

// 下单买入
// 双向持仓
func ContractInOrder(c *gin.Context) {
	symbol := c.PostForm("symbol")
	quantity := c.PostForm("quantity")
	price := c.PostForm("price")
	orderT := c.PostForm("orderT")
	positionSide := c.PostForm("positionSide")
	side := c.PostForm("side")
	res, err := UMContract.ContractOrderBuy(symbol, quantity, price, orderT, futures.PositionSideType(positionSide), futures.SideType(side))
	if err != nil {
		c.JSON(200, err)
		return
	}
	c.JSON(200, res)
}

// 撤单批量
func CancelHoldOrderList(c *gin.Context) {
	symbol := c.PostForm("symbol")
	if symbol == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "symbol parameter is missing",
		})
		return
	}
	err := UMContract.CancelHoldOrderList(symbol)
	if err != nil {
		c.JSON(200, err)
	}
	c.JSON(200, "ok")
}

// 撤某一条挂单
func CancelHoldOrderDetail(c *gin.Context) {
	symbol := c.Query("symbol")
	orderStr := c.Query("orderID")
	if symbol == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "symbol parameter is missing",
		})
		return
	}
	if orderStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "symbol parameter is missing",
		})
		return
	} else {
		orderID, err := strconv.ParseInt(orderStr, 10, 64)
		if err != nil {
			fmt.Printf("转换失败: %v\n", err)
			return
		}
		resp, err := UMContract.CancelOrder(symbol, orderID)
		if err != nil {
			c.JSON(200, err)
		}
		c.JSON(200, resp)
	}
}

// 获得挂单信息
func GetHoldOrderList(c *gin.Context) {
	symbol := c.Query("symbol")
	if symbol == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "symbol parameter is missing",
		})
		return
	}
	res, err := UMContract.GetHoldOrderList(symbol)
	if err != nil {
		c.JSON(200, gin.H{
			"err": "挂单信息获取错误",
			"res": err,
		})
		return
	}
	c.JSON(200, res)
}

func GetHoldOrderDetail(c *gin.Context) {
	symbol := c.PostForm("symbol")
	orderStr := c.PostForm("orderID")
	if symbol == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "symbol parameter is missing",
		})
		return
	}
	if orderStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "orderID parameter is missing",
		})
		return
	} else {
		orderID, err := strconv.ParseInt(orderStr, 10, 64)
		if err != nil {
			fmt.Printf("转换失败: %v\n", err)
			return
		}
		resp, err := UMContract.GetHoldOrder(symbol, orderID)
		if err != nil {
			c.JSON(200, err)
		}
		c.JSON(200, resp)
	}
}

func GetUMOrder(c *gin.Context) {
	symbol := c.Query("symbol")
	if symbol == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "symbol parameter is missing",
		})
		return
	}
	resp, err := UMContract.GetUMOrder(symbol)
	if err != nil {
		c.JSON(200, err)
	}
	c.JSON(200, resp)
}
func GetUMOrderDetail(c *gin.Context) {
	symbol := c.PostForm("symbol")
	orderStr := c.PostForm("orderID")
	if symbol == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "symbol parameter is missing",
		})
		return
	}
	if orderStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "orderID parameter is missing",
		})
		return
	} else {
		orderID, err := strconv.ParseInt(orderStr, 10, 64)
		if err != nil {
			fmt.Printf("转换失败: %v\n", err)
			return
		}
		resp, err := UMContract.GetUMOrderDetail(symbol, orderID)
		if err != nil {
			c.JSON(200, err)
		}
		c.JSON(200, resp)
	}
}

// 获得K线
func GetUMKline(c *gin.Context) {
	symbol := "BTCUSDT" // 从消息中解析 symbol
	interval := "1m"    // 从消息中解析 interval
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		fmt.Printf("Failed to set websocket upgrade: %+v\n", err)
		return
	}
	defer conn.Close()
	// 解析消息并调用 GetUMKline
	stopC := make(chan struct{})
	go func() {
		UMContract.GetUMKline(symbol, interval, conn, stopC)
	}()
	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			fmt.Printf("Read error: %+v\n", err)
			close(stopC) // 通知 GetUMKline 函数停止
			break
		}
	}
	fmt.Println("WebSocket connection closed")
}
