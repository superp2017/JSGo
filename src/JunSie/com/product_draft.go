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
//  产品草稿的保存、删除、获取
//

//新建产品草稿
func NewProductDraft(s *JsHttp.Session) {
	st := &XM_Product{}
	if err := s.GetPara(st); err != nil {
		info := "newProductDraft:" + err.Error()
		JsLogger.Error(info)
		s.Forward("1", info, nil)
		return
	}
	if st.UID == "" {
		info := "newProductDraft,UID is empty\n"
		JsLogger.Error(info)
		s.Forward("1", info, nil)
		return
	}
	st.ID = util.IDer(constant.DB_Product)
	st.CreatDate = util.CurTime()
	if err := addProDraft(constant.DB_ProDraft, st.UID, st); err != nil {
		s.Forward("1", err.Error(), nil)
		return
	}
	s.Forward("0", "success", st)
}

//删除某一个草稿
func RemoveProductDraft(s *JsHttp.Session) {
	type draft struct {
		UID   string //用户id
		ProID string //产品id
	}
	st := &draft{}
	if err := s.GetPara(st); err != nil {
		info := "RemoveProductDraft:" + err.Error()
		JsLogger.Error(info)
		s.Forward("1", info, nil)
		return
	}
	if st.ProID == "" || st.UID == "" {
		info := fmt.Sprintf("removeProductDraft,param is empty,ProID=%s,UID=%s\n", st.ProID, st.UID)
		JsLogger.Error(info)
		s.Forward("1", info, nil)
		return
	}
	if st.UID == "" {
		st.UID = "Admin"
	}
	if err := removeProDraft(constant.DB_ProDraft, st.UID, st.ProID); err != nil {
		info := "removeProDraft:" + err.Error()
		JsLogger.Error(info)
		s.Forward("1", info, nil)
		return
	}
	s.Forward("0", "success", nil)
}

//获取草稿列表
func GetProductDraft(s *JsHttp.Session) {
	type draft struct {
		UID string //用户id
	}
	st := &draft{}
	if err := s.GetPara(st); err != nil {
		info := "getProductDraft:" + err.Error()
		JsLogger.Error(info)
		s.Forward("1", info, nil)
		return
	}
	if st.UID == "" {
		info := "getProductDraft,UID is empty\n"
		JsLogger.Error(info)
		s.Forward("1", info, nil)
		return
	}
	data, err := getProDraft(constant.DB_ProDraft, st.UID)
	if err != nil {
		info := "getProductDraft:" + err.Error()
		JsLogger.Error(info)
		s.Forward("1", info, nil)
		return
	}
	s.Forward("0", "success", data)

}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

///添加一个产品草稿到数据库
func addProDraft(db, key string, pro *XM_Product) error {
	data := []*XM_Product{}
	if err := JsRedis.Redis_hget(db, key, &data); err != nil {
		JsLogger.Error("addProDraft not exist Draft:" + err.Error())
	}
	data = append(data, pro)
	if len(data) > constant.ProDraft_SIZE {
		data = data[1:]
	}
	return JsRedis.Redis_hset(db, key, data)
}

//删除某一个草稿
func removeProDraft(db, key, id string) error {
	data := []*XM_Product{}
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

//获取草稿列表
func getProDraft(db, key string) ([]*XM_Product, error) {
	data := []*XM_Product{}
	return data, JsRedis.Redis_hget(db, key, &data)
}
