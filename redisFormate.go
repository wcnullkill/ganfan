package main

import "fmt"

// userTokenRDB 个人token
func userTokenRDB(token string) string {
	str := fmt.Sprintf("token_user_%s", token)
	logger.Info(str)
	return str
}

// userRDB 个人信息
func userRDB(username string) string {
	str := fmt.Sprintf("user_%s", username)
	logger.Info(str)
	return str
}

// loginRDB 个人登录记录列表
func loginRDB(username string) string {
	str := fmt.Sprintf("user_login_%s", username)
	logger.Info(str)
	return str
}

// logoutRDB 个人登出记录列表
func logoutRDB(username string) string {
	str := fmt.Sprintf("user_logout_%s", username)
	logger.Info(str)
	return str
}

// userDailyCodeRDB 每日个人预约结果
func userDailyReservationResultRDB(username, date string) string {
	str := fmt.Sprintf("user_code_%s_%s", date, username)
	logger.Info(str)
	return str
}

// dailyCodeRDB 每日预约列表
func dailyReservationRDB(date string) string {
	str := fmt.Sprintf("user_code_%s", date)
	logger.Info(str)
	return str
}
