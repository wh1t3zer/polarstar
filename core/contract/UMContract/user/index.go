package user

import (
	"context"
	"github.com/adshao/go-binance/v2/futures"
	"polarstar/util"
)

var (
	logger = util.NewCustomLogger()
)

// 获取U本位合约用户信息
func GetUMUserDetail() (userInfo *futures.Account) {
	client := futures.NewClient(util.ApiKey, util.SecretKey)
	ctx := context.Background()
	userInfo, err := client.NewGetAccountService().Do(ctx)
	if err != nil {
		logger.Error("获取U本位合约信息失败: %v", err)
		return nil
	}
	return userInfo
}
