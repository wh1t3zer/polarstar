package user

import (
	"context"
	"github.com/adshao/go-binance/v2"
	"polarstar/util"
)

var (
	logger = util.NewCustomLogger()
)

// GetUserDetail 获得现货用户详细信息
func GetSpotUserDetail() (userInfo *binance.Account) {
	client := binance.NewClient(util.ApiKey, util.SecretKey)
	ctx := context.Background()
	userInfo, err := client.NewGetAccountService().Do(ctx)
	if err != nil {
		logger.Error("获取现货用户信息失败: %v", err)
		return nil
	}
	// 币数量大于0的才展示
	var tmpInfo []binance.Balance
	for _, info := range userInfo.Balances {
		if info.Free > "0.00000000" {
			tmpInfo = append(tmpInfo, info)
		}
	}
	userInfo.Balances = tmpInfo
	logger.Info("获得用户信息成功")
	return userInfo
}
