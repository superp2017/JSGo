package main

import (
	"JsGo/JsHttp"
	"fmt"
	"sync"
	"time"
)

var mutex sync.Mutex
var g_u int = 1

func helloCallback(session *JsHttp.Session) {

	fmt.Printf("token = %s\n", session.Ctx.Input.Header("X-Token"))
	session.Ctx.Output.Header("X-Token", "123")

	session.Forward("1", "Hi, u!", nil)

}

func userLogin(session *JsHttp.Session) {
	fmt.Printf("token = %s\n", session.Ctx.Input.Header("X-Token"))

	type Para struct {
		Username string
		Password string
	}

	para := &Para{}
	e := session.GetPara(para)
	if e != nil {
		session.Forward("1", e.Error(), nil)
		return
	}

	token := fmt.Sprintf("%d", time.Now().Nanosecond())

	session.Ctx.Output.Header("X-Token", token)
	session.GenSession()
	session.Forward("0", "success", token)
}

func userInfo(session *JsHttp.Session) {
	type Token struct {
		Token string
	}
	token := &Token{}
	session.GetPara(&token)
	x_token := session.Ctx.Input.Header("X-Token")

	fmt.Printf("%s, %s\n", token.Token, x_token)
	if token.Token != x_token {
		session.Forward("1", "token not match", nil)
		return
	}

	type Entity struct {
		Roles  string
		Name   string
		Avatar string
	}

	session.Ctx.Output.Header("X-Token", token.Token)

	entity := &Entity{"admin", fmt.Sprintf("mengzhaofeng: %d", time.Now().Nanosecond()), "http://www.qqzhi.com/uploadpic/2014-09-23/000247589.jpg"}

	session.Forward("0", "success", entity)
}

func formIndex(session *JsHttp.Session) {
	x_token := session.Ctx.Input.Header("X-Token")

	if "" == x_token {
		session.Forward("1", "token not match", nil)
		return
	}

	session.Forward("0", "success", fmt.Sprintf("%d", time.Now().Nanosecond()))
}

func main() {
	JsHttp.Https("/hello", helloCallback)
	JsHttp.Https("/user/login", userLogin)
	JsHttp.Http("/user/info", userInfo)
	JsHttp.Http("/form/index", formIndex)

	JsHttp.Run()
}
