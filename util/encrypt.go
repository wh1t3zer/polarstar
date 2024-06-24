package util

import (
	"crypto/ed25519"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"log"
	"os"
)

// ComputeEd25519 用来加密
func ComputeEd25519(message string) string {
	// 读取PEM文件
	pemData, err := os.ReadFile("./misc/Private_key.pem")
	if err != nil {
		log.Fatal("Failed to read PEM file: ", err)
	}
	// 解析PEM数据
	block, _ := pem.Decode(pemData)
	if block == nil {
		log.Fatal("Failed to decode PEM block")
	}
	// 解析私钥
	privateKey, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		log.Fatal("Failed to parse private key: ", err)
	}
	// 转换为Ed25519私钥类型
	ed25519PrivateKey, ok := privateKey.(ed25519.PrivateKey)
	if !ok {
		log.Fatal("Invalid private key type")
	}
	sign := ed25519.Sign(ed25519PrivateKey, []byte(message))
	b64Str := base64.StdEncoding.EncodeToString(sign)
	return b64Str
}
