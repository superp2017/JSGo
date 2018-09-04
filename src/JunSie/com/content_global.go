package com

import (
	"JsGo/JsBench/JsContent"
	"JsGo/JsHttp"
	"JsGo/JsLogger"
	"JsGo/JsStore/JsRedis"
	"JunSie/constant"
	"fmt"
)

//获取分状态的分页产品列表
func GetPageContent(s *JsHttp.Session) {
	type Info struct {
		Status string //状态
		SIndex int    //启始索引
		Size   int    //个数
	}
	st := &Info{}
	if err := s.GetPara(st); err != nil {
		info := "getPageProducts,GetPara:" + err.Error()
		JsLogger.Error(info)
		s.Forward("1", info, nil)
		return
	}
	if st.Status == "" || (st.Status != constant.Status_ON && st.Status != constant.Status_OFF) || st.SIndex < 0 || st.Size <= 0 {
		info := fmt.Sprintf("getPageProducts param error,Status=%s,SIndex=%d,Size=%d\n", st.Status, st.SIndex, st.Size)
		JsLogger.Error(info)
		s.Forward("1", info, nil)
		return
	}
	data := []*JsContent.ContentsAbs{}
	ids, err := getGlobalContent()
	if err != nil {
		s.Forward("1", err.Error(), nil)
		return
	}
	list := getMoreContents(ids)
	for _, v := range list {
		if v.Status == st.Status {
			data = append(data, v)
		}
	}
	if st.SIndex < len(data) {
		if st.Size < len(data) && (st.SIndex+st.Size) < len(data) {
			s.Forward("0", "success", data[st.SIndex:st.SIndex+st.Size])
			return
		}
		s.Forward("0", "success", data[st.SIndex:])
		return
	}
	s.Forward("0", "success", nil)
}

//获取各种状态的内容的数量
func GetContentNums(s *JsHttp.Session) {
	num := make(map[string]int)
	ids, err := getGlobalContent()
	if err != nil {
		s.Forward("0", err.Error(), num)
		return
	}
	list := getMoreContents(ids)
	for _, v := range list {
		if v.Status != "" {
			if n, ok := num[v.Status]; ok {
				num[v.Status] = n + 1
			} else {
				num[v.Status] = 1
			}
		}
	}

	if _,ok:=num[constant.Status_ON];!ok{
		num[constant.Status_ON]=0
	}
	if _,ok:=num[constant.Status_OFF];!ok{
		num[constant.Status_OFF]=0
	}

	s.Forward("0", "success", num)
}

//添加一个产品id全局的内容id列表
func appendToGlobalContent(ID string) error {
	return JsRedis.Redis_Sset(constant.KEY_GlobalContent, ID)
}

//从全局内容列表中删除某一个内容
func delFromGlobalContent(ID string) error {
	return JsRedis.Redis_Sdel(constant.KEY_GlobalContent, ID)
}

//返回全局的内容id列表
func getGlobalContent() ([]string, error) {
	data := []string{}
	d, err := JsRedis.Redis_Sget(constant.KEY_GlobalContent)
	for _, v := range d {
		data = append(data, string(v.([]byte)))
	}
	return data, err
}
