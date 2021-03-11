package main

import (
	"crypto/md5"
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// User 用户基础信息
type User struct {
	NickName string `json:"nickname" form:"nickname" xml:"nickname" binding:"-"`
	//全部转为小写
	UserName string `json:"usern" form:"usern" xml:"usern" binding:"-"`
	Email    string `json:"email" form:"email" xml:"email" binding:"required"`
	P        int    `json:"p" xml:"p" bindiong:"-"`
}

// UserToken 用户基础信息+token
type UserToken struct {
	User  `json:"user"`
	Token string `json:"token"`
}

const emailre string = `^([\w-\.]+)@ctrchina\.cn$`
const cookieMaxAge int = 30 * 24 * 60 * 60

// login 处理登录
// 校验邮件，设置cookie usertoken
func login(c *gin.Context) {
	var user User

	if err := c.ShouldBind(&user); err != nil {
		c.String(http.StatusBadRequest, "user valid")
		return
	}

	userName, err := getUserName(user.Email)

	if err != nil {
		c.String(http.StatusBadRequest, "email error")
		return
	}

	user.UserName = userName

	token := makeUserToken(&user)

	rdb.SetEX(ctx, rdbUserToken(token), "", time.Duration(expireUserAuth)).Result()

	rdb.HMSet(ctx, rdbToken(token), "user", user.UserName, "email", user.Email, "nickname", user.NickName, "p", defaultP)

	rdb.SAdd(ctx, rdbTokenList(), token)

	// 存储登录信息
	rdb.SAdd(ctx, rdbLogin(userName), time.Now().Format("2006-01-02 15:04:05")).Result()

	// cookie设置，包含token和用户名
	c.SetCookie("token", token, cookieMaxAge, "/", "", false, true)
	c.SetCookie("user", userName, cookieMaxAge, "/", "", false, true)
	c.String(http.StatusOK, "OK")
}

// logout 处理登出
//
func logout(c *gin.Context) {

	cookie, _ := c.Request.Cookie("token")
	// 删除token
	rdb.Del(ctx, rdbUserToken(cookie.Value))
	userName, _ := rdb.HGet(ctx, rdbToken(cookie.Value), "user").Result()

	// 记录登出
	rdb.SAdd(ctx, rdbLogout(userName), time.Now().Format("2006-01-02 15:04:05")).Result()

	//删除cookie中usertoken信息
	c.SetCookie("token", "", -1, "/", "localhost", false, true)
	c.SetCookie("user", "", -1, "/", "localhost", false, true)
	c.String(http.StatusOK, "OK")
}

// getUserName 根据邮箱获取username
func getUserName(email string) (usern string, err error) {
	emailReg := regexp.MustCompile(emailre)

	if !emailReg.MatchString(email) {
		return "", &EmailFormateError{email: email}
	}

	return strings.ToLower(string(emailReg.FindAllSubmatch([]byte(email), -1)[0][1])), nil

}

// makeCookie 根据user的UserName生成Cookie
func makeUserToken(user *User) string {
	b := md5.Sum([]byte(user.UserName + "_ganfan"))
	return fmt.Sprintf("%x", b)
}

func checkAuth(c *gin.Context) {
	var token, user string
	cookies := c.Request.Cookies()

	for _, cookie := range cookies {
		switch cookie.Name {
		case "user":
			user = cookie.Value
		case "token":
			token = cookie.Value
		}
	}
	if len(token) <= 0 || len(user) <= 0 {
		c.String(http.StatusUnauthorized, "token或user不能未空")
		c.Abort()
		return
	}
	// 验证token是否有效
	_, err := rdb.Get(ctx, rdbUserToken(token)).Result()
	//无效token
	if err != nil {
		c.String(http.StatusUnauthorized, "token无效")
		c.Abort()
		return
	}
	userName, err := rdb.HGet(ctx, rdbToken(token), "user").Result()
	//信息被修改
	// todo  补充错误信息
	if user != userName {
		c.String(http.StatusForbidden, "信息被修改")
		c.Abort()
		return
	}

	//重置token超时
	rdb.Expire(ctx, rdbUserToken(token), time.Duration(expireUserAuth))
	c.Next()
}
