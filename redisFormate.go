package main

import "fmt"

func userTokenRDB(token string) string {
	str := fmt.Sprintf("token_user_%s", token)
	logger.Info(str)
	return str
}

func userRDB(username string) string {
	str := fmt.Sprintf("user_%s", username)
	logger.Info(str)
	return str
}

func loginRDB(username string) string {
	str := fmt.Sprintf("user_login_%s", username)
	logger.Info(str)
	return str
}

func logoutRDB(username string) string {
	str := fmt.Sprintf("user_logout_%s", username)
	logger.Info(str)
	return str
}

func userDailyCodeRDB(username, date string) string {
	str := fmt.Sprintf("user_code_%s_%s", date, username)
	logger.Info(str)
	return str
}
