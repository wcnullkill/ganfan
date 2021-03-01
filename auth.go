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

// User 用户
type User struct {
	NickName string `json:"nickname"`
	UserName string `json:"username"` //全部转为小写
	Email    string `json:"email"`
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
	js, _ := json.Marshal(user)

	rdb.SetEX(ctx, token+"_user", string(js), time.Duration(expireUserAuth))
	rdb.SAdd(ctx, "ganfan_login_"+user.UserName, time.Now().Format("2006-01-02 15:04:05"))
	fmt.Printf(token)

}

func logout(token string) {
	rdb.Del(ctx, token+"_user")
}

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
