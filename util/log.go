package util

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"
)

// 自定义日志记录器
type CustomLogger struct {
	mu         sync.Mutex
	logger     *log.Logger
	logFile    *os.File // 当前日志文件
	lastLogDay int      // 上次创建日志的日期（天）
}

// 创建一个新的自定义日志记录器
func NewCustomLogger() *CustomLogger {
	logger := log.New(os.Stdout, "", log.LstdFlags)
	c := &CustomLogger{
		logger:     logger,
		lastLogDay: time.Now().Day(),
	}
	err := c.createLogFile()
	if err != nil {
		return nil
	}
	return c
}

// 检查是否需要创建新的日志文件
func (c *CustomLogger) checkCreateLogFile() error {
	currentDay := time.Now().Day()
	if currentDay != c.lastLogDay {
		// 关闭当前日志文件
		if c.logFile != nil {
			c.logFile.Close()
		}

		// 创建新的日志文件
		err := c.createLogFile()
		if err != nil {
			return fmt.Errorf("创建新的日志文件失败: %v", err)
		}
		c.lastLogDay = currentDay
	}
	return nil
}

// 创建日志文件
func (c *CustomLogger) createLogFile() error {
	// 根据当前日期创建日志文件名
	timestamp := time.Now().Format("2006-01-02")
	fileName := fmt.Sprintf("polarstar_%s.log", timestamp)
	filePath := filepath.Join("logs", fileName)

	// 打开日志文件
	logFile, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		return fmt.Errorf("打开日志文件失败: %v", err)
	}

	// 设置日志输出到控制台和文件
	c.logger.SetOutput(io.MultiWriter(os.Stdout, logFile))
	c.logFile = logFile

	return nil
}

// 获取调用者信息
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
	return fmt.Sprintf("%s.%s:%d", pkgName, fileName, line)
}

// 格式化日志输出
func (c *CustomLogger) log(level string, msg string, args ...interface{}) {
	pid := os.Getpid()
	callerInfo := getCallerInfo()

	// 检查是否需要创建新的日志文件
	c.mu.Lock()
	defer c.mu.Unlock()
	if err := c.checkCreateLogFile(); err != nil {
		c.logger.Printf("ERROR %d --- [main] util.CustomLogger: %s", pid, err.Error())
		return
	}

	// 格式化日志消息
	formattedMsg := fmt.Sprintf(msg, args...)
	c.logger.Printf("%s %d --- [main] %s : %s", level, pid, callerInfo, formattedMsg)
}

// 日志记录方法
func (c *CustomLogger) Info(msg string, args ...interface{}) {
	c.log("INFO", msg, args...)
}

func (c *CustomLogger) Debug(msg string, args ...interface{}) {
	c.log("DEBUG", msg, args...)
}

func (c *CustomLogger) Error(msg string, args ...interface{}) {
	c.log("ERROR", msg, args...)
}

func (c *CustomLogger) Fatal(msg string, args ...interface{}) {
	c.log("FATAL", msg, args...)
	os.Exit(1)
}
