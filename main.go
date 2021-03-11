package main

import (
	"context"
	"fmt"
	"net/http"
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
	fmt.Println("init over")
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
	v1.Use(cors)

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

func cors(c *gin.Context) {
	method := c.Request.Method
	origin := c.Request.Header.Get("Origin") //请求头部
	if origin != "" {

		c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
		c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		c.Header("Access-Control-Allow-Headers", "Authorization, Content-Length, X-CSRF-Token, Token,session")
		c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers")
		c.Header("Access-Control-Max-Age", "172800")
		c.Header("Access-Control-Allow-Credentials", "true")
	}

	//允许类型校验
	if method == "OPTIONS" {
		c.JSON(http.StatusOK, "ok!")
	}

	c.Next()
}
