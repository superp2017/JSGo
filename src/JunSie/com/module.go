package com

import (
	"JsGo/JsBench/JsModule"
	"JsGo/JsHttp"
	"JsGo/JsLogger"
	"JsGo/JsStore/JsRedis"

	"JunSie/constant"
	"JunSie/util"
	"fmt"
	"time"
)

func InitModule() {
	JsHttp.WhiteHttp("/creathomemodule", CreatHomeModule)                       //新建、修改模块
	JsHttp.WhiteHttp("/gethomemodule", GetHomeModule)                           //获取首页模块
	JsHttp.WhiteHttp("/downhomemodule", DownHomeModule)                         //下架首页模块
	JsHttp.WhiteHttp("/getmodulelist", GetModuleList)                           //获取所有的模块列表，包含线上线下和配置的
	JsHttp.WhiteHttp("/setchoicenessmodulelist", SetChoicenessConfigModuleList) //设置更新用户精选配置的模块列表
	JsHttp.WhiteHttp("/sethomemodulelist", SetHomeConfigModuleList)             //设置更新用户首页配置的模块列表
	JsHttp.WhiteHttp("/onlinemodule", OnlineModule)                             //上线某一个模块
}
func InitModuleMall() {
	JsHttp.WhiteHttps("/gethomemodule", GetHomeModule)                           //获取首页模块
	JsHttp.WhiteHttps("/getmodulelist", GetModuleList)                           //获取所有的模块列表，包含线上线下和配置的
	JsHttp.WhiteHttps("/setchoicenessmodulelist", SetChoicenessConfigModuleList) //设置更新用户精选配置的模块列表
	JsHttp.WhiteHttps("/sethomemodulelist", SetHomeConfigModuleList)             //设置更新用户首页配置的模块列表

}

type ModuleInfo struct {
	Name   string
	PageID string
}

type ModuleList struct {
	Online     []ModuleInfo //线上的
	Preview    []ModuleInfo //预览的
	Home       []ModuleInfo //当前首页配置的
	Choiceness []ModuleInfo //当前精选配置的
}

func (this *ModuleList) appendOnline(PageID string, Name string) {
	exist := false
	for _, v := range this.Online {
		if v.PageID == PageID {
			exist = true
			break
		}
	}
	if !exist {
		this.Online = append(this.Online, ModuleInfo{
			Name:   Name,
			PageID: PageID,
		})
	}
}

func (this *ModuleList) appendPreview(PageID string, Name string) {
	exist := false
	for _, v := range this.Preview {
		if v.PageID == PageID {
			exist = true
			break
		}
	}
	if !exist {
		this.Preview = append(this.Preview, ModuleInfo{
			Name:   Name,
			PageID: PageID,
		})
	}
}

func CreatHomeModule(session *JsHttp.Session) {
	st := &JsModule.Module{}
	if err := session.GetPara(st); err != nil {
		JsLogger.Error(err.Error())
		session.Forward("1", err.Error(), nil)
		return
	}
	if st.PageID == "" || st.Name == "" {
		str := fmt.Sprintf("CreatHomeModule failed,PageID=%s,Name=%s\n", st.PageID, st.Name)
		JsLogger.Error(str)
		session.Forward("1", str, nil)
		return
	}
	data := &JsModule.Module{}
	key := make(map[string]string)
	exist := false
	JsRedis.Redis_hget(constant.HomeModule, constant.ModuleMap, &key)
	if moduleID, ok := key[st.PageID]; ok && moduleID != "" {
		if err := JsRedis.Redis_hget(constant.HomeModule, moduleID, data); err == nil {
			exist = true
		}
	}
	if exist {
		st.ModuleID = data.ModuleID
		st.Status = data.Status
		st.CreatTime = data.CreatTime
	} else {
		st.ModuleID = util.IDer(constant.HomeModule)
		st.CreatTime = util.CurTime()
	}
	st.Status = "0"
	if err := JsRedis.Redis_hset(constant.HomeModule, st.ModuleID, st); err != nil {
		JsLogger.Error(err.Error())
		session.Forward("1", err.Error(), nil)
		return
	}
	if !exist {
		key[st.PageID] = st.ModuleID
		if err := JsRedis.Redis_hset(constant.HomeModule, constant.ModuleMap, &key); err != nil {
			JsLogger.Error(err.Error())
			session.Forward("1", err.Error(), nil)
			return
		}
	}

	if st.Islimit && st.Limit != nil {
		cur := time.Now().Unix()
		if st.Limit.EndTime == 0 || st.Limit.StartTime < cur || st.Limit.EndTime < st.Limit.StartTime || st.Limit.EndTime <= cur {
			str := fmt.Sprintf("限时模块创建失败,StartTime=%d,EndTime=%d", st.Limit.StartTime, st.Limit.EndTime)
			JsLogger.Error(str)
			session.Forward("1", str, nil)
			return
		}
	}

	list := &ModuleList{}
	JsRedis.Redis_hget(constant.HomeModule, constant.KEY_ModuleConfig, &list)
	list.appendPreview(st.PageID, st.Name)
	go JsRedis.Redis_hset(constant.HomeModule, constant.KEY_ModuleConfig, &list)

	session.Forward("0", "creat success", st)
}

//获取首页模块
func GetHomeModule(session *JsHttp.Session) {
	type Para struct {
		PageID string //版面号
	}
	st := &Para{}
	if err := session.GetPara(st); err != nil {
		JsLogger.Error(err.Error())
		session.Forward("1", err.Error(), nil)
		return
	}
	if st.PageID == "" {
		str := "GetHomeModule failed,PageID is empty\n"
		session.Forward("1", str, JsLogger.ErrorLog(str))
		return
	}
	data := &JsModule.Module{}
	key := make(map[string]string)
	if err := JsRedis.Redis_hget(constant.HomeModule, constant.ModuleMap, &key); err != nil {
		JsLogger.Error("Redis_hget(ModuleMap) failed:%s\n", err.Error())
		session.Forward("1", err.Error(), nil)
		return
	}
	if ModuleID, ok := key[st.PageID]; ok {
		if err := JsRedis.Redis_hget(constant.HomeModule, ModuleID, data); err != nil {
			JsLogger.Error("Redis_hget(ModuleID) failed:%s\n", err.Error())
			session.Forward("1", err.Error(), nil)
			return
		}
		if data.Status == "1" {
			session.Forward("1", "get failed,Module Status is 1", data)
			return
		}
		if data.Islimit && data.Limit != nil {
			cur := time.Now().Unix()
			if data.Limit.EndTime == 0 || data.Limit.StartTime < cur || data.Limit.EndTime < data.Limit.StartTime || data.Limit.EndTime <= cur {
				str := fmt.Sprintf("限时模块已经超时,StartTime=%d,EndTime=%d,CurTime=%d", data.Limit.StartTime, data.Limit.EndTime, cur)
				JsLogger.Error(str)
				session.Forward("1", str, nil)
				data.Status = "1"
				go JsRedis.Redis_hset(constant.HomeModule, ModuleID, data)
				return
			}
		}
		for i, v := range data.Data {
			if v.Type == "0" {
				//if pro, err := product.GetProductInfo(v.ID); err == nil {
				//	data.Data[i].Name = pro.Title
				//	data.Data[i].SubName = pro.SubTitle
				//	data.Data[i].Price = pro.NowPrice
				//	data.Data[i].OrigPrice = pro.OriPrice
				//	if data.Data[i].Pic == "" && len(pro.Images) > 0 {
				//		data.Data[i].Pic = pro.Images[0]
				//	}
				//	data.Data[i].Tags = product.GetProTags(v.ID)
				//}
				if sta, err := ProStatics(v.ID); err == nil {
					data.Data[i].VisitNum = sta.VisitNum
					data.Data[i].PraiseNum = sta.PraiseNum
					data.Data[i].AttentionNums = sta.AttentionNums
					data.Data[i].CommentNum = sta.CommentNum
				}
			} else {
				//if con, err := content.Getcontent(v.ID); err == nil {
				//	data.Data[i].UID = con.UID
				//	data.Data[i].UserHead = con.UserHead
				//	data.Data[i].UserName = con.Author
				//	data.Data[i].Name = con.Title
				//	data.Data[i].SubName = con.SubTitle
				//	if data.Data[i].Pic == "" && len(con.Images) > 0 {
				//		data.Data[i].Pic = con.Images[0]
				//	}
				//	data.Data[i].Tags = content.ContentTags(v.ID)
				//}
				if sta, err := ContentStatics(v.ID); err == nil {
					data.Data[i].VisitNum = sta.VisitNum
					data.Data[i].PraiseNum = sta.PraiseNum
					data.Data[i].AttentionNums = sta.AttentionNums
					data.Data[i].CommentNum = sta.CommentNum
				}
			}
		}
		session.Forward("0", "get success", data)
		return
	}

	session.Forward("1", "get faild not exist such module", nil)
}

func DownHomeModule(session *JsHttp.Session) {
	type Para struct {
		PageID string //版面号
	}
	st := &Para{}
	if err := session.GetPara(st); err != nil {
		session.Forward("1", err.Error(), nil)
		return
	}
	if st.PageID == "" {
		str := "DownHomeModule failed,PageID is empty\n"
		session.Forward("1", str, JsLogger.ErrorLog(str))
		return
	}
	data := &JsModule.Module{}
	key := make(map[string]string)
	if err := JsRedis.Redis_hget(constant.HomeModule, constant.ModuleMap, &key); err != nil {
		JsLogger.Error("DownHomeModule Redis_hget(ModuleMap) failed:%s\n", err.Error())
		session.Forward("1", err.Error(), nil)
		return
	}
	if ModuleID, ok := key[st.PageID]; ok {
		if err := JsRedis.Redis_hget(constant.HomeModule, ModuleID, data); err != nil {
			JsLogger.Error("DownHomeModule Redis_hget(ModuleID) failed:%s\n", err.Error())
			session.Forward("1", err.Error(), nil)
			return
		}
		data.Status = "1"
		if err := JsRedis.Redis_hset(constant.HomeModule, ModuleID, data); err != nil {
			JsLogger.Error("DownHomeModule Redis_hget(ModuleID) failed:%s\n", err.Error())
			session.Forward("1", err.Error(), nil)
			return
		}
		session.Forward("0", "down module success\n", data)
		return
	}
	session.Forward("1", "DownHomeModule failed, module is not exist\n", nil)
}

//获取首页模块配置的列表
func GetModuleList(session *JsHttp.Session) {
	list := &ModuleList{}
	if err := JsRedis.Redis_hget(constant.HomeModule, constant.KEY_ModuleConfig, list); err != nil {
		str := "GetModuleList failed:" + err.Error()
		JsLogger.Error(str)
		session.Forward("0", str, list)
		return
	}
	session.Forward("0", "GetModuleList success\n", list)
}

//更新首页模块的配置列表
func SetHomeConfigModuleList(session *JsHttp.Session) {
	setconfig(session, true)
}

//更新精选模块的配置列表
func SetChoicenessConfigModuleList(session *JsHttp.Session) {
	setconfig(session, false)
}

func setconfig(session *JsHttp.Session, isHome bool) {
	type Para struct {
		List []ModuleInfo
	}
	st := &Para{}
	if err := session.GetPara(st); err != nil {
		JsLogger.Error(err.Error())
		session.Forward("1", err.Error(), nil)
		return
	}
	list := &ModuleList{}
	if err := JsRedis.Redis_hget(constant.HomeModule, constant.KEY_ModuleConfig, list); err != nil {
		str := "SetConfigModuleList Redis_hget failed:" + err.Error()
		JsLogger.Error(str)
		session.Forward("1", str, nil)
		return
	}
	if isHome {
		list.Home = st.List
	} else {
		list.Choiceness = st.List
	}
	if err := JsRedis.Redis_hset(constant.HomeModule, constant.KEY_ModuleConfig, list); err != nil {
		str := "SetConfigModuleList Redis_hset failed:" + err.Error()
		JsLogger.Error(str)
		session.Forward("1", str, nil)
		return
	}
	session.Forward("0", "SetConfigModuleList success\n", list)
}

func OnlineModule(session *JsHttp.Session) {
	type Para struct {
		PageID     string //正式上线的id
		PageID_Pre string //预览的id
	}
	st := &Para{}
	if err := session.GetPara(st); err != nil {
		JsLogger.Error("OnlineModule faild:" + err.Error())
		session.Forward("1", "OnlineModule faild:"+err.Error(), nil)
		return
	}

	if st.PageID == "" || st.PageID_Pre == "" || st.PageID == st.PageID_Pre {
		str := fmt.Sprintf("OnlineModule failed ID is empty:PageID=%s,PageID_Pre=%s\n", st.PageID, st.PageID_Pre)
		JsLogger.Error(str)
		session.Forward("1", str, nil)
		return
	}

	key := make(map[string]string)
	if err := JsRedis.Redis_hget(constant.HomeModule, constant.ModuleMap, &key); err != nil {
		JsLogger.Error("Redis_hget(ModuleMap) failed:%s\n", err.Error())
		session.Forward("1", err.Error(), nil)
		return
	}
	ModuleID, ok := key[st.PageID_Pre]
	if !ok || ModuleID == "" {
		str := fmt.Sprintf("OnlineModule ,PageID_Pre =%s not exist module\n", st.PageID_Pre)
		JsLogger.Error(str)
		session.Forward("1", str, nil)
		return
	}
	data := &JsModule.Module{}
	if err := JsRedis.Redis_hget(constant.HomeModule, ModuleID, data); err != nil {
		JsLogger.Error(err.Error())
		session.Forward("1", err.Error(), nil)
		return
	}
	////////////////更新map
	key[st.PageID] = ModuleID
	delete(key, st.PageID_Pre)

	/////更新配置列表
	list := &ModuleList{}
	JsRedis.Redis_hget(constant.HomeModule, constant.KEY_ModuleConfig, &list)
	list.appendOnline(st.PageID, data.Name)
	list.appendPreview(st.PageID_Pre, data.Name)
	/////更新模块详情
	data.PageID = st.PageID
	if err := JsRedis.Redis_hset(constant.HomeModule, ModuleID, data); err != nil {
		JsLogger.Error(err.Error())
		session.Forward("1", err.Error(), nil)
		return
	}

	if err := JsRedis.Redis_hset(constant.HomeModule, constant.ModuleMap, &key); err != nil {
		JsLogger.Error(err.Error())
		session.Forward("1", err.Error(), nil)
		return
	}

	if err := JsRedis.Redis_hset(constant.HomeModule, constant.KEY_ModuleConfig, &list); err != nil {
		JsLogger.Error(err.Error())
		session.Forward("1", err.Error(), nil)
		return
	}

	session.Forward("0", "OnlineModule success", list)
}
