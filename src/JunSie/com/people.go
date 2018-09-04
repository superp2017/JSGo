package com

import (
	"JsGo/JsHttp"
	"JsGo/JsLogger"
	"JsGo/JsStore/JsRedis"
	"JunSie/constant"
	"fmt"
)

func Peoplerecord() {
	// JsHttp.Http("/recordlike", recordLike)             //2.记录用户赞过的//2.增加点赞数量（）
	// JsHttp.Http("/recordcollection", recordCollection) //3.记录用户收藏///3.增加收藏数量（我喜欢的文章产品） collec
	// JsHttp.Http("/recordattention", recordAttention)   //3.1记录关注（作者）////3.1增加关注数量(我关注的是医生，作者)
	// JsHttp.Http("/recordcomment", recordComment)       ////4.记录用户(评论过的)(X)（用于修改和删除评论）/4.增加评论数量
	JsHttp.WhiteHttps("/getrecorduser", getRecordUser) //5.获取某个用户的所有统计信息（用户加所需要的东西字段）
	// JsHttp.Http("/recordpageview", recordPageView)     //1.记录用户(浏览)，//（猜你喜欢 ）
}

// rec := &PersonalRecord{}
////5.获取某个用户的所有统计信息
func getRecordUser(s *JsHttp.Session) {

	type rec struct {
		UID  string //用户ID
		TYPE string //请求种
	}
	// "MyLikeArt"       //2.01我赞过的文章
	// "MyLikePro"       //2.02我赞过的产品
	// "MyCollectionArt" //3.01我收藏的文章
	// "MyCollectionPro" //3.02我收藏的产品
	// "MyAttention"     //3.1记录关注（作者）
	// "MyComment"       //4.记录(评论过的)(X)
	// "MyPageViewPro"      //1.我浏览过的产品
	// "MyPageViewArt"      //1.我浏览过的文章
	st := &rec{}
	err := s.GetPara(st)
	fmt.Printf("%v\n", st)
	if err != nil {
		info := "get para error:" + err.Error()
		JsLogger.Error(info)
		s.Forward("1", info, st)
		return
	}
	data := make([]string, 0)

	if err := JsRedis.Redis_hdbget(9, st.TYPE, st.UID, &data); err != nil {
		info := "get error" + err.Error()
		JsLogger.Error(info)
		s.Forward("1", info, data)
		return
	}
	s.Forward("0", "success", data)
}

//1.记录用户(浏览)，///1.增加访问量-------------add page view
const Mexlen = 15

func RecordPageView(Regulation bool, UID, IDap string, Type int8) error {
	//Regulation增加或减少。IDap，ID。Type产品还是文章
	data := []string{} //浏览条数
	var ht string
	if Type == 1 { //文章
		ht = constant.MYPageViewArt
	}
	if Type == 2 { //产品
		ht = constant.MYPageViewPro
	}
	JsRedis.Redis_hdbget(9, ht, UID, &data)
	exist := false
	index := 0
	for i, v := range data { //判断有没有在，有的话在那个位置
		if v == IDap {
			exist = true
			index = i
			break
		}
	}

	if Regulation {
		//追加到头部尾部溢出

		if exist || index != 0 { //如果存在且不是第一个那么，放在第一个位置
			data = append(data[:index], data[index+1:]...)
			temp := append([]string{}, IDap)
			data = append(temp, data...)

		} else { //为空，或没找到添加过，放在第一个。
			temp := append([]string{}, IDap)
			data = append(temp, data...)
			if len(data) > Mexlen { //如果大于最大设定量，则取前面的
				data = data[:Mexlen]
			}
		}

	}
	if !Regulation {
		if exist == true {
			data = append(data[:index], data[index+1:]...)
		}
	}
	if err := JsRedis.Redis_hdbset(9, ht, UID, &data); err != nil {
		info := "set error:" + err.Error()
		JsLogger.Error(info)
		return err
	}
	return nil
}

//2.记录用户赞过的//2.增加点赞数量（）
func RecordLike(Regulation bool, UID, IDap string, Type int8) error {
	//Regulation增加或减少。IDap，ID。Type产品还是文章
	data := make([]string, 0)
	var ht string
	if Type == 1 { //文章
		ht = constant.MYLikeArt
	}
	if Type == 2 { //产品
		ht = constant.MYLikePro
	}
	JsRedis.Redis_hdbget(9, ht, UID, &data)
	index := -1
	for i, v := range data { //判断有没有在，有的话在那个位置
		if v == IDap {
			index = i
			break
		}
	}
	if Regulation {
		if index == -1 { //为空，或没找到添加过，
			temp := data
			data := make([]string, 0, len(temp)+1)
			data = append(data, IDap)
			data = append(data, temp...)
		}
	} else {
		if index != -1 {
			data = append(data[:index], data[index+1:]...)
		}
	}
	if err := JsRedis.Redis_hdbset(9, ht, UID, &data); err != nil {
		info := "set error:" + err.Error()
		JsLogger.Error(info)
		return err
	}
	return nil

}

//3.记录用户收藏///3.增加收藏数量（我喜欢的文章产品） collec
func RecordCollection(Regulation bool, UID, IDap string, Type int8) error { //记录收藏
	data := make([]string, 0)
	var ht string
	if Type == 1 { //文章
		ht = constant.MYCollectionArt
	}
	if Type == 2 { //产品
		ht = constant.MYCollectionPro
	}
	JsRedis.Redis_hdbget(9, ht, UID, &data)
	index := -1
	for i, v := range data { //判断有没有在，有的话在那个位置
		if v == IDap {
			index = i
			break
		}
	}
	if Regulation {
		if index == -1 { //为空，或没找到添加过，
			temp := data
			data = make([]string, 0, len(temp)+1)
			data = append(data, IDap)
			data = append(data, temp...)
		}
	} else {
		if index != -1 {
			data = append(data[:index], data[index+1:]...)
		}
	}
	if err := JsRedis.Redis_hdbset(9, ht, UID, &data); err != nil {
		info := "set error:" + err.Error()
		JsLogger.Error(info)
		return err
	}
	return nil
}

// b, e := JsRedis.Redis_hdbexists(9, UID, K) //exist存在
// if e != nil {
// 	JsLogger.Error(e.Error())
// 	return
// }
// if b {
// 	err := JsRedis.Redis_hdbget(9, UID, K, &data)
// 	if err != nil {
// 		info := "get error:" + err.Error()
// 		JsLogger.Error(info)
// 		return
// 	}
// }

// if Regulation {
// 	data = append(data, IDap)
// 	if err := JsRedis.Redis_hdbset(9, UID, K, &data); err != nil {
// 		info := "set error:" + err.Error()
// 		JsLogger.Error(info)
// 		return
// 	}
// } else {
// 	dataT := make([]string, 0, len(data))
// 	for i, _ := range data {
// 		if data[i] != IDap {
// 			dataT = append(dataT, data[i])
// 		}
// 	}
// 	if err := JsRedis.Redis_hdbset(9, UID, K, &dataT); err != nil {
// 		info := "set error:" + err.Error()
// 		JsLogger.Error(info)
// 		return
// 	}
// }
// return
// }

// for i, _ := range data {
// 		var dataT []string
// 		if data[i] == IDap {
// 			dataT = append(data[:i], data[i+1:]...)
// 			if err := JsRedis.Redis_hdbset(9, UID, K, &dataT); err != nil {
// 				info := "set error:" + err.Error()
// 				JsLogger.Error(info)
// 				return
// 			}
// 			break;
// 		}

// 	}

// //3.1记录关注（作者）////3.1增加关注数量(我关注的是医生，作者)
// func RecordAttention(Regulation bool, UID string, IDpeo string) {

// }

// //4.记录用户(评论过的)(X)（用于修改和删除评论）/4.增加评论数量
// func RecordComment(Regulation bool, UID string, IDpeo string) {
// 	type people struct {
// 		Regulation bool   //增加为TRUE，减少为false
// 		UID        string //用户ID
// 		// UserHead //用户头像
// 		IDpro string //产品ID
// 		IDart string //文章ID

// 	}

// }
