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
	NickName string `json:"nickname" form:"nickname" xml:"nickname" binding:"required"`
	//全部转为小写
	UserName string `json:"username" form:"username" xml:"username" binding:"-"`
	Email    string `json:"email" form:"email" xml:"email" binding:"required"`
}

// UserToken 用户基础信息+token
type UserToken struct {
	User  `json:"user"`
	Token string `json:"token"`
}

const emailre string = `^([\w-\.]+)@ctrchina\.cn$`

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
	userJSON, _ := json.Marshal(user)

	usertoken := &UserToken{
		User:  user,
		Token: token,
	}
	userTokenJSON, _ := json.Marshal(usertoken)
	// 存储token
	rdb.SetEX(ctx, userTokenRDB(token), string(userJSON), time.Duration(expireUserAuth)).Result()

	// 存储user信息
	rdb.SetEX(ctx, userRDB(userName), string(userTokenJSON), time.Duration(expireUserAuth)).Result()

	// 存储登录信息
	rdb.SAdd(ctx, loginRDB(userName), time.Now().Format("2006-01-02 15:04:05")).Result()

	c.SetCookie("usertoken", token, 0, "/", "localhost", false, true)
	c.String(http.StatusOK, "OK")
}

// logout 处理登出
//
func logout(c *gin.Context) {

	cookie, _ := c.Request.Cookie("usertoken")
	// 删除token
	rdb.Del(ctx, userTokenRDB(cookie.Value))
	//删除cookie中usertoken信息
	c.SetCookie("usertoken", "", -1, "/", "localhost", false, true)
	c.String(http.StatusOK, "OK")
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
	return fmt.Sprintf("%x", b)
}

func checkAuth(c *gin.Context) {
	cookie, err := c.Request.Cookie("usertoken")
	if err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	_, err = rdb.Get(ctx, userTokenRDB(cookie.Value)).Result()
	//无效token
	if err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	//重置token超时
	rdb.Expire(ctx, userTokenRDB(cookie.Value), time.Duration(expireUserAuth))
	c.Next()
}
