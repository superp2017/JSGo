package com

import (
	"JsGo/JsHttp"
	. "JsGo/JsLogger"
	"JsGo/JsStore/JsRedis"
	"JunSie/constant"
)

type HomeTemp struct {
	Title     string   //标题
	Author    string   //作者
	Status    string   //状态
	Timestamp int64    //
	Type      string   //类型
	IDs       []string //ID
}

// popularity 受大众欢迎一、人气推荐（4）
//getpopularity
/*
建立一个放ID的表

getnewpro

四、专题精选 轮播图关联文章，specialPro
五、猜你喜欢（多个产品）？？？？guessLike

*/ // 	a = "popularity" // A, popularity 受大众欢迎一、人气推荐（4）
func NewPooularity(s *JsHttp.Session) {
	IDs := []string{}
	e := s.GetPara(IDs)
	if e != nil {
		Error(e.Error())
		s.Forward("1", e.Error(), "")
		return
	}
	eo := JsRedis.Redis_set(constant.POPULARITY, IDs)
	if eo != nil {
		Error(eo.Error())
		s.Forward("1", eo.Error(), "")
		return
	}
	s.Forward("0", "success", "")
}

////新建新品首发
func NewfreshPro(s *JsHttp.Session) {
	IDs := []string{}
	e := s.GetPara(IDs)
	if e != nil {
		Error(e.Error())
		s.Forward("1", e.Error(), "")
		return
	}
	eo := JsRedis.Redis_set(constant.FRESHPRO, IDs)
	if eo != nil {
		Error(eo.Error())
		s.Forward("2", eo.Error(), "")
		return
	}

}

//topic 主题板块(简化版)
func getHomePro(s *JsHttp.Session) {
	data := make([]string, 0)
	Topic := "" //板块单元
	er := s.GetPara(&Topic)
	if er != nil {
		Error(er.Error())
		s.Forward("1", er.Error(), "")
		return
	}
	err := JsRedis.Redis_get(Topic, &data)
	if err != nil {
		Error(err.Error())
		s.Forward("2", err.Error(), "")
		return
	}
	s.Forward("0", "success", GetMoreProductInfo(data))
}

// const (
// 	FRESHPRO   = "FreshPro"   //新品首发（6）
// 	POPULARITY = "Popularity" //人气推荐（4）受大众欢迎
// 	FLASHSALE  = "FlashSale"  //限时抢购（3）
// 	GUESSSLIKE = "GuessLike"  //guessLike猜你喜欢
// )

// func getHomePromore(s *JsHttp.Session) {//获取比较多的
// var e error
// e = JsRedis.Redis_get(Topic, &data)

// home := &HomeTemp{}
// e = JsRedis.Redis_hget("Template", Topic, home)
// articles := []string{}
// if home.Type == "0" {
// 	ids := home.IDs

// }
// ret := make(map[string]interface{})
// ret["A"] = home
// ret["B"] = articles
// s.Forward("0", "success", ret)

// if topic == "FreshPro" {
// 	e = JsRedis.Redis_get(constant.FRESHPRO, &data)
// }
// if topic == "FlashSale" {
// 	e = JsRedis.Redis_get(constant.FRESHPRO, &data)
// }
// if topic == "GuessLike" {
// 	e = JsRedis.Redis_get(constant.FRESHPRO, &data)
// }

//获取首页海报
// func getPoster(s *JsHttp.Session) {
// 	data := make([]poster, 0)
// 	e := JsRedis.Redis_get(constant.POSTER, &data)
// 	if e != nil {
// 		Error(e.Error())
// 		s.Forward("2", e.Error(), "")
// 		return
// 	}
// 	s.Forward("0", "sucess", data)
// }

// //修改首页海报
