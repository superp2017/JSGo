package com

import (
	"JsGo/JsHttp"
	. "JsGo/JsLogger"
	"JsGo/JsStore/JsRedis"
	"JunSie/constant"
	"time"
)

func InitTag() {
	JsHttp.Https("/newtag", newTag)
	JsHttp.Http("/queryalltag", queryAllTag)
	JsHttp.Https("/deltag", delTag)

}

type Tagnr struct { //标签内容
	Labe      string //标签明字
	CountPro  int    //产品数量
	CountArt  int    //文章（内容）数量
	TagFinPro string // tag find product (标签查找产品;放最新的)
	TagFinArt string // tag find article（标签查找文章；放最新的）
	Stamp     int64  //时间
}

//新建标签
func newTag(s *JsHttp.Session) {
	type Para struct { //parameter参数
		Tag string
	}
	para := &Para{}
	if e := s.GetPara(para); e != nil {
		Error(e.Error())
		s.Forward("1", e.Error(), "")
		return
	}
	if para.Tag == "" {
		s.Forward("1", "tag=nil", "")
		return
	}

	b, e := JsRedis.Redis_hexists(constant.TAG, para.Tag) //exist存在
	if e != nil {
		Error(e.Error())
		s.Forward("2", e.Error(), "")
		return
	}
	if b {
		s.Forward("0", "success", "")
		return
	}

	tag := &Tagnr{}
	tag.Labe = para.Tag
	tag.CountPro = 0
	tag.CountArt = 0
	tag.TagFinPro = ""
	tag.TagFinArt = ""
	tag.Stamp = time.Now().Unix()

	e = JsRedis.Redis_hset(constant.TAG, para.Tag, tag)
	if e != nil {
		Error(e.Error())
		s.Forward("3", e.Error(), "")
		return
	}
	s.Forward("0", "success", "")
}

//查询所有标签
func queryAllTag(s *JsHttp.Session) {
	keys, e := JsRedis.Redis_hkeys(constant.TAG)
	if e != nil {
		Error(e.Error())
		s.Forward("1", e.Error(), "")
		return
	}
	tags := make(map[string]interface{})
	for _, v := range keys {
		tag := &Tagnr{}
		tags[v] = tag
	}
	e = JsRedis.Redis_hmget(constant.TAG, &tags)
	if e != nil {
		Error(e.Error())
		s.Forward("2", e.Error(), "")
		return
	}
	s.Forward("0", "success", tags)
}

//删除标签
func delTag(s *JsHttp.Session) {
	type Para struct {
		Tag string
	}
	para := &Para{}
	e := s.GetPara(para)
	if e != nil {
		Error(e.Error())
		s.Forward("1", e.Error(), "")
		return
	}
	e = JsRedis.Redis_hdel(constant.TAG, para.Tag)
	if e != nil {
		Error(e.Error())
		s.Forward("3", e.Error(), "")
		return
	}
	s.Forward("0", "success", "")
}
