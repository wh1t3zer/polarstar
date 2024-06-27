package user

import (
	"context"
	"github.com/adshao/go-binance/v2/delivery"
	"polarstar/util"
)

var (
	logger = util.NewCustomLogger()
)

// 获取币本位合约用户信息
func GetCMUserDetail() (userInfo *delivery.Account) {
	client := delivery.NewClient(util.ApiKey, util.SecretKey)
	ctx := context.Background()
	userInfo, err := client.NewGetAccountService().Do(ctx)
	if err != nil {
		logger.Error("获取币本位合约信息失败: %v", err)
		return
	}
	return userInfo
}
