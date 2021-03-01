package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
)

const (
	expireUserAuth int = int(time.Hour * 24 * 30)
)

var rdb *redis.Client
var ctx = context.Background()

func init() {
	rdb = redis.NewClient(&redis.Options{
		Addr:     "10.20.12.80:6379",
		Password: "",
		DB:       0,
	})
}

func main() {
	err := rdb.Set(ctx, "wc", "niubi", 0).Err()
	if err != nil {
		panic(err)
	}
	v, err := rdb.Get(ctx, "wc").Result()
	if err != nil {
		panic(err)
	}

	fmt.Println(v)
}

// RouterDefault get default router
func RouterDefault() *gin.Engine {
	router := gin.Default()

	v1 := router.Group("/v1")

	authGroup := v1.Group("/auth")
	{
		authGroup.POST("/", login)
	}

	// orderGroup := v1.Group("/order")
	// {
	// 	orderGroup.Post("/")
	// }

	return router
}

func checkAuth(c *gin.Context) {
	userToken, err := c.Request.Cookie("UserToken")

	//未找到cookie
	if err == http.ErrNoCookie {
		c.String(http.StatusUnauthorized, "Unauthorized")
	}

	_, err = rdb.Get(ctx, userToken.Value+"_user").Result()
	//无效token
	if err != nil {
		c.String(http.StatusUnauthorized, "Unauthorized")
	}
	//重置token超时
	rdb.Expire(ctx, userToken.Value+"_user", time.Duration(expireUserAuth))
}
