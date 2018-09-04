package com

import (
	. "JsGo/JsLogger"
	"JsGo/JsStore/JsRedis"
	"JunSie/constant"
	_ "time"
)

//tag link article 建立标签查找文章链接
func TagLinkA(NAid string, tags []string) { //文章ID  多个标签

	for _, v := range tags { //标签的循环，
		//从数据库中查找(标签关联产品列表）内容中的K
		tagk := v + "_" + NAid
		b, ea := JsRedis.Redis_hexists(constant.TAGPRO, tagk) //exist存在
		if ea != nil {
			Error("2.tagFpro.go", " info  set Tag find error")
			return
		}
		if b {
			Info("3.tagFpro.go", "already set Tag find")
			return
		}

		d, eb := GetOneTag(v) //获取标签内容放到d中
		if eb != nil {
			Error("1.tagFpro.go", " get tag error")
			return
		}
		idOld := d.TagFinArt //从内容中取出 TagFinPro //产品ID旧的

		//创建一个关联的表放在数据库 NPid关联取出来的 (TagFinPro //产品ID)
		errf := JsRedis.Redis_hset(constant.TAGPRO, tagk, idOld)
		if errf != nil {
			Info("4.tagFpro.go ", "tag connect errf")
			return
		}
		d.TagFinArt = NAid //更新标签中的最新ID
		d.CountArt++
		JsRedis.Redis_hset(constant.TAG, v, d)
		Info("0", "success", "")
	}

}

//收到一个产品ID,获取其内容中关联的ID
func TagNextA(tag, id string) (needID string) {
	tagk := tag + "_" + id
	errf := JsRedis.Redis_hget(constant.TAGART, tagk, needID)
	if errf != nil {
		Info("5.tagFart.go ", "Server cannot find,Nothing else")
		return
	}
	Info("0", "success", "")
	return
}

//修改文章链接表，修改产品链接表。//  modifyLinkp //modifylink product修改链接
func DelteTageConnectA(tag, uptagid, delid string) error {
	threetagid := TagNextA(tag, delid)
	//第一步，建立与要删除的下一个链接。(重复自动对换)
	tagko := tag + "_" + uptagid
	err := JsRedis.Redis_hset(constant.TAGART, tagko, threetagid)
	if err != nil {
		return err
	}
	Info("connect success")
	//第二步，删除链接
	tagkt := tag + "_" + delid
	e := JsRedis.Redis_hdel(constant.TAGART, tagkt)
	if e != nil {
		return e
	}
	Info("delete success")
	return nil

}
