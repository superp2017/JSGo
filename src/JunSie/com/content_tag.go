package com

import (
	"JsGo/JsHttp"
	"JsGo/JsLogger"
	"JsGo/JsStore/JsRedis"
	"fmt"
)

//获取内容的所有标签
func GetContentTags(s *JsHttp.Session) {
	type Para struct {
		ID string //产品id
	}
	st := &Para{}
	if err := s.GetPara(st); err != nil {
		s.Forward("1", err.Error(), nil)
		return
	}
	if st.ID == "" {
		info := "GetContentTags ID is empty\n"
		JsLogger.Error(info)
		s.Forward("1", info, nil)
		return
	}
	s.Forward("0", "success", ContentTags(st.ID))

}

//删除内容标签
func RemoveContentTag(s *JsHttp.Session) {
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
		info := fmt.Sprintf("RemoveContentTag failed,ID=%s,Tag=%s\n", st.ID, st.Tag)
		JsLogger.Error(info)
		s.Forward("1", info, nil)
		return
	}
	if err := removeContentTag(st.ID, st.Tag); err != nil {
		s.Forward("1", err.Error(), nil)
		return
	}
	s.Forward("0", "success", nil)
}

//内容关联标签
func Content2Tag(id string, tags []string) {
	for _, v := range tags {
		if err := JsRedis.Redis_Sset(id, v); err != nil {
			JsLogger.Error(err.Error())
		}
	}
}

//获取内容所有的标签
func ContentTags(id string) (list []string) {
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
func removeContentTag(id, tag string) error {
	return JsRedis.Redis_Sdel(id, tag)
}

//更新产品标签
func changeContentTag(id string, tags []string) {
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
			go removeContentTag(id, v)
		}
	}
	if len(tags) > 0 {
		go Content2Tag(id, tags)
	}
}

//由前端查询调用网络接口
//by tag query  article 标签查询对应的文章链接 返回内容或列表
// func TagQueryLinkA(s *JsHttp.Session) {
// 	type Para struct {
// 		Tag          string
// 		AP           bool   //FALSE=文章或TRUE=产品
// 		Num          int    //所需要查找的数量
// 		startID      string //起始查找的ID，如果为空表示从第一个查找
// 		ListOrDetail bool   //要返回的是列表还是详细内容，false=详情，TRUE=列表
// 	}
// 	idList := []string{}              //放ID的集合(正常ID)
// 	artCommonList := []*XM_Contents{} //放内容详情的集合content
// 	var Needid string                 //循环暂存Needid
// 	var RecordIdArt string            //记录标签中文章ID避免删除找不到下一个
// 	var upidt string                  //上一个IP地址缓存

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
// 		Needid = d.TagFinArt //取标签中的ID
// 		RecordIdArt = d.TagFinArt
// 	} else {
// 		upidt = st.startID
// 		Needid = tage.TagNextA(st.Tag, st.startID) //取出入ID的下一个ID
// 	}
// 	for i := 0; i <= st.Num; { //循环获取详细内容
// 		if Needid == "" { //如果去到的ID为空说明上一次已经是最后一个
// 			return
// 		}
// 		//在每次查找时都去读取其内容
// 		dataart, err := Getcontent(Needid) //查询文章（内容），对应的详情。
// 		if err != nil {
// 			JsLogger.Error("detail fild error")
// 			s.Forward("6", err.Error(), nil)
// 			return
// 		}

// 		if dataart.DelTag && Needid != RecordIdArt { //并且判断中的关键字，
// 			//是否被标记为删除//如果被删除了，就返回顶部再去获取
// 			err := tage.DelteTageConnectA(st.Tag, upidt, Needid)
// 			if err != nil {
// 				JsLogger.Info("detail failure Article")
// 				return
// 			}

// 		} else {
// 			i++ //找到正常的然后进行加一
// 			idList = append(idList, Needid)
// 			artCommonList = append(artCommonList, dataart)
// 		}
// 		upidt = Needid
// 		Needid = tage.TagNextA(st.Tag, upidt) //读取下一个到的ID放在Needid中

// 	}
// 	if !st.ListOrDetail {
// 		s.Forward("0", "success", artCommonList)
// 	} else {
// 		s.Forward("0", "success", idList)
// 	}
// }
