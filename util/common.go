package util

import (
	"fmt"
	"github.com/e421083458/golang_common/lib"
	"io"
	"log"
	"os"
)

func init() {
	lib.InitModule("./conf/", []string{"base"})
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
