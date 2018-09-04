package com

//首页海报
import (
	"JsGo/JsHttp"
	. "JsGo/JsLogger"
	"JsGo/JsStore/JsRedis"
	"JunSie/constant"
)

type poster struct {
	Title     string //标题
	Image     string //海报图
	ContentID string //文章id
	Status    string //状态
}

func InitHome() {
	JsHttp.WhiteHttps("/newposter", newPoster) // 新建首页海报
	JsHttp.WhiteHttps("/getposter", getPoster) //获取首页海报
	// JsHttp.Https("/modifyposter", ModPoster) //修改首页海报
	JsHttp.WhiteHttps("/gethomepro", getHomePro) //获取首页产品板块
}

//新建首页海报
func newPoster(s *JsHttp.Session) {
	posterlist := &[]poster{}
	e := s.GetPara(posterlist) //从前端获取
	if e != nil {
		Error(e.Error())
		s.Forward("1", e.Error(), "")
		return
	}
	eo := JsRedis.Redis_Sset(constant.POSTER, posterlist)
	if eo != nil {
		Error(eo.Error())
		s.Forward("3", e.Error(), "")
		return
	}
	s.Forward("0", "success", "")
}

//获取首页海报
func getPoster(s *JsHttp.Session) {
	data := make([]poster, 0)
	e := JsRedis.Redis_get(constant.POSTER, &data)
	if e != nil {
		Error(e.Error())
		s.Forward("2", e.Error(), "")
		return
	}
	s.Forward("0", "sucess", data)
}

// //修改首页海报
// func ModPoster(s.*JsHttp.Session) {
// data, e := JsRedis.Redis_Sget(constant.POSTER)
// if e !=nil {
// 	Error
// }
// s.Forward("0", "sucess", data)
// }
