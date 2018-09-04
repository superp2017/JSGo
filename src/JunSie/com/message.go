package com

import (
	"JsGo/JsHttp"
	"JsGo/JsLogger"
	"JsGo/JsStore/JsRedis"
	"JunSie/constant"
	"JunSie/util"
	"fmt"
)

//消息推送
//推送草稿暂时存放。

//消息放一个切片数组中(放最近的三到五条)需要一个页面展示的

//历史消息（放需要保存的历史消息）其他的都放在历史

//添加推送消息

//用户推送的消息

type Message struct {
	MessTitle   string    //推送内容标题(不能为空)
	WordPicture []WordPic //图片文字
	MessPro     string    //推送关联产品
	MessArt     string    //推送关联文章
	MessTime    string    //推送添加时间
	TimeStamp   int64     //时间（用于判断是否查看消息）
}
type WordPic struct {
	Type     string
	MessWord string //推送内容
	MessPic  string //推送图片

}

func Init_Message() {
	JsHttp.WhiteHttp("/messageready", MessageReady)       //存为草稿
	JsHttp.WhiteHttp("/getmessageready", GetMessageReady) //获取推送草稿
	JsHttp.WhiteHttp("/messageadd", MessageAdd)           //添加推送消息
	JsHttp.WhiteHttp("/getmessage", GetMessage)           //获取推送消息
	JsHttp.WhiteHttp("/messageDel", MessageDel)           //删除推送消息
	JsHttp.WhiteHttp("/messagesend", MessageSend)         //商家直接推送消息
	JsHttp.WhiteHttp("/chatuseradd", ChatUserAdd)         //1用户留言
	JsHttp.WhiteHttp("/chatuserget", ChatUserGet)         //2用户获取留言//商家获取用户留言(单个用户留言)
	JsHttp.WhiteHttp("/chatbusreply", ChatBusReply)       //3商家回复留言（）
	JsHttp.WhiteHttp("/chantnumuser", Chantnumuser)       //4用户获取商家回复的未读留言条数（）//Chantnumuser
	JsHttp.WhiteHttp("/chatbusgetlist", ChatGetList)      //5商家获取用户留言（分页列表）(获取未读条数)
}

func Init_MessageMall() {
	JsHttp.WhiteHttps("/getmessage", GetMessage)     //获取推送消息
	JsHttp.WhiteHttps("/chatuseradd", ChatUserAdd)   //1用户留言
	JsHttp.WhiteHttps("/chatuserget", ChatUserGet)   //2用户获取留言//商家获取用户留言(单个用户留言)
	JsHttp.WhiteHttps("/chantnumuser", Chantnumuser) //4用户获取商家回复的未读留言条数（）//Chantnumuser
}

//首先创建推送草稿，
//草稿可以获取后再进行修改
//将草稿添加到推送消息，正式推送消息到用户
//用户和商家都可以获取到推送消息
//用户可以通过在本地存储消息推送时间，从而判断用户是否看过，
//也可以通过在用户的结构中增加一个时间字段，存放最新查看时间。

//存为草稿（草稿内容不检测）
func MessageReady(s *JsHttp.Session) {
	st := &Message{}
	err := s.GetPara(st)
	if err != nil {
		JsLogger.Error(err.Error())
		s.Forward("1", err.Error(), nil)
		return
	}
	//草稿本来就保存不完整数据属于不检查完整性
	err = JsRedis.Redis_set(constant.MessageReadyBus_001, st)
	if err != nil {
		JsLogger.Error(err.Error())
		s.Forward("3", err.Error(), nil)
		return
	}
	s.Forward("0", "success", nil)
}

//获取推送草稿
func GetMessageReady(s *JsHttp.Session) {
	st := &Message{}
	err := JsRedis.Redis_get(constant.MessageReadyBus_001, st)
	if err != nil {
		JsLogger.Error(err.Error())
		s.Forward("3", err.Error(), nil)
		return
	}
	s.Forward("0", "success", st)
}

//添加推送消息(将草稿推送）
func MessageAdd(s *JsHttp.Session) {
	data := []Message{}
	st := Message{}
	err := JsRedis.Redis_get(constant.MessageReadyBus_001, &st) //获取草稿
	if err != nil {
		JsLogger.Error(err.Error())
		s.Forward("3", err.Error(), nil)
		return
	}
	if st.MessTitle == "" { //检查逻辑正确性
		info := fmt.Sprintf("st.MessTitle=%s", st.MessTitle)
		JsLogger.Error(info)
		s.Forward("2", info, nil)
		return
	}
	st.MessTime = util.CurTime() //推送时自动添加时间
	st.TimeStamp = util.CurStamp()
	err = JsRedis.Redis_get(constant.MessageBus_001, &data) //获取消息数组
	if err != nil {
		JsLogger.Error(err.Error())
	}
	temp := append([]Message{}, st)
	data = append(temp, data...)
	err = JsRedis.Redis_set(constant.MessageBus_001, &data)
	if err != nil {
		JsLogger.Error(err.Error())
		s.Forward("3", err.Error(), nil)
		return
	}
	s.Forward("0", "success", nil)
}

//获取推送消息
func GetMessage(s *JsHttp.Session) {
	Sto := make([]Message, 0)
	err := JsRedis.Redis_get(constant.MessageBus_001, &Sto)
	if err != nil {
		JsLogger.Error(err.Error())
		s.Forward("3", err.Error(), Sto)
		return
	}
	s.Forward("0", "success", Sto)
}

//删除推送消息
func MessageDel(s *JsHttp.Session) {
	//发送要删除的标题
	type title struct {
		MessTitle string
	}
	para := &title{}
	err := s.GetPara(para)
	if err != nil {
		JsLogger.Error(err.Error())
		s.Forward("1", err.Error(), "")
		return
	}
	Sto := make([]Message, 0)
	err = JsRedis.Redis_get(constant.MessageBus_001, &Sto)
	if err != nil {
		JsLogger.Error(err.Error())
		s.Forward("3", err.Error(), nil)
		return
	}
	index := -1
	for i, V := range Sto {
		if V.MessTitle == para.MessTitle {
			index = i
		}
	}
	if index == -1 {
		s.Forward("5", "unfind", para)
		return
	} else {
		Sto = append(Sto[:index], Sto[index+1:]...)
	}
	err = JsRedis.Redis_set(constant.MessageBus_001, &Sto)
	if err != nil {
		JsLogger.Error(err.Error())
		s.Forward("4", err.Error(), &Sto)
		return
	}
	s.Forward("0", "success", &Sto)
}

//直接发送消息接口

func MessageSend(s *JsHttp.Session) {
	st := Message{}
	data := []Message{}
	err := s.GetPara(&st)
	if err != nil {
		JsLogger.Error(err.Error())
		s.Forward("1", err.Error(), nil)
		return
	}
	if st.MessTitle == "" { //检查逻辑正确性
		info := fmt.Sprintf("st.MessTitle=%s", st.MessTitle)
		JsLogger.Error(info)
		s.Forward("2", info, nil)
		return
	}
	st.MessTime = util.CurTime() //推送时自动添加时间
	st.TimeStamp = util.CurStamp()
	err = JsRedis.Redis_get(constant.MessageBus_001, &data) //获取消息数组
	if err != nil {
		JsLogger.Error(err.Error())
	}
	temp := append([]Message{}, st)
	data = append(temp, data...)
	err = JsRedis.Redis_set(constant.MessageBus_001, &data)
	if err != nil {
		JsLogger.Error(err.Error())
		s.Forward("3", err.Error(), nil)
		return
	}
	s.Forward("0", "success", &st)
}

//删除推送消息
//func DeleteMessage(s *JsHttp.Session){
//
//}
//商家ID
//获取历史推送消息
//func GetMessageH(s JsHttp.Session) {
//	//商家ID
//
//}

//客服
//用户消息留言代码逻辑
//一个用户id,为一个K值存储消息，商家和用户都对这个消息内容进行添加
//商家和用户使用的 消息结构相同
//整个消息是以时间顺序存储的一个切片
//用户和商家都往这个切片中追加消息结构（回复留言）
//UID为空表示商家，根据UID是否为空把消息内容放在 对话 的左右两侧

//添加商家和用户查看信息时间，用于处理商家是否有新的用户消息
//如果商家标记已经读取消息（在聊天列表中删除客户）//标记为已经读取消息。
//那么把标记为已读不在消息列表中获取到 （）
//如果判断有则把消息给商家。
//p排序最新的消息放在最前
//标记已经读取的是给消息的最新的一条加一个时间int64  （LookTime  int64  //	商家或用户查看时间），
// 如果检测到则表示已读取消息
//用户获取消息，，，，用户或取新消息（用户在消息页面等待）
//商家
//大于30天自动删除（没有）
//用户留言结构222222222222222222222222222222222222222222222222222222222222222222222222222222222222222222222222222222222222222222
type ChatMessage struct {
	UID      string //用户,为空时表示商家。
	UserName string //用户名 this.user.username
	Type     string //消息类型       Image图片  Text//文字   PID /关联的产品
	Text     string //消息内容'OH CRAP!!'    Image图片   ""//为空表示文字   PID /关联的产品
}

type Chat struct { //聊天最基本
	UID       string //用户,为空时表示商家。
	Type      string //消息类型   0 Text//文字   1 Image图片    PID /关联的产品
	Text      string //消息内容
	Date      string //new Date().toString(),//日期时间字符串
	TimeStamp int64  //时间戳
}
type ChatDB struct {
	ID       string
	UID      string //用户,为空时表示商家。
	UserName string //用户名 this.user.username
	//Header   string //用户头像 this.user.pic
	//TimeStamp int64  //时间戳(用户查看时间)
	//LookTime  int64  //	商家或用户查看时间
	Numbus  int //商家未读数量
	Numuser int //用户未读数量
	Chatmor []Chat
}
type Chatnum struct {
	Numbus int //商家未读数量
}

//1用户留言
func ChatUserAdd(s *JsHttp.Session) {
	para := ChatMessage{} //参数字段
	err := s.GetPara(&para)
	if err != nil {
		JsLogger.Error(err.Error())
		s.Forward("1", err.Error(), nil)
		return
	}
	if para.UID == "" {
		info := fmt.Sprintf("data.UID=%s", para.UID)
		JsLogger.Error(info)
		s.Forward("2", info, nil)
		return
	}
	sdata := ChatDB{} //数据库中保存的聊天结构
	err = JsRedis.Redis_hdbget(0, constant.Chat, para.UID, &sdata)
	if err != nil {
		JsLogger.Error(err.Error())
	}
	sdata.UID = para.UID
	sdata.UserName = para.UserName
	sdata.Numbus++ //商家未读消息加一
	data := Chat{} //最基本结构
	data.UID = "user"
	data.TimeStamp = util.CurStamp()
	data.Date = util.CurTimef()
	data.Type = para.Type
	data.Text = para.Text
	sdata.Chatmor = append(sdata.Chatmor, data)
	err = JsRedis.Redis_hdbset(0, constant.Chat, para.UID, &sdata)
	if err != nil {
		JsLogger.Error(err.Error())
		s.Forward("3", err.Error(), nil)
		return
	}

	//给商家消息队列进行添加和排序（将新消息排到最前面）
	go ChatBusaddlist(para.UID)
	s.Forward("0", "success", data)
}

//给商家消息队列进行添加和排序（将新消息排到最前面）
func ChatBusaddlist(UID string) {
	sdata := []string{}
	JsRedis.Redis_get(constant.ChatList, &sdata)
	var i int
	exist := false
	for n, v := range sdata {
		if v == UID {
			exist = true
			i = n
		}
	}

	if exist {
		if i == 0 { //如果是第一个直接返回
			return
		}
		temp := make([]string, 0, len(sdata)+1)
		temp = append([]string{}, UID)
		sdata = append(temp, sdata...)
		sdata = append(sdata[:i+1], sdata[i+2:]...)
	} else {
		temp := append([]string{}, UID)
		sdata = append(temp, sdata...)
	}

	JsRedis.Redis_set(constant.ChatList, &sdata)
	return
}

//4.用户获取商家回复未读条数
func Chantnumuser(s *JsHttp.Session) {
	type Para struct {
		UID string
	}
	para := &Para{}
	err := s.GetPara(para)
	if err != nil {
		JsLogger.Error(err.Error())
		s.Forward("1", err.Error(), nil)
		return
	}
	if para.UID == "" {
		s.Forward("2", "uid= nil ", nil)
		return
	}
	type Number struct {
		Numuser int //用户未读数量
	}
	data := &Number{}
	err = JsRedis.Redis_hdbget(0, constant.Chat, para.UID, data)
	if err != nil {
		JsLogger.Error(err.Error())
	}
	s.Forward("0", "success", data)
}

//2用户获取留言
// 用户和商家获取留言后自己的未读消息被清空为零前端拿到的也是零。
func ChatUserGet(s *JsHttp.Session) {
	type Para struct {
		Who string
		UID string
	}
	para := &Para{}
	err := s.GetPara(para)
	if err != nil {
		JsLogger.Error(err.Error())
		s.Forward("1", err.Error(), nil)
		return
	}
	if para.UID == "" {
		s.Forward("2", "uid= nil ", nil)
		return
	}
	sdata := ChatDB{} //返回结构
	err = JsRedis.Redis_hdbget(0, constant.Chat, para.UID, &sdata)
	if err != nil {
		JsLogger.Error(err.Error())
		s.Forward("3", err.Error(), nil)
		return
	}
	if para.Who == "" { //"0"表示商家   ""表示用户
		sdata.Numuser = 0
		JsRedis.Redis_hdbset(0, constant.Chat, para.UID, &sdata)
	} else {
		sdata.Numbus = 0
		JsRedis.Redis_hdbset(0, constant.Chat, para.UID, &sdata)

	}
	s.Forward("0", "success", sdata)

}

//商家回复留言（）
func ChatBusReply(s *JsHttp.Session) {
	para := Chat{} //参数字段//(UID 回复用户的UID)
	err := s.GetPara(&para)
	if err != nil {
		JsLogger.Error(err.Error())
		s.Forward("1", err.Error(), nil)
		return
	}
	if para.UID == "" {
		info := fmt.Sprintf("data.UID=%s", para.UID)
		JsLogger.Error(info)
		s.Forward("2", "uid = nil", nil)
	}
	sdata := ChatDB{}
	err = JsRedis.Redis_hdbget(0, constant.Chat, para.UID, &sdata)
	if err != nil {
		JsLogger.Error(err.Error())
	}
	sdata.Numuser++ //用户未读消息加一
	data := Chat{}  //最基本结构
	data.UID = ""   //空表示商家回复的消息
	data.TimeStamp = util.CurStamp()
	data.Date = util.CurTimef()
	data.Type = para.Type
	data.Text = para.Text
	sdata.Chatmor = append(sdata.Chatmor, data)
	err = JsRedis.Redis_hdbset(0, constant.Chat, para.UID, &sdata)
	if err != nil {
		JsLogger.Error(err.Error())
		s.Forward("3", err.Error(), nil)
		return
	}
	s.Forward("0", "success", data)
}

//商家获取客户消息列表(获取未读条数)
func ChatGetList(s *JsHttp.Session) {
	type Info struct {
		SIndex int //启始索引（从列表数组的第几个开始）
		Size   int //个数
	}
	st := &Info{}
	if err := s.GetPara(st); err != nil {
		info := "ChatGetList:" + err.Error()
		JsLogger.Error(info)
		s.Forward("1", info, nil)
		return
	}
	if st.SIndex < 0 || st.Size <= 0 {
		info := fmt.Sprintf("getPageProducts param error,SIndex=%d,Size=%d\n", st.SIndex, st.Size)
		JsLogger.Error(info)
		s.Forward("1", info, nil)
		return
	}
	idlist := []string{}
	err := JsRedis.Redis_get(constant.ChatList, &idlist)
	if err != nil {
		JsLogger.Error(err.Error())
		s.Forward("3", err.Error(), idlist)
		return
	}
	lenList := len(idlist)
	if st.SIndex > lenList {
		//超过长度直接返回空
		s.Forward("0", "success", nil)
		return
	} else {
		if (st.SIndex + st.Size) < lenList { //所取的在范围内
			idlist = idlist[st.SIndex : st.SIndex+st.Size]
		} else { //不在范围内
			idlist = idlist[st.SIndex:lenList]
		}
	}

	Mess := make(map[string]int)

	for _, v := range idlist {
		chat := ChatDB{}
		if err := JsRedis.Redis_hdbget(0, constant.Chat, v, &chat); err == nil {
			Mess[v] = chat.Numbus
		} else {
			JsLogger.Error(err.Error())
		}
	}
	type Para struct {
		Msg map[string]int
		Ids []string
	}
	data := &Para{
		Msg: Mess,
		Ids: idlist,
	}
	s.Forward("0", "success", data)
	//Mess := make(map[int][]ChatMessage)
	//for _, v := range idlist {
	//	i:=0 ;i++
	//	chat := []ChatMessage{}
	//	err := JsRedis.Redis_hdbget(0, constant.Chat, v, &chat)
	//	if err != nil {
	//		JsLogger.Error(err.Error())
	//	} else {
	//		Mess[i] = chat
	//	}
	//}
	//s.Forward("0", "success", Mess)
}
func GetMessageNum() (Message int) {
	idlist := []string{}
	err := JsRedis.Redis_get(constant.ChatList, &idlist)
	if err != nil {
		JsLogger.Error(err.Error())
	}
	for _, v := range idlist {
		chat := Chatnum{}
		err := JsRedis.Redis_hdbget(0, constant.Chat, v, &chat)
		if err != nil {
			JsLogger.Error(err.Error())
		}
		if chat.Numbus > 0 {
			Message++
		}
	}
	return Message
}
