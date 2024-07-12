package user

import (
	"context"
	"github.com/adshao/go-binance/v2/futures"
	"polarstar/util"
)

var (
	logger = util.NewCustomLogger()
)

// GetUMUserDetail 获取U本位合约用户信息
func GetUMUserDetail() (userInfo *futures.Account) {
	ctx := context.Background()
	//client := futures.NewClient(util.BC.Binance.MainNet.ApiKey, util.BC.Binance.MainNet.SecretKey)
	client := futures.NewClient(util.BC.Binance.TestNet.ApiKeyTest1, util.BC.Binance.TestNet.SecretKeyTest1)
	client.BaseURL = "https://testnet.binancefuture.com"
	userInfo, err := client.NewGetAccountService().Do(ctx)
	if err != nil {
		logger.Error("获取U本位合约信息失败: %v", err)
		return nil
	}
	//"maintMargin": "254.95913021" 维持保证金

	return userInfo
}
