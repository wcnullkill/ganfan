package main

import (
	"context"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/mattn/go-colorable"
	"github.com/robfig/cron/v3"
)

const (
	expireUserAuth time.Duration = time.Hour * 24 * 30 //30天
	expireCode     time.Duration = time.Hour * 24 * 5  //5天
	codeMax        int           = 999999
	codeMin        int           = 100000
)

var rdb *redis.Client
var ctx = context.Background()
var randArray [codeMax - codeMin]int

func init() {
	rdb = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	if _, err := rdb.Ping(ctx).Result(); err != nil {
		panic(err)
	}

	initLogger()

	// 创建抽奖定时任务
	c := cron.New()
	c.AddFunc(jobCronS, execute)
	c.Start()

	//初始化随机数组
	for i := 0; i < codeMax-codeMin; i++ {
		randArray[i] = i + codeMin
	}
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

	codeGroup := v1.Group("/code")
	{
		codeGroup.Use(checkAuth)
		codeGroup.POST("", submitCode)
		codeGroup.GET("", queryCode)
	}

	return router
}
