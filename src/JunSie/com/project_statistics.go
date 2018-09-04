package com

import (
	"JsGo/JsBench/JsStatistics"
	"JsGo/JsHttp"
	"JsGo/JsLogger"
	"JsGo/JsStore/JsRedis"
	"JunSie/constant"
	"fmt"
)

//获取完整的产品统计
func GetProductStatics(session *JsHttp.Session) {
	type Para struct {
		ProID string //产品id
	}
	st := &Para{}
	if err := session.GetPara(st); err != nil {
		JsLogger.Error(err.Error())
		session.Forward("1", err.Error(), nil)
		return
	}
	if st.ProID == "" {
		JsLogger.Error("GetProductStatics failed,ProID is empty\n")
		session.Forward("1", "GetProductStatics failed,ProID is empty\n", nil)
		return
	}

	data, err := ProStatics(st.ProID)
	if err != nil {
		JsLogger.Error(err.Error())
		session.Forward("1", err.Error(), nil)
		return
	}
	session.Forward("0", "success", data)
}

func ProStatics(ProID string) (*JsStatistics.Statistics, error) {
	data := &JsStatistics.Statistics{}
	err := JsRedis.Redis_hget(constant.ProStatistics, ProID, data)
	return data, err
}

//增加产品访问量
func NewProVisit(session *JsHttp.Session) {
	type Para struct {
		ProID string //产品id
		UID   string //用户id
	}
	st := &Para{}
	if err := session.GetPara(st); err != nil {
		JsLogger.Error(err.Error())
		session.Forward("1", err.Error(), nil)
		return
	}
	if st.ProID == "" || st.UID == "" {
		info := fmt.Sprintf("NewProVisit failed,ProID =%s,UID =%s\n", st.ProID, st.UID)
		JsLogger.Error(info)
		session.Forward("1", info, nil)
		return
	}
	data, err := JsStatistics.NewVisit(constant.ProStatistics, st.ProID, "0")
	if err != nil {
		JsLogger.Error(err.Error())
		session.Forward("1", err.Error(), nil)
		return
	}
	//将浏览记录到个人信息中（+）
	if err := RecordPageView(true, st.UID, st.ProID, 2); err != nil {
		JsLogger.Error(err.Error())
		session.Forward("3", err.Error(), nil)
		return
	}
	session.Forward("0", "success", data)
}

//增加产品点赞量
func NewProPraise(session *JsHttp.Session) {
	type Para struct {
		ProID string //产品id
		UID   string //用户id
	}
	st := &Para{}
	if err := session.GetPara(st); err != nil {
		JsLogger.Error(err.Error())
		session.Forward("1", err.Error(), nil)
		return
	}
	if st.ProID == "" || st.UID == "" {
		info := fmt.Sprintf("NewProPraise failed,ProID =%s,UID =%s\n", st.ProID, st.UID)
		JsLogger.Error(info)
		session.Forward("1", info, nil)
		return
	}
	data, err := JsStatistics.NewPraise(constant.ProStatistics, st.ProID, "0", st.UID)
	if err != nil {
		JsLogger.Error(err.Error())
		session.Forward("1", err.Error(), nil)
		return
	}
	//将点赞记录到个人信息中1111111111111111111111
	if err := RecordLike(true, st.UID, st.ProID, 2); err != nil {
		JsLogger.Error(err.Error())
		session.Forward("3", err.Error(), nil)
		return
	}
	session.Forward("0", "success", data)
}

//关注产品
func NewProAttention(session *JsHttp.Session) {
	type Para struct {
		ProID string //产品id
		UID   string //用户id
	}
	st := &Para{}
	if err := session.GetPara(st); err != nil {
		JsLogger.Error(err.Error())
		session.Forward("1", err.Error(), nil)
		return
	}
	if st.ProID == "" || st.UID == "" {
		info := fmt.Sprintf("NewProAttention failed,ProID =%s,UID =%s\n", st.ProID, st.UID)
		JsLogger.Error(info)
		session.Forward("1", info, nil)
		return
	}
	data, err := JsStatistics.NewAttention(constant.ProStatistics, st.ProID, "0", st.UID)
	if err != nil {
		JsLogger.Error(err.Error())
		session.Forward("1", err.Error(), nil)
		return
	}
	if err := RecordCollection(true, st.UID, st.ProID, 2); err != nil {
		JsLogger.Error(err.Error())
		session.Forward("3", err.Error(), nil)
		return
	}
	session.Forward("0", "success", data)
}

//评论产品
func NewProComment(ProID string, score float64) error {
	if ProID == "" {
		return JsLogger.ErrorLog("NewProComment failed,ProID is empty\n")
	}
	_, err := JsStatistics.NewComment(constant.ProStatistics, ProID, "0", score)
	if err != nil {
		return JsLogger.ErrorLog(err.Error())
	}
	return nil
}

//取消产品点赞
func RemoveProPraise(session *JsHttp.Session) {
	type Para struct {
		ProID string //产品id
		UID   string //用户id
	}
	st := &Para{}
	if err := session.GetPara(st); err != nil {
		JsLogger.Error(err.Error())
		session.Forward("1", err.Error(), nil)
		return
	}
	if st.ProID == "" || st.UID == "" {
		info := fmt.Sprintf("RemoveProPraise failed,ProID =%s,UID =%s\n", st.ProID, st.UID)
		JsLogger.Error(info)
		session.Forward("1", info, nil)
		return
	}
	data, err := JsStatistics.RemovePraise(constant.ProStatistics, st.ProID, "0", st.UID)
	if err != nil {
		JsLogger.Error(err.Error())
		session.Forward("1", err.Error(), nil)
		return
	}
	if err := RecordLike(false, st.UID, st.ProID, 2); err != nil {
		JsLogger.Error(err.Error())
		session.Forward("4", err.Error(), nil)
		return
	}
	session.Forward("0", "success", data)
}

//取消关注产品
func RemoveProAttention(session *JsHttp.Session) {
	type Para struct {
		ProID string //产品id
		UID   string //用户id
	}
	st := &Para{}
	if err := session.GetPara(st); err != nil {
		JsLogger.Error(err.Error())
		session.Forward("1", err.Error(), nil)
		return
	}
	if st.ProID == "" || st.UID == "" {
		info := fmt.Sprintf("RemoveProAttention failed,ProID =%s,UID =%s\n", st.ProID, st.UID)
		JsLogger.Error(info)
		session.Forward("1", info, nil)
		return
	}
	data, err := JsStatistics.RemoveAttention(constant.ProStatistics, st.ProID, "0", st.UID)
	if err != nil {
		JsLogger.Error(err.Error())
		session.Forward("1", err.Error(), nil)
		return
	}
	if err := RecordCollection(false, st.UID, st.ProID, 2); err != nil {
		JsLogger.Error(err.Error())
		session.Forward("3", err.Error(), nil)
		return
	}
	session.Forward("0", "success", data)
}

//删除访问历史记录
func RemoveProVisit(session *JsHttp.Session) {
	type Para struct {
		ProID string //产品id
		UID   string //用户id
	}
	st := &Para{}
	if err := session.GetPara(st); err != nil {
		JsLogger.Error(err.Error())
		session.Forward("1", err.Error(), nil)
		return
	}
	if st.ProID == "" || st.UID == "" {
		info := fmt.Sprintf("RemoveProVisit failed,ProID =%s,UID =%s\n", st.ProID, st.UID)
		JsLogger.Error(info)
		session.Forward("2", info, nil)
		return
	}

	if err := RecordPageView(false, st.UID, st.ProID, 2); err != nil {
		JsLogger.Error(err.Error())
		session.Forward("3", err.Error(), nil)
		return
	}
	session.Forward("0", "success", nil)
}

//新产品的销量统计
func NewProductSaleNum(ProID string, nums int) {
	_, e := JsStatistics.NewSales(constant.ProStatistics, ProID, "0", nums)
	if e == nil {

		go UpdateHotList(ProID)
	}
}
