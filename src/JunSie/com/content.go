package com

//内容||文章
import (
	"JsGo/JsBench/JsContent"
	"JsGo/JsHttp"
	"JsGo/JsLogger"
	"JsGo/JsStore/JsRedis"
	"JunSie/constant"
	"JunSie/util"
	"encoding/json"
	"fmt"
)

type XM_Contents struct {
	JsContent.Contents
}

func Init_content() {
	JsHttp.WhiteHttps("/newcontent", NewContent)                         //创建内容
	JsHttp.WhiteHttp("/querycontent", QueryContent)                      //查询内容
	JsHttp.WhiteHttp("/querymorecontents", QueryMoreContents)            //查询多个内容
	JsHttp.WhiteHttps("/modcontent", ModContent)                         //修改多个内容
	JsHttp.WhiteHttps("/delcontentdb", DelContentDB)                     //永久删除内容
	JsHttp.WhiteHttps("/delcontentmark", DelContentMark)                 //标记删除内容
	JsHttp.WhiteHttps("/updowncontent", UpdownContent)                   //上下架内容
	JsHttp.WhiteHttp("/getpagecontent", GetPageContent)                  //获取分页内容
	JsHttp.WhiteHttp("/getcontentnums", GetContentNums)                  //获取不同状态内容的数量
	JsHttp.WhiteHttps("/newcontentdraft", SaveContentDraft)              //保存内容草稿
	JsHttp.WhiteHttps("/removecontentdraft", QueryContentDraft)          //查询内容草稿
	JsHttp.WhiteHttp("/getcontentductdraft", RemoveContentDraft)         ///删除某一个内容草稿
	JsHttp.WhiteHttps("/removecontenttag", RemoveContentTag)             //移除产品的某一个标签
	JsHttp.WhiteHttp("/getcontenttags", GetContentTags)                  //获取产品的所有标签
	JsHttp.WhiteHttp("/getcontentstatics", GetContentStatics)            //获取内容统计
	JsHttp.WhiteHttps("/newcontentvisit", NewContentVisit)               //内容访问
	JsHttp.WhiteHttps("/newcontentpraise", NewContentPraise)             //内容点赞
	JsHttp.WhiteHttps("/newcontentattention", NewContentAttention)       //内容关注
	JsHttp.WhiteHttps("/cancelcontentpraise", RemoveContentPraise)       //内容取消点赞
	JsHttp.WhiteHttps("/cancelcontentattention", RemoveContentAttention) //内容取消关注
	JsHttp.WhiteHttps("/cancelcontentatvisit", RemoveArtVisit)           //删除内容访问（+）

}

func Init_contentMall() {
	JsHttp.WhiteHttps("/newcontent", NewContent)                         //创建内容
	JsHttp.WhiteHttps("/querycontent", QueryContent)                     //查询内容
	JsHttp.WhiteHttps("/querymorecontents", QueryMoreContents)           //查询多个内容
	JsHttp.WhiteHttps("/modcontent", ModContent)                         //修改多个内容
	JsHttp.WhiteHttps("/delcontentdb", DelContentDB)                     //永久删除内容
	JsHttp.WhiteHttps("/delcontentmark", DelContentMark)                 //标记删除内容
	JsHttp.WhiteHttps("/updowncontent", UpdownContent)                   //上下架内容
	JsHttp.WhiteHttps("/getpagecontent", GetPageContent)                 //获取分页内容
	JsHttp.WhiteHttps("/getcontentnums", GetContentNums)                 //获取不同状态内容的数量
	JsHttp.WhiteHttps("/newcontentdraft", SaveContentDraft)              //保存内容草稿
	JsHttp.WhiteHttps("/removecontentdraft", QueryContentDraft)          //查询内容草稿
	JsHttp.WhiteHttps("/getcontentductdraft", RemoveContentDraft)        ///删除某一个内容草稿
	JsHttp.WhiteHttps("/removecontenttag", RemoveContentTag)             //移除产品的某一个标签
	JsHttp.WhiteHttps("/getcontenttags", GetContentTags)                 //获取产品的所有标签
	JsHttp.WhiteHttps("/getcontentstatics", GetContentStatics)           //获取内容统计
	JsHttp.WhiteHttps("/newcontentvisit", NewContentVisit)               //内容访问
	JsHttp.WhiteHttps("/newcontentpraise", NewContentPraise)             //内容点赞
	JsHttp.WhiteHttps("/newcontentattention", NewContentAttention)       //内容关注
	JsHttp.WhiteHttps("/cancelcontentpraise", RemoveContentPraise)       //内容取消点赞
	JsHttp.WhiteHttps("/cancelcontentattention", RemoveContentAttention) //内容取消关注
	JsHttp.WhiteHttps("/cancelcontentatvisit", RemoveArtVisit)           //删除内容访问（+）
}

//新建内容
func NewContent(s *JsHttp.Session) {
	type Para struct {
		CT   XM_Contents //内容结构
		Tags []string    //标签列表
	}
	st := &Para{}
	if err := s.GetPara(st); err != nil {
		s.Forward("1", err.Error(), st)
		return
	}
	st.CT.ID = util.IDer(constant.DB_Content)
	st.CT.CreatDate = util.CurDate()

	if st.CT.UID == "" {
		info := "NewContent,UID is empty\n"
		JsLogger.Error(info)
		s.Forward("1", info, nil)
		return
	}

	if st.CT.Status == "" {
		st.CT.Status = constant.Status_ON
	}
	if err := JsRedis.Redis_hset(constant.DB_Content, st.CT.ID, &st.CT); err != nil {
		s.Forward("1", err.Error(), st)
		return
	}

	//建立内容搜索索引
	go creatContentSearchIndex(&st.CT)

	go appendToGlobalContent(st.CT.ID)
	if len(st.Tags) > 0 {
		go Content2Tag(st.CT.ID, st.Tags)
	}
	TagLinkA(st.CT.ID, st.Tags)
	// if erro := tage.TagLinkA(st.CT.ID, st.Tags); erro != nil {
	// 	info := "Tage contect fail:" + erro.Error()
	// 	JsLogger.Error(info)
	// 	s.Forward("1", info, nil)
	// 	return
	// }
	s.Forward("0", "success", st)
}

//查询内容
func QueryContent(s *JsHttp.Session) {
	type info struct {
		ID string //id
	}
	st := &info{}
	if err := s.GetPara(st); err != nil {
		s.Forward("1", err.Error(), nil)
		return
	}
	data, err := Getcontent(st.ID)
	if err != nil {
		s.Forward("1", err.Error(), nil)
		return
	}
	s.Forward("0", "success", data)
}

//查询多个内容
func QueryMoreContents(s *JsHttp.Session) {
	type Info struct {
		IDs []string ////产品id列表
	}
	st := &Info{}
	if err := s.GetPara(st); err != nil {
		info := "QueryMoreContents,GetPara:" + err.Error()
		JsLogger.Error(info)
		s.Forward("1", info, nil)
		return
	}
	if len(st.IDs) < 0 {
		info := "QueryMoreContents param IDs is empty\n"
		JsLogger.Error(info)
		s.Forward("1", info, nil)
		return
	}
	s.Forward("0", "success", getMoreContents(st.IDs))
}

//修改内容
func ModContent(s *JsHttp.Session) {
	type Para struct {
		CT   XM_Contents //内容结构
		Tags []string    //标签列表
	}
	st := &Para{}
	if err := s.GetPara(st); err != nil {
		info := "ModContent,GetPara:" + err.Error()
		JsLogger.Error(info)
		s.Forward("1", info, nil)
		return
	}
	data := &XM_Contents{}
	if err := JsRedis.Redis_hget(constant.DB_Content, st.CT.ID, data); err != nil {
		info := "ModContent Redis_get:" + err.Error()
		JsLogger.Error(info)
		s.Forward("1", info, data)
		return
	}
	d, e := json.Marshal(st.CT)
	if e != nil {
		info := e.Error()
		JsLogger.Error(info)
		s.Forward("1", info, nil)
		return
	}
	if err := json.Unmarshal(d, data); err != nil {
		info := err.Error()
		JsLogger.Error(info)
		s.Forward("1", info, nil)
		return
	}
	if err := JsRedis.Redis_hset(constant.DB_Content, data.ID, data); err != nil {
		info := err.Error()
		JsLogger.Error(info)
		s.Forward("1", info, nil)
		return
	}
	go changeContentTag(data.ID, st.Tags)
	s.Forward("0", "success", data)
}

//数据库直接删除内容
func DelContentDB(s *JsHttp.Session) {
	type Para struct {
		ID string
	}
	st := &Para{}
	if err := s.GetPara(st); err != nil {
		s.Forward("1", err.Error(), nil)
		return
	}
	if st.ID == "" {
		info := "DelContentDB param ID is empty\n"
		JsLogger.Error(info)
		s.Forward("1", info, nil)
		return
	}
	if ok, e := JsRedis.Redis_hexists(constant.DB_Content, st.ID); e == nil && ok {
		if err := JsRedis.Redis_hdel(constant.DB_Content, st.ID); err != nil {
			info := "DelContentDB" + err.Error()
			JsLogger.Error(info)
			s.Forward("1", info, nil)
			return
		}
	}
	go delFromGlobalContent(st.ID)
	s.Forward("0", "success", nil)
}

//标记删除内容
func DelContentMark(s *JsHttp.Session) {
	type Info struct {
		ID string
	}
	st := &Info{}
	if err := s.GetPara(st); err != nil {
		info := "DelContentMark,GetPara:" + err.Error()
		JsLogger.Error(info)
		s.Forward("1", info, nil)
		return
	}
	if st.ID == "" {
		info := "DelContentMark param ID is empty\n"
		JsLogger.Error(info)
		s.Forward("1", info, nil)
		return
	}
	data := &XM_Contents{}
	if err := JsRedis.Redis_hget(constant.DB_Content, st.ID, data); err != nil {
		info := "DelContentMark,Redis_hset:" + err.Error()
		JsLogger.Error(info)
		s.Forward("1", info, nil)
		return
	}
	data.DelTag = true
	if err := JsRedis.Redis_hset(constant.DB_Content, st.ID, data); err != nil {
		info := "DelContentMark,Redis_hset:" + err.Error()
		JsLogger.Error(info)
		s.Forward("1", info, nil)
		return
	}
	go delFromGlobalContent(st.ID)
	s.Forward("0", "success", data)
}

//上下架内容
func UpdownContent(s *JsHttp.Session) {
	type Info struct {
		ID   string //id
		IsUp bool   //是否上架
	}
	st := &Info{}
	if err := s.GetPara(st); err != nil {
		info := "UpdownContent,GetPara:" + err.Error()
		JsLogger.Error(info)
		s.Forward("1", info, nil)
		return
	}
	if st.ID == "" {
		info := "UpdownContent param ID is empty\n"
		JsLogger.Error(info)
		s.Forward("1", info, nil)
		return
	}
	data := &XM_Contents{}
	if err := JsRedis.Redis_hget(constant.DB_Content, st.ID, data); err != nil {
		info := "UpdownContent,Redis_get:" + err.Error()
		JsLogger.Error(info)
		s.Forward("1", info, nil)
		return
	}
	if st.IsUp {
		data.Status = constant.Status_ON
	} else {
		data.Status = constant.Status_OFF
	}
	if err := JsRedis.Redis_hset(constant.DB_Content, st.ID, data); err != nil {
		info := "UpdownContent,Redis_get:" + err.Error()
		JsLogger.Error(info)
		s.Forward("1", info, nil)
		return
	}
	s.Forward("0", "success", data)
}

//查询单个内容
func Getcontent(id string) (*XM_Contents, error) {
	data := &XM_Contents{}
	return data, JsRedis.Redis_hget(constant.DB_Content, id, data)
}

//获取多个内容信息（摘要）
func getMoreContents(ids []string) []*JsContent.ContentsAbs {
	data := []*JsContent.ContentsAbs{}
	fmt.Printf("id = %v\n", ids)
	for _, v := range ids {
		retData := &JsContent.ContentsAbs{}
		err := JsRedis.Redis_hget(constant.DB_Content, v, retData)
		if err != nil {
			JsLogger.Error(err.Error())
			continue
		}
		data = append(data, retData)
	}
	return data
}
