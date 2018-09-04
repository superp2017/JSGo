package com

import (
	"JsGo/JsBench/JsStatistics"
	"JsGo/JsHttp"
	"JsGo/JsLogger"
	"JsGo/JsStore/JsRedis"
	"JunSie/constant"
	"fmt"
)

//获取完整的内容统计
func GetContentStatics(session *JsHttp.Session) {
	type Para struct {
		ContentID string //内容id
	}
	st := &Para{}
	if err := session.GetPara(st); err != nil {
		JsLogger.Error(err.Error())
		session.Forward("1", err.Error(), nil)
		return
	}
	if st.ContentID == "" {
		JsLogger.Error("GetContentStatics failed,ContentID is empty\n")
		session.Forward("1", "GetContentStatics failed,ContentID is empty\n", nil)
		return
	}
	data, err := ContentStatics(st.ContentID)
	if err != nil {
		JsLogger.Error(err.Error())
		session.Forward("1", err.Error(), nil)
		return
	}
	session.Forward("0", "success", data)
}

func ContentStatics(ContentID string) (*JsStatistics.Statistics, error) {
	data := &JsStatistics.Statistics{}
	err := JsRedis.Redis_hget(constant.ContentStatistics, ContentID, data)
	return data, err
}

//增加内容访问量
func NewContentVisit(session *JsHttp.Session) {
	type Para struct {
		ContentID string //内容id
		UID       string //用户id
	}
	st := &Para{}
	if err := session.GetPara(st); err != nil {
		JsLogger.Error(err.Error())
		session.Forward("1", err.Error(), nil)
		return
	}
	if st.ContentID == "" || st.UID == "" {
		info := fmt.Sprintf("NewContentVisit failed,ContentID =%s,UID =%s\n", st.ContentID, st.UID)
		JsLogger.Error(info)
		session.Forward("2", info, nil)
		return
	}
	data, err := JsStatistics.NewVisit(constant.ContentStatistics, st.ContentID, "1")
	if err != nil {
		JsLogger.Error(err.Error())
		session.Forward("5", err.Error(), nil)
		return
	}
	//将浏览记录到个人信息中（+）
	if err := RecordPageView(true, st.UID, st.ContentID, 1); err != nil {
		JsLogger.Error(err.Error())
		session.Forward("3", err.Error(), nil)
		return
	}
	session.Forward("0", "success", data)
}

//增加内容点赞量
func NewContentPraise(session *JsHttp.Session) {
	type Para struct {
		ContentID string //内容id
		UID       string //用户id
	}
	st := &Para{}
	if err := session.GetPara(st); err != nil {
		JsLogger.Error(err.Error())
		session.Forward("1", err.Error(), nil)
		return
	}
	if st.ContentID == "" || st.UID == "" {
		info := fmt.Sprintf("NewContentPraise failed,ContentID =%s,UID =%s\n", st.ContentID, st.UID)
		JsLogger.Error(info)
		session.Forward("1", info, nil)
		return
	}
	data, err := JsStatistics.NewPraise(constant.ContentStatistics, st.ContentID, "1", st.UID)
	if err != nil {
		JsLogger.Error(err.Error())
		session.Forward("1", err.Error(), nil)
		return
	}
	//将点赞记录到个人信息中（+）
	if err := RecordLike(true, st.UID, st.ContentID, 1); err != nil {
		JsLogger.Error(err.Error())
		session.Forward("3", err.Error(), nil)
		return
	}
	session.Forward("0", "success", data)
}

//关注内容
func NewContentAttention(session *JsHttp.Session) {
	type Para struct {
		ContentID string //内容id
		UID       string //用户id
	}
	st := &Para{}
	if err := session.GetPara(st); err != nil {
		JsLogger.Error(err.Error())
		session.Forward("1", err.Error(), nil)
		return
	}
	if st.ContentID == "" || st.UID == "" {
		info := fmt.Sprintf("NewContentAttention failed,ContentID =%s,UID =%s\n", st.ContentID, st.UID)
		JsLogger.Error(info)
		session.Forward("1", info, nil)
		return
	}
	data, err := JsStatistics.NewAttention(constant.ContentStatistics, st.ContentID, "1", st.UID)
	if err != nil {
		JsLogger.Error(err.Error())
		session.Forward("1", err.Error(), nil)
		return
	}
	//收藏到个人记录里(添加)
	if err := RecordCollection(true, st.UID, st.ContentID, 1); err != nil {
		JsLogger.Error(err.Error())
		session.Forward("3", err.Error(), nil)
		return
	}
	session.Forward("0", "success", data)
}

//评论内容
func NewContentComment(ContentID string, score float64) error {
	if ContentID == "" {
		return JsLogger.ErrorLog("NewContentComment failed,ContentID is empty\n")
	}
	_, err := JsStatistics.NewComment(constant.ContentStatistics, ContentID, "1", score)
	if err != nil {
		return JsLogger.ErrorLog(err.Error())
	}
	return nil
}

//取消内容点赞
func RemoveContentPraise(session *JsHttp.Session) {
	type Para struct {
		ContentID string //内容id
		UID       string //用户id
	}
	st := &Para{}
	if err := session.GetPara(st); err != nil {
		JsLogger.Error(err.Error())
		session.Forward("1", err.Error(), nil)
		return
	}
	if st.ContentID == "" || st.UID == "" {
		info := fmt.Sprintf("RemoveContentPraise failed,ContentID =%s,UID =%s\n", st.ContentID, st.UID)
		JsLogger.Error(info)
		session.Forward("1", info, nil)
		return
	}
	data, err := JsStatistics.RemovePraise(constant.ContentStatistics, st.ContentID, "1", st.UID)
	if err != nil {
		JsLogger.Error(err.Error())
		session.Forward("1", err.Error(), nil)
		return
	}
	//将点赞记录到个人信息中（去除）
	if err := RecordLike(false, st.UID, st.ContentID, 1); err != nil {
		JsLogger.Error(err.Error())
		session.Forward("3", err.Error(), nil)
		return
	}
	session.Forward("0", "success", data)
}

//取消关注内容
func RemoveContentAttention(session *JsHttp.Session) {
	type Para struct {
		ContentID string //内容id
		UID       string //用户id
	}
	st := &Para{}
	if err := session.GetPara(st); err != nil {
		JsLogger.Error(err.Error())
		session.Forward("1", err.Error(), nil)
		return
	}
	if st.ContentID == "" || st.UID == "" {
		info := fmt.Sprintf("RemoveContentAttention failed,ContentID =%s,UID =%s\n", st.ContentID, st.UID)
		JsLogger.Error(info)
		session.Forward("1", info, nil)
		return
	}
	data, err := JsStatistics.RemoveAttention(constant.ContentStatistics, st.ContentID, "1", st.UID)
	if err != nil {
		JsLogger.Error(err.Error())
		session.Forward("1", err.Error(), nil)
		return
	}
	//收藏添加到个人记录里
	if err := RecordCollection(false, st.UID, st.ContentID, 1); err != nil {
		JsLogger.Error(err.Error())
		session.Forward("3", err.Error(), nil)
		return
	}
	session.Forward("0", "success", data)
}

//删除访问历史记录
func RemoveArtVisit(session *JsHttp.Session) {
	type Para struct {
		ContentID string //内容id
		UID       string //用户id
	}
	st := &Para{}
	if err := session.GetPara(st); err != nil {
		JsLogger.Error(err.Error())
		session.Forward("1", err.Error(), nil)
		return
	}
	if st.ContentID == "" || st.UID == "" {
		info := fmt.Sprintf("RemoveArtVisit failed,ProID =%s,UID =%s\n", st.ContentID, st.UID)
		JsLogger.Error(info)
		session.Forward("2", info, nil)
		return
	}

	if err := RecordPageView(false, st.UID, st.ContentID, 1); err != nil {
		JsLogger.Error(err.Error())
		session.Forward("5", err.Error(), nil)
		return
	}
	session.Forward("0", "success", nil)
}
