package util

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

// 自定义日志记录器
type CustomLogger struct {
	logger *log.Logger
}

// 创建一个新的自定义日志记录器
func NewCustomLogger() *CustomLogger {
	timestamp := time.Now().Format("2006-01-02")
	logFile, err := os.OpenFile("logs/"+"polarstar_"+timestamp+".log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatalf("Failed to open log file: %v", err)
	}
	// 设置日志输出到控制台和文件
	multiWriter := io.MultiWriter(os.Stdout, logFile)
	return &CustomLogger{
		logger: log.New(multiWriter, "", 0),
	}
}

func getCallerInfo() string {
	pc, file, line, ok := runtime.Caller(3)
	if !ok {
		return "unknown"
	}
	fn := runtime.FuncForPC(pc)
	if fn == nil {
		return "unknown"
	}
	fileName := filepath.Base(file)
	fileName = strings.TrimSuffix(fileName, ".go")
	pkgName := strings.TrimPrefix(filepath.Dir(file), "/")
	pkgName = strings.ReplaceAll(pkgName, "/", ".")
	funcName := fn.Name()
	return fmt.Sprintf("%s.%s:%d %s", pkgName, fileName, line, funcName)
}

// 格式化日志输出
func (c *CustomLogger) log(level string, msg string, args ...interface{}) {
	timestamp := time.Now().Format("2006-01-02T15:04:05.000-07:00")
	pid := os.Getpid()
	callerInfo := getCallerInfo()

	// 格式化日志消息
	formattedMsg := fmt.Sprintf(msg, args...)
	c.logger.Printf("%s %s %d --- [main] %s : %s", timestamp, level, pid, callerInfo, formattedMsg)
}

// 日志记录方法
func (c *CustomLogger) Info(msg string, args ...interface{}) {
	c.log("INFO", msg, args...)
}

func (c *CustomLogger) Debug(msg string, args ...interface{}) {
	c.log("DEBUG", msg, args...)
}

func (c *CustomLogger) ERROR(msg string, args ...interface{}) {
	c.log("ERROR", msg, args...)
}
