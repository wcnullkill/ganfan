package main

import (
	"context"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/mattn/go-colorable"
)

const (
	expireUserAuth int = int(time.Hour * 24 * 30)
)

var rdb *redis.Client
var ctx = context.Background()

func init() {
	rdb = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	initLogger()
}

func main() {
	router := RouterDefault()
	router.Run(":50001")
}

// RouterDefault get default router
func RouterDefault() *gin.Engine {
	// 启用gin的日志输出带颜色
	gin.ForceConsoleColor()
	// 替换默认Writer（关键步骤）
	gin.DefaultWriter = colorable.NewColorableStdout()
	router := gin.Default()

	v1 := router.Group("/v1")
	v1.Use(logToFile)

	authGroup := v1.Group("/auth")
	{
		authGroup.POST("", login)
		authGroup.Use(checkAuth)
		authGroup.DELETE("", logout)
	}

	orderGroup := v1.Group("/order")
	{
		orderGroup.Use(checkAuth)
		orderGroup.POST("", submitOrder)
		orderGroup.GET("", queryOrder)
	}

	return router
}
