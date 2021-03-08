package main

import "fmt"

// rdbUserToken kv,个人token,用于验证token是否有效
func rdbUserToken(token string) string {
	str := fmt.Sprintf("user_token_%s", token)
	logger.Info(str)
	return str
}

// rdbToken hash 根据token，存储个人信息，用于根据token信息，查询user信息
func rdbToken(token string) string {
	str := fmt.Sprintf("token_%s", token)
	logger.Info(str)
	return str
}

// rdbTokenList set 存储所有token
func rdbTokenList() string {
	str := fmt.Sprint("tokenlist")
	logger.Info(str)
	return str
}

// rdbLogin set 个人登录记录列表
func rdbLogin(username string) string {
	str := fmt.Sprintf("user_login_%s", username)
	logger.Info(str)
	return str
}

// rdbLogout set 个人登出记录列表
func rdbLogout(username string) string {
	str := fmt.Sprintf("user_logout_%s", username)
	logger.Info(str)
	return str
}

// rdbDailyReservationCode kv 每日预约结果,存储code
func rdbDailyReservationCode(token, date string) string {
	str := fmt.Sprintf("reservation_code_%s_%s", date, token)
	logger.Info(str)
	return str
}

// rdbDailyReservation set 每日预约列表，存储token信息
func rdbDailyReservation(date string) string {
	str := fmt.Sprintf("reservation_%s", date)
	logger.Info(str)
	return str
}
