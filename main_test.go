package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/go-redis/redis/v8"
)

type testUser struct {
	Email string `json:"email" xml:"email"`
}

var (
	username = "testuser"
	token    = "f692af1e749cfb6a8933b27a5fe18973"
	email    = username + "@ctrchina.cn"
	ts       *httptest.Server
)

func init() {
	setup()
	ts = httptest.NewServer(RouterDefault())
	fmt.Println("test init")
}

func setup() {
	rdb = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	if _, err := rdb.Ping(ctx).Result(); err != nil {
		panic(err)
	}

	initLogger()

	logger.SetOutput(os.Stderr)
}

type setupFunc func(token string) error

// setupLogin 正常登录
// 重复登录暂时没有问题
func setupLogin() error {

	u := &testUser{
		Email: email,
	}
	c := http.DefaultClient
	us, _ := json.Marshal(u)
	reader := bytes.NewReader(us)
	req, err := http.NewRequest("POST", ts.URL+"/v1/auth", reader)
	req.Header["Content-Type"] = []string{"application/json"}
	if err != nil {
		return fmt.Errorf("setupLogin error:%v", err)
	}
	resp, err := c.Do(req)
	if err != nil {
		return fmt.Errorf("setupLogin error:%v", err)
	}

	if resp.StatusCode != 200 {
		return errors.New(fmt.Sprint("setupLogin login error"))
	}
	return nil
}

// setupRemoveCode 移除当日中奖结果
func setupRemoveCode(token string) error {
	date := time.Now().Format("20060102")
	// 移除中奖结果
	_, err := rdb.Del(ctx, rdbDailyReservationCode(token, date)).Result()
	return err
}

// setupRemoveReservation 移除当日预约
func setupRemoveReservation(token string) error {
	date := time.Now().Format("20060102")
	_, err := rdb.SRem(ctx, rdbDailyReservation(date), token).Result()
	return err
}

// setupAddCode 移除当日中奖结果后，重新生成中奖结果(code=123456)
func setupAddCode(token string) error {
	date := time.Now().Format("20060102")
	setupRemoveCode(token)
	_, err := rdb.Set(ctx, rdbDailyReservationCode(token, date), "123456", time.Duration(20*time.Second)).Result()
	return err
}

// setupAddReservation 移除当日预约列表，重新加入预约列表
func setupAddReservation(token string) error {
	date := time.Now().Format("20060102")
	setupRemoveReservation(token)
	_, err := rdb.SAdd(ctx, rdbDailyReservation(date), token).Result()
	return err
}

// TestLogin 测试登录，响应post auth
func TestLogin(t *testing.T) {
	table := []struct {
		email       string
		code        int
		msg         string
		cookieUser  string
		cookieToken string
	}{
		{email, 200, "OK", username, token},                    //正常登录
		{email[:len(email)-1], 400, "email error", "", ""},     //邮箱不是ctr邮箱
		{"", 400, "user valid", "", ""},                        //邮箱为空
		{"asdf.fdsa@@ctrchina.cn", 400, "email error", "", ""}, //邮箱不符合常规邮箱正则
	}
	c := http.DefaultClient
	for _, v := range table {
		u := &testUser{
			Email: v.email,
		}
		us, _ := json.Marshal(u)
		reader := bytes.NewReader(us)
		req, err := http.NewRequest("POST", ts.URL+"/v1/auth", reader)
		req.Header["Content-Type"] = []string{"application/json"}
		checkError(t, err)
		resp, err := c.Do(req)
		checkError(t, err)
		body, err := ioutil.ReadAll(resp.Body)
		checkError(t, err)
		cookies := resp.Cookies()
		var cookiesUser, cookieToken string

		for _, cookie := range cookies {
			switch cookie.Name {
			case "user":
				cookiesUser = cookie.Value
			case "token":
				cookieToken = cookie.Value
			}
		}
		if resp.StatusCode != v.code || string(body) != v.msg || cookieToken != v.cookieToken || cookiesUser != v.cookieUser {
			t.Errorf("%v", v)
		}
	}

}

// TestCheckAuth 测试auth中间件，目前除了Login外，
// 其余情况都需要通过auth中间件.
// auth的测试，统一在这里了.
func TestCheckAuth(t *testing.T) {

	table := []struct {
		token string
		user  string
		code  int
		msg   string
	}{
		{token + "1", username, 401, "token无效"}, //错误的token值
		{token, username + "1", 403, "信息被修改"},   //user与token不匹配
		{"", username, 401, "token或user不能未空"},   //token为空
		{token, "", 401, "token或user不能未空"},      //user为空
	}

	for _, v := range table {
		c := http.DefaultClient
		req, err := http.NewRequest("DELETE", ts.URL+"/v1/auth", nil)
		checkError(t, err)
		req.Header.Set("Cookie", fmt.Sprintf("token=%s;user=%s", v.token, v.user))
		resp, err := c.Do(req)
		checkError(t, err)
		body, err := ioutil.ReadAll(resp.Body)
		checkError(t, err)
		if resp.StatusCode != v.code || string(body) != v.msg {
			t.Errorf("%v", v)
		}

	}
}

// TestLogout 测试登出，响应delet auth
func TestLogout(t *testing.T) {

	// 前置条件
	setupLogin()

	table := []struct {
		token string
		user  string
		code  int
		msg   string
	}{
		{token, username, 200, "OK"},
	}

	for _, v := range table {
		c := http.DefaultClient
		req, err := http.NewRequest("DELETE", ts.URL+"/v1/auth", nil)
		checkError(t, err)
		req.Header.Set("Cookie", fmt.Sprintf("token=%s;user=%s", v.token, v.user))
		resp, err := c.Do(req)
		checkError(t, err)
		body, err := ioutil.ReadAll(resp.Body)
		checkError(t, err)
		if resp.StatusCode != v.code || string(body) != v.msg {
			t.Errorf("%v", v)
		}

	}
}

// TestReserv 预约，响应post code
func TestReserv1(t *testing.T) {
	table := []struct {
		funcs []setupFunc
		code  int
		msg   string
	}{
		{
			funcs: []setupFunc{setupRemoveCode, setupRemoveReservation},
			code:  200,
			msg:   "OK",
		}, //正常预约
		{
			funcs: []setupFunc{setupRemoveReservation, setupAddCode},
			code:  403,
			msg:   "已生成预约结果，本次预约无效",
		}, //已生成预约结果，本次预约无效，比如过了11点20后，再执行预约请求
		{
			funcs: []setupFunc{setupRemoveCode, setupAddReservation},
			code:  403,
			msg:   "重复预约",
		}, //重复预约，比如11点10分预约一次，11点11分预约一次
	}

	for _, v := range table {
		// 固定前置条件
		setupLogin()
		//按顺序执行自定义前置条件
		for _, f := range v.funcs {
			err := f(token)
			checkError(t, err)
		}
		c := http.DefaultClient
		req, err := http.NewRequest("POST", ts.URL+"/v1/code", nil)

		req.Header.Set("cookie", fmt.Sprintf("token=%s;user=%s", token, username))
		resp, err := c.Do(req)
		checkError(t, err)

		body, err := ioutil.ReadAll(resp.Body)
		checkError(t, err)
		if resp.StatusCode != v.code || string(body) != v.msg {
			t.Fail()
		}
	}

}

// TestReserv 查询预约，响应get code
func TestQueryReservetion(t *testing.T) {
	table := []struct {
		funcs []setupFunc
		code  int
		msg   string
	}{
		{
			funcs: []setupFunc{setupAddCode},
			code:  200,
			msg:   "123456",
		}, //正常查询，即当日已生成预约结果
		{
			funcs: []setupFunc{setupRemoveReservation, setupRemoveCode},
			code:  404,
			msg:   "code not found",
		}, //没有预约，且当日还未生成预约结果
		{
			funcs: []setupFunc{setupAddReservation, setupRemoveCode},
			code:  404,
			msg:   "code not found",
		}, //有预约，且当日还未生成预约结果，比如在11点20以前查询
		//不存在没有预约，但有user预约结果的情况
	}

	for _, v := range table {
		// 固定前置条件
		setupLogin()
		//按顺序执行自定义前置条件
		for _, f := range v.funcs {
			err := f(token)
			checkError(t, err)
		}
		c := http.DefaultClient
		req, err := http.NewRequest("GET", ts.URL+"/v1/code", nil)

		req.Header.Set("cookie", fmt.Sprintf("token=%s;user=%s", token, username))
		resp, err := c.Do(req)
		checkError(t, err)

		body, err := ioutil.ReadAll(resp.Body)
		checkError(t, err)
		if resp.StatusCode != v.code || string(body) != v.msg {
			fmt.Println(v)
			t.Fail()
		}
	}

}

func checkError(t *testing.T, err error) {
	if err != nil {
		t.Error(err)
	}
}
