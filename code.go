package main

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// Code 兑换码
type Code struct {
	Code     int `json:"code" xml:"code"`         //6位数字兑换码
	Expirate int `json:"expirate" xml:"expirate"` //有效期
	Status   int `json:"status" xml:"status"`     //兑换状态
}

func submitCode(c *gin.Context) {
	time.Now()
	c.String(http.StatusOK, "OK")
}

// queryOrder 查询本人当日订单
func queryCode(c *gin.Context) {
	date := time.Now().Format("20060102")
	user, _ := c.Request.Cookie("user")

	codeJSON, err := rdb.Get(ctx, userDailyCodeRDB(user.Value, date)).Result()

	if err != nil {
		c.String(http.StatusNotFound, "code not found")
		return
	}

	c.Header("Content-Type", "application/json")
	c.String(http.StatusOK, codeJSON)

}
