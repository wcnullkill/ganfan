package main

import "fmt"

func userTokenRDB(token string) string {
	return fmt.Sprintf("token_user_%s", token)
}

func userRDB() string {
	return fmt.Sprintf("user")
}

func loginRDB(username string) string {
	return fmt.Sprintf("user_login_%s", username)
}

func logoutRDB(username string) string {
	return fmt.Sprintf("user_logout_%s", username)
}
