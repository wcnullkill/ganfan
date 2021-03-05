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

// submitCode 预约
func submitCode(c *gin.Context) {

	c.String(http.StatusOK, "OK")
}

// queryOrder 查询本人当日预约结果
func queryCode(c *gin.Context) {
	date := time.Now().Format("20060102")
	user, _ := c.Request.Cookie("user")

	codeJSON, err := rdb.Get(ctx, userDailyReservationResultRDB(user.Value, date)).Result()

	if err != nil {
		c.String(http.StatusNotFound, "code not found")
		return
	}

	c.Header("Content-Type", "application/json")
	c.String(http.StatusOK, codeJSON)

}

//excute 抽奖
func excute() {

}

// getMembers
func getMembers() []*member {
	//查询当天预约名单
	//date := time.Now().Format("20060102")
	return nil
}
