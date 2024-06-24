package user

import (
	"context"
	"fmt"
	"github.com/adshao/go-binance/v2"
	"polarstar/aa"
)

// GetUserDetail 获得用户详细信息
func GetUserDetail() (userInfo *binance.Account) {
	client := binance.NewClient(aa.ApiKey, aa.SecretKey)
	ctx := context.Background()
	userInfo, err := client.NewGetAccountService().Do(ctx)
	if err != nil {
		fmt.Println("获取用户信息失败：", err)
		return nil
	}
	return userInfo
}
