package main

import (
	"fmt"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

const logFilePath = "./gin.log"

var logger = logrus.New()

// initLogger 初始化logger
func initLogger() {

	src, err := createLoggerFile(logFilePath)

	if err != nil {
		fmt.Println("err", err)
	}
	//logger.Out = os.Stdout
	logger.Out = src
	logger.SetLevel(logrus.InfoLevel)
	logger.SetFormatter(&logrus.TextFormatter{})
}

// logToFile 日志输出到log文件
func logToFile(c *gin.Context) {

	// 开始时间
	startTime := time.Now()

	// 处理请求
	c.Next()

	// 结束时间
	endTime := time.Now()

	// 执行时间
	latencyTime := endTime.Sub(startTime)

	// 请求方式
	reqMethod := c.Request.Method

	// 请求路由
	reqURI := c.Request.RequestURI

	// 状态码
	statusCode := c.Writer.Status()

	// 请求IP
	clientIP := c.ClientIP()

	// 日志格式
	logger.Infof("| %3d | %13v | %15s | %s | %s |",
		statusCode,
		latencyTime,
		clientIP,
		reqMethod,
		reqURI,
	)

}

func createLoggerFile(path string) (f *os.File, err error) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		f, err = os.Create(path)
		if err != nil {
			panic(err)
		}
	} else {
		f, err = os.OpenFile(path, os.O_RDWR|os.O_APPEND, 0660)
	}
	return
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}
