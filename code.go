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

const (
	defaultP int    = 50                                   //初始中奖概率
	jobCron  int64  = int64(time.Hour*11 + time.Minute*20) //抽奖定时任务时间的unix,11点20分
	jobCronS string = "20 11 * * * *"                      //抽奖定时任务时间字符串形式
)

// submitCode 预约
// 将预约名单插入redis
func submitCode(c *gin.Context) {
	date := time.Now().Format("20060102")
	token, _ := c.Request.Cookie("token")

	// 检查是否已经生成结果
	code, _ := rdb.Get(ctx, rdbDailyReservationCode(token.Value, date)).Result()
	if len(code) > 0 {
		c.String(http.StatusForbidden, "已生成预约结果，本次预约无效")
		return
	}
	// 加入预约名单
	result, _ := rdb.SAdd(ctx, rdbDailyReservation(date), token).Result()
	if result != 1 {
		c.String(http.StatusForbidden, "重复预约")
		return
	}

	c.String(http.StatusOK, "OK")
}

// queryOrder 查询本人当日预约结果
func queryCode(c *gin.Context) {
	date := time.Now().Format("20060102")
	token, _ := c.Request.Cookie("token")

	code, err := rdb.Get(ctx, rdbDailyReservationCode(token.Value, date)).Result()

	if err != nil {
		c.String(http.StatusNotFound, "code not found")
		return
	}

	c.String(http.StatusOK, code)

}

// excute 抽奖
func excute() {

}

// getMembers 查询出当天预约名单
func getMembers() []*member {
	// 查询当天预约名单
	date := time.Now().Format("20060102")
	tokens, _ := rdb.SMembers(ctx, rdbDailyReservation(date)).Result()
	members := make([]*member, 0, len(tokens))

	for _, token := range tokens {
		results, _ := rdb.HMGet(ctx, rdbToken(token), "username", "email", "p").Result()
		members = append(members, &member{
			name:  results[0].(string),
			email: results[1].(string),
			p:     results[2].(int),
		})
	}
	// TODO	pipeline批量处理

	return members
}
