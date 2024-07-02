// util/util.go
package util

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"os"
)

var (
	ApiKey         string
	SecretKey      string
	ApiKeyTest     string
	SecretKeyTest  string
	ApiKeyTest1    string
	SecretKeyTest1 string
	BarkKey        string
	BarkURL        string
	BusinessURL    string
	BusinessKey    string
)

type BaseConfig struct {
	Binance struct {
		MainNet struct {
			ApiKey    string `yaml:"apiKey"`
			SecretKey string `yaml:"secretKey"`
		} `yaml:"MainNet"`
		TestNet struct {
			ApiKeyTest     string `yaml:"apiKeyTest"`
			SecretKeyTest  string `yaml:"secretKeyTest"`
			ApiKeyTest1    string `yaml:"apiKeyTest1"`
			SecretKeyTest1 string `yaml:"secretKeyTest1"`
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

type CommonConfig struct {
	Env struct {
		IsProxy bool   `yaml:"isProxy"`
		Address string `yaml:"Address"`
	} `yaml:"Env"`
	Http struct {
		Addr           string   `yaml:"addr"`
		ReadTimeout    int      `yaml:"read_timeout"`
		WriteTimeout   int      `yaml:"write_timeout"`
		MaxHeaderBytes int      `yaml:"max_header_bytes"`
		AllowIP        []string `yaml:"allow_ip"`
	} `yaml:"http"`
}

var (
	BC BaseConfig
	CC CommonConfig
)

// InitConfig initializes the configuration from the given file path and file name
func InitConfig(filePath string, configFile string) *BaseConfig {
	configPath := fmt.Sprintf("%s/%s.yaml", filePath, configFile)
	yamlFile, err := os.ReadFile(configPath)
	if err != nil {
		logger.Error("读取文件失败: %v", err)
		os.Exit(1)
	}
	err = yaml.Unmarshal(yamlFile, &BC)
	if err != nil {
		logger.Error("解析配置文件失败: %v", err)
		os.Exit(1)
	}
	if &BC != nil {
		logger.Info("交易基本配置初始化成功")
		return &BC
	} else {
		logger.Error("交易基本配置初始化失败")
		return nil
	}
}
func InitCommon(filePath string, configFile string) *CommonConfig {
	configPath := fmt.Sprintf("%s/%s.yaml", filePath, configFile)
	yamlFile, err := os.ReadFile(configPath)
	if err != nil {
		logger.Error("读取文件失败: %v", err)
		os.Exit(1)
	}
	err = yaml.Unmarshal(yamlFile, &CC)
	if err != nil {
		logger.Error("解析配置文件失败: %v", err)
		os.Exit(1)
	}
	if &CC != nil {
		logger.Info("服务器基本配置初始化成功")
		return &CC
	} else {
		logger.Error("服务器基本配置初始化失败")
		return nil
	}
}
