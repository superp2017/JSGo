package com

import (
	. "JsGo/JsLogger"
	"JsGo/JsStore/JsRedis"
	"JunSie/constant"
	_ "time"
)

//获取一个标签的内容
func GetOneTag(tag string) (*Tagnr, error) {
	st := &Tagnr{} //st存放标签内容
	err := JsRedis.Redis_hget(constant.TAG, tag, st)
	return st, err
}

//tag link product  建立标签查找产品链接
func TagLinkP(npid string, tags []string) { //产品ID  多个标签
	if tags == nil {
		return
	}

	for _, v := range tags { //标签的循环，

		//从数据库中(标签关联产品列表）内容中的K
		tagk := v + "_" + npid
		b, e := JsRedis.Redis_hexists(constant.TAGPRO, tagk) //exist存在
		if e != nil {
			Error("2.tagFpro.go", " info already set Tag find")
			return
		}
		if b {
			Info("3.tagFpro.go", "already set Tag find")
			return
		}

		d, err := GetOneTag(v) //获取标签内容
		if err != nil {
			Info("1.tagFpro.go", " get tag error")
			return
		}
		idOld := d.TagFinPro //从内容中取出 TagFinPro //产品ID旧的

		//创建一个关联的表放在数据库 NPid关联取出来的 (TagFinPro //产品ID)
		errf := JsRedis.Redis_hset(constant.TAGPRO, tagk, idOld)
		if errf != nil {
			Info("4.tagFpro.go ", "tag connect errf")
			return
		}
		d.TagFinPro = npid //更新标签中的最新ID
		d.CountPro++
		err = JsRedis.Redis_hset(constant.TAG, v, d)
		if err != nil {
			Error(err.Error())
		}
		return
	}

}

//收到一个产品ID,获取其内容中关联的ID
func TagNextP(tag, id string) (needID string) {
	tagk := tag + "_" + id
	err := JsRedis.Redis_hget(constant.TAGPRO, tagk, needID)
	if err != nil {
		Info("5.tagFpro.go ", "Server cannot find,Nothing else")
		return
	}
	Info("TagNextP success")
	return
}

//修改文章链接表，修改产品链接表。//  modifyLinkp //modifylink product修改链接
func DelteTageConnectP(tag, uptagid, delid string) error {
	threetagid := TagNextA(tag, delid)
	//第一步，建立与要删除的下一个链接。(重复自动对换)
	tagko := tag + "_" + uptagid
	err := JsRedis.Redis_hset(constant.TAGPRO, tagko, threetagid)
	if err != nil {
		return err
	}
	Info("connect success")
	//第二步，删除链接
	tagkt := tag + "_" + delid
	e := JsRedis.Redis_hdel(constant.TAGPRO, tagkt)
	if e != nil {
		return e
	}
	Info("delete success")
	return nil

}
