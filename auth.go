package main

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// User 用户基础信息
type User struct {
	NickName string `json:"nickname"`
	UserName string `json:"username"` //全部转为小写
	Email    string `json:"email"`
}

// UserToken 用户基础信息+token
type UserToken struct {
	User
	Token string `json:"token"`
}

const emailre string = `^([\w-\.]+)@ctrchina\.cn$`

// login 处理登录
// 校验邮件，设置cookie usertoken
func login(c *gin.Context) {

	email := c.PostForm("email")
	userName, err := getUserName(email)

	if err != nil {
		c.String(http.StatusBadRequest, "email error")
	}

	nickName := c.PostForm("nickname")

	user := &User{
		Email:    email,
		UserName: userName,
		NickName: nickName,
	}

	token := makeUserToken(user)
	userJSON, _ := json.Marshal(user)

	usertoken := &UserToken{
		User:  *user,
		Token: token,
	}
	userTokenJSON, _ := json.Marshal(usertoken)
	// 存储token
	rdb.SetEX(ctx, token+"_user", string(userJSON), time.Duration(expireUserAuth))
	// 存储user信息
	rdb.SetEX(ctx, userRDB(), string(userTokenJSON), 0)
	// 存储登录信息
	rdb.SAdd(ctx, "login_"+user.UserName, time.Now().Format("2006-01-02 15:04:05"))

	fmt.Printf(token)
	c.SetCookie("usertoken", token, 0, "/", "localhost", false, true)
	c.String(http.StatusOK, "OK")
}

// logout 处理登出
//
func logout(token string) {
	// 删除token
	rdb.Del(ctx, token+"_user")
}

// getUserName 根据邮箱获取username
func getUserName(email string) (username string, err error) {
	emailReg := regexp.MustCompile(emailre)

	if !emailReg.MatchString(email) {
		return "", &EmailFormateError{email: email}
	}

	return strings.ToLower(string(emailReg.FindAllSubmatch([]byte(email), -1)[0][1])), nil

}

// makeCookie 根据user的UserName生成Cookie
func makeUserToken(user *User) string {
	b := md5.Sum([]byte(user.UserName))
	return string(b[:])
}
