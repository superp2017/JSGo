package main

import (
	"JsGo/JsHttp"
	. "JsGo/JsLogger"
	"JunSie/com"
	"fmt"
)

func InitLogin() {
	JsHttp.WhiteHttps("/login", Login) //登录密码账号
	JsHttp.Https("/exchange", Exchange)
}

func Login(s *JsHttp.Session) {
	type Para struct {
		Name     string
		Password string
	}
	para := &Para{}
	e := s.GetPara(para)
	if e != nil {
		Error(e.Error())
		s.Forward("1", e.Error(), "")
		return
	}

	if para.Name == "" || para.Password == "" {
		str := fmt.Sprintf("Login failed:Name=%s,Pwd=%\n", para.Name, para.Password)
		Error(str)
		s.Forward("1", str, nil)
		return
	}

	admin, err := com.GetAccountInfo(para.Name)
	if err != nil {
		Error(e.Error())
		s.Forward("2", e.Error(), "")
		return
	}
	if para.Name == admin.Account && para.Password == admin.Password {
		s.MarkSession()
		s.Forward("0", "success", admin)
	} else {
		s.Forward("3", "name password not match", "")
	}
}

func Exchange(s *JsHttp.Session) {
	type Para struct {
		Name  string
		Token string
	}

	para := &Para{}
	e := s.GetPara(para)
	if e != nil {
		Error(e.Error())
		s.Forward("1", e.Error(), "")
		return
	}

	if s.IsExpired(para.Token) {
		Error("token is expired")
		s.Forward("4", "token is expired", "")
		return
	}

	admin, err := com.GetAccountInfo(para.Name)
	if err != nil {
		Error(e.Error())
		s.Forward("2", e.Error(), "")
		return
	}

	if para.Name == admin.Account {
		s.MarkSession()
		s.Forward("0", "success", admin)
	} else {
		s.Forward("3", "name password not match", "")
	}
}
