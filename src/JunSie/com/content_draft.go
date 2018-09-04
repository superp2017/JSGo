package com

import (
	"JsGo/JsHttp"
	"JsGo/JsLogger"
	"JsGo/JsStore/JsRedis"
	"JunSie/constant"
	"JunSie/util"
	"fmt"
)

//
//  内容草稿的保存、删除、获取
//

//保存内容草稿
func SaveContentDraft(s *JsHttp.Session) {
	st := &XM_Contents{}
	if err := s.GetPara(st); err != nil {
		info := "SaveContentDraft:" + err.Error()
		JsLogger.Error(info)
		s.Forward("1", info, nil)
		return
	}
	if st.UID == "" {
		info := "SaveContentDraft,UID is empty\n"
		JsLogger.Error(info)
		s.Forward("1", info, nil)
		return
	}
	st.ID = util.IDer(constant.DB_Product)
	st.CreatDate = util.CurTime()
	if err := addContentDraft(constant.DB_ContentDraft, st.UID, st); err != nil {
		s.Forward("1", err.Error(), nil)
		return
	}
	s.Forward("0", "success", st)
}

//查询内容草稿
func QueryContentDraft(s *JsHttp.Session) {
	type draft struct {
		UID string //用户id
	}
	st := &draft{}
	if err := s.GetPara(st); err != nil {
		info := "QueryContentDraft:" + err.Error()
		JsLogger.Error(info)
		s.Forward("1", info, nil)
		return
	}
	if st.UID == "" {
		info := "QueryContentDraft,UID is empty\n"
		JsLogger.Error(info)
		s.Forward("1", info, nil)
		return
	}
	data, err := getContentDraft(constant.DB_ContentDraft, st.UID)
	if err != nil {
		info := "QueryContentDraft:" + err.Error()
		JsLogger.Error(info)
		s.Forward("1", info, nil)
		return
	}
	s.Forward("0", "success", data)
}

//删除某一个内容草稿
func RemoveContentDraft(s *JsHttp.Session) {
	type draft struct {
		UID string //用户id
		CID string //内容id
	}
	st := &draft{}
	if err := s.GetPara(st); err != nil {
		info := "RemoveContentDraft:" + err.Error()
		JsLogger.Error(info)
		s.Forward("1", info, nil)
		return
	}
	if st.CID == "" || st.UID == "" {
		info := fmt.Sprintf("RemoveContentDraft,param is empty,ProID=%s,UID=%s\n", st.CID, st.UID)
		JsLogger.Error(info)
		s.Forward("1", info, nil)
		return
	}
	if st.UID == "" {
		st.UID = "Admin"
	}
	if err := removeContentDraft(constant.DB_ProDraft, st.UID, st.CID); err != nil {
		info := "RemoveContentDraft:" + err.Error()
		JsLogger.Error(info)
		s.Forward("1", info, nil)
		return
	}
	s.Forward("0", "success", nil)
}

//删除某一个内容草稿
func removeContentDraft(db, key, id string) error {
	data := []*XM_Contents{}
	if err := JsRedis.Redis_hget(db, key, &data); err != nil {
		return err
	}
	index := -1
	for i, v := range data {
		if v.ID == id {
			index = i
		}
	}
	if index != -1 {
		data = append(data[index:], data[:index+1]...)
	}
	return JsRedis.Redis_hset(db, key, data)
}

//增加一个内容草稿
func addContentDraft(db, key string, v *XM_Contents) error {
	data := []*XM_Contents{}
	if err := JsRedis.Redis_hget(db, key, &data); err != nil {
		JsLogger.Error("addContentDraft not exist Draft:" + err.Error())
	}

	data = append(data, v)
	if len(data) > constant.ContentDraft_SIZE {
		data = data[1:]
	}
	return JsRedis.Redis_hset(db, key, data)
}

//获取内容草稿列表
func getContentDraft(db, key string) ([]*XM_Contents, error) {
	data := []*XM_Contents{}
	return data, JsRedis.Redis_hget(db, key, &data)
}
