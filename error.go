package main

import "fmt"

// EmailFormateError 邮件格式错误
type EmailFormateError struct {
	email string
}

func (err EmailFormateError) String() string {
	return "email error"
}
func (err EmailFormateError) Error() string {
	return fmt.Sprintf("%s,%s", err.String(), err.msg())
}

func (err EmailFormateError) msg() string {
	return fmt.Sprintf("{\"email\":\"%s\"}", err.email)
}
