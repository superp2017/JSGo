package com

import (
	"JsGo/JsHttp"
	"JsGo/JsLogger"
	"JsGo/JsStore/JsRedis"
	"fmt"
)

//获取产品的所有标签
func GetProductTags(s *JsHttp.Session) {
	type Para struct {
		ID string //产品id
	}
	st := &Para{}
	if err := s.GetPara(st); err != nil {
		s.Forward("1", err.Error(), nil)
		return
	}
	if st.ID == "" {
		info := "GetProductTags ID is empty\n"
		JsLogger.Error(info)
		s.Forward("1", info, nil)
		return
	}
	s.Forward("0", "success", GetProTags(st.ID))

}

//删除产品标签
func RemoveProductTag(s *JsHttp.Session) {
	type Para struct {
		ID  string //产品id
		Tag string //标签
	}
	st := &Para{}
	if err := s.GetPara(st); err != nil {
		s.Forward("1", err.Error(), nil)
		return
	}
	if st.ID == "" || st.Tag == "" {
		info := fmt.Sprintf("RemoveProductTag failed,ID=%s,Tag=%s\n", st.ID, st.Tag)
		JsLogger.Error(info)
		s.Forward("1", info, nil)
		return
	}
	if err := removeProTag(st.ID, st.Tag); err != nil {
		s.Forward("1", err.Error(), nil)
		return
	}
	s.Forward("0", "success", nil)
}

//产品关联标签
func pro2Tag(id string, tags []string) {
	for _, v := range tags {
		if err := JsRedis.Redis_Sset(id, v); err != nil {
			JsLogger.Error(err.Error())
		}
	}
}

//获取产品所有的标签
func GetProTags(id string) (list []string) {
	d, err := JsRedis.Redis_Sget(id)
	for _, v := range d {
		list = append(list, string(v.([]byte)))
	}
	if err != nil {
		JsLogger.Error(err.Error())
	}
	return
}

//移除一个产品的标签
func removeProTag(id, tag string) error {
	return JsRedis.Redis_Sdel(id, tag)
}

//更新产品标签
func changeProTag(id string, tags []string) {
	data := []string{}
	d, _ := JsRedis.Redis_Sget(id)
	for _, v := range d {
		data = append(data, string(v.([]byte)))
	}
	for _, v := range data {
		exist := false
		for _, v1 := range tags {
			if v == v1 {
				exist = true
				break
			}
		}
		if !exist {
			go removeProTag(id, v)
		}
	}
	if len(tags) > 0 {
		go pro2Tag(id, tags)
	}
}

/////////////////////////////////////////////////////////////////////////////////////////
//由前端查询调用网络接口
// //by tag query  article product标签查询对应的文章或产品链接返回内容或列表
// func TagQueryLinkP(s *JsHttp.Session) {
// 	type Para struct {
// 		Tag          string
// 		Num          int    //所需要查找的数量
// 		startID      string //起始查找的ID，如果为空表示从第一个查找
// 		ListOrDetail bool   //要返回的是列表还是详细内容，false=详情，TRUE=列表
// 	}
// 	idList := []string{}             //放ID的集合(正常ID)
// 	proCommonList := []*XM_Product{} //放产品内容详情的集合
// 	var Needid string                //循环暂存Needid
// 	var RecordIdPro string           //记录标签中产品ID避免删除找不到下一个
// 	var upidt string                 //上一个IP地址缓存

// 	st := &Para{}
// 	if err := s.GetPara(st); err != nil {
// 		s.Forward("1", err.Error(), nil)
// 		return
// 	}
// 	// 标签名字或数量,有一个不满足返回
// 	if st.Num <= 0 || st.Tag == "" {
// 		info := fmt.Sprintf("TagQueryLinkAP,tag=%s,Num=%d\n", st.Tag, st.Num)
// 		JsLogger.Error(info)
// 		s.Forward("1", info, nil)
// 		return
// 	}

// 	//取10个内容详情，两种取法，只有第一个取法不同，
// 	///可以把这两种取法的第一个单独拉出去,之后的几种一起处理。

// 	if st.startID == "" { //如果ID为空
// 		d, eb := tage.GetOneTag(st.Tag) //获取标签内容放到d中
// 		if eb != nil {
// 			JsLogger.Error("tag detail unfild") //标签中的内容相关ID找不到
// 			s.Forward("5", eb.Error(), nil)
// 			return
// 		}
// 		Needid = d.TagFinPro //取标签中的ID
// 		RecordIdPro = d.TagFinPro
// 	} else {
// 		Needid = tage.TagNextA(st.Tag, st.startID) //取出入ID的下一个ID
// 	}
// 	for i := 0; i <= st.Num; { //循环获取详细内容
// 		if Needid == "" { //如果去到的ID为空说明上一次已经是最后一个
// 			break
// 		}
// 		//在每次查找时都去读取其内容
// 		dataPro, err := GetProductInfo(Needid) //查询产品对应的详情。
// 		if err != nil {
// 			JsLogger.Error("detail fild error")
// 			s.Forward("6", err.Error(), nil)
// 			return
// 		}
// 		if dataPro.DelTag && Needid != RecordIdPro { //并且判断中的关键字，
// 			//是否被标记为删除//如果被删除了，就返回顶部再去获取
// 			err := tage.DelteTageConnectP(st.Tag, upidt, Needid)
// 			if err != nil {
// 				JsLogger.Info("detail failure Product")
// 			} else {
// 				i++ //找到正常的然后进行加一
// 				idList = append(idList, Needid)
// 				proCommonList = append(proCommonList, dataPro)
// 			}
// 			upidt = Needid
// 			Needid = tage.TagNextA(st.Tag, upidt) //读取下一个到的ID放在Needid中
// 		}
// 	}

// 	if !st.ListOrDetail {
// 		s.Forward("0", "success", proCommonList)
// 	} else {
// 		s.Forward("0", "success", idList)

// 	}
// }
