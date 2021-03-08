package main

import (
	"math/rand"
	"net/http"
	"strconv"
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
	result, _ := rdb.SAdd(ctx, rdbDailyReservation(date), token.Value).Result()
	if result != 1 {
		c.String(http.StatusForbidden, "重复预约")
		//logger.Info(err)
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

// execute 抽奖
func execute() {

	date := time.Now().Format("20060102")
	logger.Infof("抽奖%s开始", date)
	// 预约名单
	members := getMembers()
	// 所有预约都能加权
	for _, member := range members {
		p, _ := rdb.HGet(ctx, rdbToken(member.Token), "p").Result()
		i, _ := strconv.Atoi(p)
		// 50,75,87,93,99,99
		i = i + (100-i)>>1
		rdb.HSet(ctx, rdbToken(member.Token), "p", i)
	}

	// 中奖名单
	results := make([]*member, size)
	if len(members) <= size {
		copy(results, members)
	} else {
		// 抽奖
		m := randAsListNode(members)
		copy(results, m)
	}
	rand.Seed(time.Now().Unix())
	rand.Shuffle(len(randArray), func(i, j int) {
		randArray[i], randArray[j] = randArray[j], randArray[i]
	})

	// 处理中奖名单
	for index, member := range results {
		// 将中奖token和code放入redis
		rdb.Set(ctx, rdbDailyReservationCode(member.Token, date), randArray[index], expireCode)
		// 重置中奖者P，并更新至user信息
		rdb.HSet(ctx, rdbToken(member.Token), "p", defaultP)
	}
	logger.Infof("抽奖%s结束", date)
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
			Token: token,
			User: User{
				UserName: results[0].(string),
				Email:    results[1].(string),
				P:        results[2].(int),
			},
		})
	}
	// TODO	pipeline批量处理

	return members
}
