package com

//在文章或产品的下面有文章和产品标签
//不同文章产品显示它自己的标签，
//点击标签跳转到列表页面，文章或产品列表

import (
	"JsGo/JsHttp"
	"JsGo/JsLogger"
	"fmt"
)

func Find() {
	JsHttp.Http("/tagquerylinkap", TagQueryLinkAP) //通过标签查询文章或产品
}

//由前端查询调用网络接口
//by tag query  article product标签查询对应的文章或产品链接返回内容或列表
func TagQueryLinkAP(s *JsHttp.Session) {
	type Para struct {
		Tag          string
		AP           bool   //FALSE=文章或TRUE=产品
		Num          int    //所需要查找的数量
		startID      string //起始查找的ID，如果为空表示从第一个查找
		ListOrDetail bool   //要返回的是列表还是详细内容，false=详情，TRUE=列表
	}
	idList := []string{}              //放ID的集合(正常ID)
	artCommonList := []*XM_Contents{} //放内容详情的集合content
	proCommonList := []*XM_Product{}  //放产品内容详情的集合
	var Needid string                 //循环暂存Needid
	var RecordIdArt string            //记录标签中文章ID避免删除找不到下一个
	var RecordIdPro string            //记录标签中产品ID避免删除找不到下一个
	var upidt string                  //上一个IP地址缓存

	st := &Para{}
	if err := s.GetPara(st); err != nil {
		s.Forward("1", err.Error(), nil)
		return
	}
	// 标签名字或数量,有一个不满足返回
	if st.Num <= 0 || st.Tag == "" {
		info := fmt.Sprintf("TagQueryLinkAP,tag=%s,Num=%d\n", st.Tag, st.Num)
		JsLogger.Error(info)
		s.Forward("1", info, nil)
		return
	}

	//取10个内容详情，两种取法，只有第一个取法不同，
	///可以把这两种取法的第一个单独拉出去,之后的几种一起处理。

	if st.AP == false { //如果是零FALSE，查找文章//文章调用
		if st.startID == "" { //如果ID为空
			d, eb := GetOneTag(st.Tag) //获取标签内容放到d中
			if eb != nil {
				JsLogger.Error("tag detail unfild") //标签中的内容相关ID找不到
				s.Forward("5", eb.Error(), nil)
				return
			}
			Needid = d.TagFinArt //取标签中的ID
			RecordIdArt = d.TagFinArt
		} else {
			upidt = st.startID
			Needid = TagNextA(st.Tag, st.startID) //取出入ID的下一个ID
		}
		for i := 0; i <= st.Num; { //循环获取详细内容
			if Needid == "" { //如果去到的ID为空说明上一次已经是最后一个
				return
			}
			//在每次查找时都去读取其内容
			dataart, err := Getcontent(Needid) //查询文章（内容），对应的详情。
			if err != nil {
				JsLogger.Error("detail fild error")
				s.Forward("6", err.Error(), nil)
				return
			}
			if dataart.DelTag && Needid != RecordIdArt { //并且判断中的关键字，
				//是否被标记为删除//如果被删除了，就返回顶部再去获取
				err := DelteTageConnectA(st.Tag, upidt, Needid)
				if err != nil {
					JsLogger.Info("detail failure Article")
					return
				}
			} else {
				i++ //找到正常的然后进行加一
				idList = append(idList, Needid)
				artCommonList = append(artCommonList, dataart)
			}
			upidt = Needid
			Needid = TagNextA(st.Tag, upidt) //读取下一个到的ID放在Needid中
		}
		if !st.ListOrDetail {
			s.Forward("0", "success", artCommonList)
		}
	}
	if st.AP == true { //如果是零TRUE，查找产品//产品调用
		if st.startID == "" { //如果ID为空
			d, eb := GetOneTag(st.Tag) //获取标签内容放到d中
			if eb != nil {
				JsLogger.Error("tag detail unfild") //标签中的内容相关ID找不到
				s.Forward("5", eb.Error(), nil)
				return
			}
			Needid = d.TagFinPro //取标签中的ID
			RecordIdPro = d.TagFinPro
		} else {
			Needid = TagNextA(st.Tag, st.startID) //取出入ID的下一个ID
		}
		for i := 0; i <= st.Num; { //循环获取详细内容
			if Needid == "" { //如果去到的ID为空说明上一次已经是最后一个
				break
			}
			//在每次查找时都去读取其内容
			dataPro, err := GetProductInfo(Needid) //查询产品对应的详情。
			if err != nil {
				JsLogger.Error("detail fild error")
				s.Forward("6", err.Error(), nil)
				return
			}
			if dataPro.DelTag && Needid != RecordIdPro { //并且判断中的关键字，
				//是否被标记为删除//如果被删除了，就返回顶部再去获取
				err := DelteTageConnectP(st.Tag, upidt, Needid)
				if err != nil {
					JsLogger.Info("detail failure Product")
				} else {
					i++ //找到正常的然后进行加一
					idList = append(idList, Needid)
					proCommonList = append(proCommonList, dataPro)
				}
				upidt = Needid
				Needid = TagNextA(st.Tag, upidt) //读取下一个到的ID放在Needid中
			}
		}

		if !st.ListOrDetail {
			s.Forward("0", "success", proCommonList)
		}
	}

	if st.ListOrDetail {
		s.Forward("0", "success", idList)
	}
}
