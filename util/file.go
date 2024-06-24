package util

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"os"
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

func InitConfig(filePath string, configFile string) {
	configPath := fmt.Sprintf("%s/%s", filePath, configFile+".yaml")
	yamlFile, err := os.ReadFile(configPath)
	if err != nil {
		logger.ERROR("读取文件失败: %v", err)
		os.Exit(1)
	}
	err = yaml.Unmarshal(yamlFile, &Config)
	if err != nil {
		logger.ERROR("解析配置文件失败: %v", err)
		os.Exit(1)
	}
}
