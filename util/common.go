package util

import (
	"fmt"
	"io"
	"log"
	"os"
)

type Config struct {
	bc *BaseConfig
	cc *CommonConfig
}

func Init() *Config {
	Banner()
	var AConfig *Config
	BC := InitConfig("./conf/", "base")
	CC := InitCommon("./conf/", "common")
	AConfig = &Config{
		bc: BC,
		cc: CC,
	}
	return AConfig
}
func Banner() {
	// 打开文件
	file, err := os.Open("banner.txt")
	if err != nil {
		log.Fatal("Error opening file:", err)
		return
	}
	defer file.Close()

	// 读取文件内容
	content, err := io.ReadAll(file)
	if err != nil {
		log.Fatal("Error reading file:", err)
		return
	}

	// 输出文件内容到控制台
	fmt.Println(string(content))
}
