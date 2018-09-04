package com

import (
	"JsGo/JsHttp"
	. "JsGo/JsLogger"
	"JsGo/JsStore/JsRedis"
	"JunSie/constant"
)

func InitShow() {

	JsHttp.WhiteHttps("/getshowartone", GetShowArtOne) //获取首页产品板块
}

type Showone struct {
	Title     string   //标题
	UserHead  string   //用户头像
	SubTitle  string   //副标题
	Brief     string   //简介
	Question  string   //问题
	Showpic   string   //展示使用图
	Timestamp int64    //时间
	Type      string   //类型
	IDs       []string //ID
}

func GetShowArtOne(s *JsHttp.Session) {
	data := make([]Showone, 0)
	e := JsRedis.Redis_get(constant.SHOWARTONE, &data)
	if e != nil {
		Error(e.Error())
		s.Forward("2", e.Error(), "")
		return
	}
	s.Forward("0", "sucess", data)
}
