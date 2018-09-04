package com

import (
	"JsGo/JsBench/JsUser"
	"JsGo/JsHttp"
	"JsGo/JsLogger"
	"JsGo/JsMobile"
	. "JsGo/JsStore/JsRedis"
	"JunSie/constant"
	"JunSie/util"
	"fmt"

	"github.com/chanxuehong/wechat/mp/user/oauth2"
)

func InitUser() {
	JsHttp.WhiteHttps("/exchangeuser", ExchangeUser) //从UID获取用户信息
	JsHttp.WhiteHttp("/moduser", ModUser)
	JsHttp.WhiteHttp("/getuserinfofromopenid", GetUserInfoFromOpenID)   //从openID获取用户信息
	JsHttp.WhiteHttp("/getuserinfofromunionid", GetUserInfoFromUnionID) //从unionID获取用户信息
	JsHttp.WhiteHttp("/getuserfrommobile", GetUserFromMobile)           //从Mobile获取用户信息
	JsHttp.WhiteHttps("/getuserinfoabs", Getuserinfoabs)                //获取用多个户简要信息
}
func InitUserMall() {
	JsHttp.WhiteHttps("/exchangeuser", ExchangeUser) //从UID获取用户信息
	JsHttp.WhiteHttps("/moduser", ModUser)
	JsHttp.WhiteHttps("/getuserinfofromopenid", GetUserInfoFromOpenID)     //从openID获取用户信息
	JsHttp.WhiteHttps("/getuserinfofromunionid", GetUserInfoFromUnionID)   //从unionID获取用户信息
	JsHttp.WhiteHttps("/mobilelogin", NewUserFromMobile)                   //手机号注册用户
	JsHttp.WhiteHttps("/getuserfrommobile", GetUserFromMobile)             //从Mobile获取用户信息
	JsHttp.WhiteHttps("/appenduserreceivingaddr", AppendUserReceivingAddr) //追加一个用户收货地址
	JsHttp.WhiteHttps("/updateuserrecevingaddr", UpdateUserRecevingAddr)   //更新用户收货地址
	JsHttp.WhiteHttps("/getuserinfoabs", Getuserinfoabs)                   //获取用多个户简要信息
}

func initNewUser(user *oauth2.UserInfo) *JsUser.User {
	newuser := &JsUser.User{}
	newuser.Nickname = user.Nickname
	newuser.Sex = user.Sex
	newuser.City = user.City
	newuser.Header = user.HeadImageURL
	newuser.Country = user.Country
	newuser.Unionid = user.UnionId
	newuser.CreatTime = util.CurTime()
	return newuser
}

func WxNewUser(user *oauth2.UserInfo, session *JsHttp.Session, userChannel string) string {

	if len(user.UnionId) > 0 {
		UID := ""
		Redis_hget(constant.UNIONUID, user.UnionId, &UID)
		newuser := &JsUser.User{}
		if UID != "" {
			go updateUserOpenID(UID, user.OpenId, userChannel)
			return UID
		}
		newuser = initNewUser(user)
		fitOpenID(newuser, user.OpenId, userChannel)
		if err := newuser.NewUser(constant.USER); err == nil {
			go Redis_hset(constant.UNIONUID, user.UnionId, newuser.ID)
			go Redis_hset(constant.OPENUID, user.OpenId, newuser.ID)
		}
		return newuser.ID
	} else if len(user.OpenId) > 0 {
		UID := ""
		Redis_hget(constant.OPENUID, user.OpenId, &UID)
		if UID != "" {
			go updateUserOpenID(UID, user.OpenId, userChannel)
			return UID
		}
		newuser := initNewUser(user)
		fitOpenID(newuser, user.OpenId, userChannel)
		if err := newuser.NewUser(constant.USER); err == nil {
			go Redis_hset(constant.UNIONUID, user.UnionId, newuser.ID)
			go Redis_hset(constant.OPENUID, user.OpenId, newuser.ID)
		}
		return newuser.ID
	}
	return ""
}

//更新用户openid
func updateUserOpenID(UID, opendID, userChannel string) {
	newuser := &JsUser.User{}
	if err := Redis_hget(constant.USER, UID, newuser); err == nil {
		fitOpenID(newuser, opendID, userChannel)
		go Redis_hset(constant.USER, UID, newuser)
	}
}

//填充用户openID
func fitOpenID(user *JsUser.User, opendID, userChannel string) {
	if userChannel == "web" {
		user.Openid = opendID
	}
	if userChannel == "app" {
		user.Openid_app = opendID
	}
	if userChannel == "minip" {
		user.Openid_small = opendID
	}
}

//从openID获取用户信息
func GetUserInfoFromOpenID(s *JsHttp.Session) {
	type para struct {
		OpendID string
	}
	st := &para{}
	if err := s.GetPara(st); err != nil {
		s.Forward("1", err.Error(), nil)
		return
	}
	var UID string
	if err := Redis_hget(constant.OPENUID, st.OpendID, &UID); err != nil {
		s.Forward("1", err.Error(), nil)
		return
	}
	data := &JsUser.User{}
	if err := Redis_hget(constant.USER, UID, data); err != nil {
		s.Forward("1", err.Error(), nil)
		return
	}
	s.Forward("0", "success", data)
}

//从unionID获取用户信息
func GetUserInfoFromUnionID(s *JsHttp.Session) {
	type para struct {
		UnionID string
	}
	st := &para{}
	if err := s.GetPara(st); err != nil {
		s.Forward("1", err.Error(), nil)
		return
	}
	var UID string
	if err := Redis_hget(constant.UNIONUID, st.UnionID, &UID); err != nil {
		s.Forward("1", err.Error(), nil)
		return
	}
	data := &JsUser.User{}
	if err := Redis_hget(constant.USER, UID, data); err != nil {
		s.Forward("1", err.Error(), nil)
		return
	}
	s.Forward("0", "success", data)
}

func ExchangeUser(s *JsHttp.Session) {
	type Para struct {
		UID string
	}
	para := &Para{}
	e := s.GetPara(para)
	go NewVisit() //增加用户访问量
	if e != nil {
		JsLogger.Error(e.Error())
		s.Forward("1", e.Error(), "")
		return
	}

	user := &JsUser.User{}
	e = Redis_hget(constant.USER, para.UID, user)
	if e != nil {
		JsLogger.Error(e.Error())
		s.Forward("2", e.Error(), "")
		return
	}

	s.MarkSession()
	s.Forward("0", "success", user)
}

func ModUser(s *JsHttp.Session) {
	user := &JsUser.User{}
	e := s.GetPara(user)
	if e != nil {
		JsLogger.Error(e.Error())
		s.Forward("1", e.Error(), "")
		return
	}

	e = user.ModUser(constant.USER)
	if e != nil {
		JsLogger.Error(e.Error())
		s.Forward("2", e.Error(), "")
		return
	}
	s.Forward("0", "success", user)
}

func NewUserFromMobile(s *JsHttp.Session) {
	type Para struct {
		Mobile string
		Code   string
	}
	st := &Para{}
	if err := s.GetPara(st); err != nil {
		s.Forward("1", err.Error(), nil)
		return
	}
	if st.Mobile == "" || st.Code == "" {
		info := fmt.Sprintf("NewUserFromMobile failed,Cell = %s,SmsCode =%s \n", st.Mobile, st.Code)
		JsLogger.Error(info)
		s.Forward("1", info, nil)
		return
	}
	if ok := JsMobile.VerifySmsCode(st.Mobile, st.Code); !ok {
		info := fmt.Sprintf("NewUserFromMobile VerifySmsCode failed,Cell = %s,SmsCode =%s \n", st.Mobile, st.Code)
		JsLogger.Error(info)
		s.Forward("1", info, nil)
		return
	}

	uid := ""
	if e2 := Redis_hget(constant.Key_Mobile, st.Mobile, &uid); e2 == nil {
		u := &JsUser.User{}
		if e3 := Redis_hget(constant.USER, uid, &u); e3 == nil {
			s.Forward("0", "success", u)
			return
		}
	}

	user := &JsUser.User{}
	user.ID = util.IDer(constant.USER)
	user.Mobile = st.Mobile
	user.CreatTime = util.CurTime()
	if err := Redis_hset(constant.USER, user.ID, user); err != nil {
		JsLogger.Error(err.Error())
		s.Forward("1", err.Error(), nil)
		return
	}
	go Redis_hset(constant.Key_Mobile, user.Mobile, user.ID)
	s.Forward("0", "success", user)
}

func GetUserFromMobile(s *JsHttp.Session) {
	type Para struct {
		Cell string
	}
	st := &Para{}
	if err := s.GetPara(st); err != nil {
		JsLogger.Error(err.Error())
		s.Forward("1", err.Error(), nil)
		return
	}
	if st.Cell == "" {
		JsLogger.Error("GetUserFromMobile failed,Cell is empty \n")
		s.Forward("1", "GetUserFromMobile failed,Cell is empty \n", nil)
		return
	}
	var uid string
	if e := Redis_hget(constant.Key_Mobile, st.Cell, &uid); e != nil {
		JsLogger.Error(e.Error())
		s.Forward("1", e.Error(), nil)
		return
	}
	user := &JsUser.User{}
	if e := Redis_hget(constant.USER, uid, user); e != nil {
		JsLogger.Error(e.Error())
		s.Forward("1", e.Error(), nil)
		return
	}
	s.Forward("0", "success", user)
}

//添加一个收货地址
func AppendUserReceivingAddr(session *JsHttp.Session) {
	type Para struct {
		UID      string //用户id
		Alias    string //收货信息别名
		Name     string //收货人姓名
		Cell     string //收货人电话
		Province string //省
		City     string //市
		Area     string //区
		Addr     string //收货地址
	}
	st := &Para{}
	if err := session.GetPara(st); err != nil {
		JsLogger.Error(err.Error())
		session.Forward("1", err.Error(), nil)
		return
	}

	if st.Addr == "" || st.Cell == "" || st.Name == "" {
		str := fmt.Sprintf("AppendReceivingAddr failed,Addr=%s,Cell=%s,Name=%s\n", st.Addr, st.Cell, st.Name)
		JsLogger.Error(str)
		session.Forward("1", str, nil)
		return
	}
	user := &JsUser.User{}
	if e := Redis_hget(constant.USER, st.UID, user); e != nil {
		JsLogger.Error(e.Error())
		session.Forward("1", e.Error(), nil)
		return
	}
	user.ReceivingAddrs = append(user.ReceivingAddrs, JsUser.RecAddr{
		Alias:    st.Alias,
		Name:     st.Name,
		Cell:     st.Cell,
		Addr:     st.Addr,
		Province: st.Province,
		City:     st.City,
		Area:     st.Area,
	})
	if e := Redis_hset(constant.USER, st.UID, user); e != nil {
		JsLogger.Error(e.Error())
		session.Forward("1", e.Error(), nil)
		return
	}
	session.Forward("0", "AppendReceivingAddr success\n", user)
}

//更新收货地址
func UpdateUserRecevingAddr(session *JsHttp.Session) {
	type Para struct {
		UID            string //用户ID
		ReceivingAddrs []JsUser.RecAddr
	}
	st := &Para{}
	if err := session.GetPara(st); err != nil {
		JsLogger.Error(err.Error())
		session.Forward("1", "UpdateUserRecevingAddr failed\n", nil)
		return
	}
	if st.UID == "" {
		JsLogger.Error("UpdateUserRecevingAddr failed,UID is empty\n")
		session.Forward("1", "UpdateUserRecevingAddr failed,UID is empty\n", nil)
		return
	}
	user := &JsUser.User{}
	if e := Redis_hget(constant.USER, st.UID, user); e != nil {
		JsLogger.Error(e.Error())
		session.Forward("1", e.Error(), nil)
		return
	}
	user.ReceivingAddrs = st.ReceivingAddrs
	if e := Redis_hset(constant.USER, st.UID, user); e != nil {
		JsLogger.Error(e.Error())
		session.Forward("1", e.Error(), nil)
		return
	}
	session.Forward("0", "UpdateUserRecevingAddr success\n", user)
}

//获取用多个户简要信息JsHttp.WhiteHttp("/getuserinfoabs", Getuserinfoabs)
func Getuserinfoabs(s *JsHttp.Session) {
	type Para struct {
		UIDS []string
	}
	para := &Para{}
	e := s.GetPara(para)
	if e != nil {
		JsLogger.Error(e.Error())
		s.Forward("1", e.Error(), "")
		return
	}
	if len(para.UIDS) == 0 {
		info := fmt.Sprintf("Getuserinfoabs UID is empty\n uidslen,UIDS==%d", len(para.UIDS))
		JsLogger.Info(info)
		s.Forward("2", info, nil)
		return
	}
	s.Forward("0", "success", GetuserinfoABS(para.UIDS))
}
func GetuserinfoABS(ids []string) []*JsUser.UserABS {
	Users := []*JsUser.UserABS{} //用户摘要信息切片返回值
	for _, v := range ids {
		user := &JsUser.UserABS{}
		e := Redis_hget(constant.USER, v, user)
		if e != nil {
			Users = append(Users, nil)
			JsLogger.Error(e.Error())
			continue
		} else {
			Users = append(Users, user)
		}
	}
	return Users
}
