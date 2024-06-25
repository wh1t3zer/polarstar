// util/util.go
package util

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"os"
)

var (
	ApiKey        string
	SecretKey     string
	ApiKeyTest    string
	SecretKeyTest string
	BarkKey       string
	BarkURL       string
	BusinessURL   string
	BusinessKey   string
)

type baseConfig struct {
	Binance struct {
		MainNet struct {
			ApiKey    string `yaml:"apiKey"`
			SecretKey string `yaml:"secretKey"`
		} `yaml:"MainNet"`
		TestNet struct {
			ApiKeyTest    string `yaml:"apiKeyTest"`
			SecretKeyTest string `yaml:"secretKeyTest"`
		} `yaml:"TestNet"`
	} `yaml:"Binance"`
	Push struct {
		Bark struct {
			BarkKey string `yaml:"barkKey"`
			BarkUrl string `yaml:"barkURL"`
		} `yaml:"Bark"`
		BusWechat struct {
			BusinessURL string `yaml:"businessURL"`
			BusinessKey string `yaml:"businessKey"`
		} `yaml:"BusWechat"`
	} `yaml:"Push"`
}

var (
	Config baseConfig
)

// InitConfig initializes the configuration from the given file path and file name
func InitConfig(filePath string, configFile string) {
	configPath := fmt.Sprintf("%s/%s.yaml", filePath, configFile)
	yamlFile, err := os.ReadFile(configPath)
	if err != nil {
		logger.Error("读取文件失败: %v", err)
		os.Exit(1)
	}
	err = yaml.Unmarshal(yamlFile, &Config)
	if err != nil {
		logger.Error("解析配置文件失败: %v", err)
		os.Exit(1)
	}
	if &Config != nil {
		ApiKey = Config.Binance.MainNet.ApiKey
		SecretKey = Config.Binance.MainNet.SecretKey
		ApiKeyTest = Config.Binance.TestNet.ApiKeyTest
		SecretKeyTest = Config.Binance.TestNet.SecretKeyTest
		BarkKey = Config.Push.Bark.BarkKey
		BarkURL = Config.Push.Bark.BarkUrl
		BusinessKey = Config.Push.BusWechat.BusinessKey
		BusinessURL = Config.Push.BusWechat.BusinessURL
		logger.Info("初始化配置成功")
	} else {
		logger.Error("初始化配置失败")
	}
}
