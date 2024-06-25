package util

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/adshao/go-binance/v2"
	"log"
)

/*
params : symbol 货币交易对
*/
func GetRateLimit() {
	client := binance.NewClient(ApiKey, SecretKey)
	ctx := context.Background()
	exchangeInfo, err := client.NewExchangeInfoService().Do(ctx)
	if err != nil {
		log.Fatalf("Error getting exchange info :%v", err)
	}
	rateLimitJson, err := json.MarshalIndent(exchangeInfo.RateLimits, "", "  ")
	if err != nil {
		log.Fatalf("将速率限制信息转换为 JSON 时出错: %v", err)
	}
	fmt.Println(string(rateLimitJson))

}
