package com

import (
	"JsGo/JsBench/JsProduct"
	"JsGo/JsHttp"
	"JsGo/JsLogger"
	"JsGo/JsStore/JsRedis"
	"JunSie/constant"
	"JunSie/util"
	"fmt"
)

type Goods struct {
	ProID        string                  //产品idJ
	ProName      string                  //产品名称J
	Tags         []string                //标签
	ProFormat    JsProduct.ProductFormat //产品规格
	Nums         int                     //数量J
	BusinessID   string                  //商家id
	BusinessName string                  //商家名称
	ExData       map[string]string       //扩充的数据
	CreatTime    string                  //添加的时间
}

func InitShoppingCart() {
	JsHttp.WhiteHttps("/getshoppingcart", GetShoppingCart)                             //获取购物车所有商品
	JsHttp.WhiteHttps("/add2shoppingcart", AddToShoppingCart)                          //添加一个商品到购物车
	JsHttp.WhiteHttps("/removefromshoppingcart", RemoveFromShoppingCart)               //从购物车删除一个商品
	JsHttp.WhiteHttps("/removemorefromshoppingcart", RemoveMoreFromShoppingCart)       //从购物车删除多个商品（从购物车支付后删除购物车对应的产品）
	JsHttp.WhiteHttps("/modifygoodnumwithshoppingcart", ModifyGoodNumWithShoppingCart) //修改购物车内的商品数量
	JsHttp.WhiteHttps("/clearusershoppingcart", ClearUserShoppingCart)                 //清空购物车
}

type ShoppingCart struct {
	UID       string  //用户id`
	Data      []Goods //收藏的商品（购物车）
	CreatTime string  //创建时间
}

//获取所有的购物车
func GetShoppingCart(s *JsHttp.Session) {
	type Para struct {
		UID string //用户id
	}
	st := &Para{}
	if e := s.GetPara(st); e != nil {
		JsLogger.Error(e.Error())
		s.Forward("1", e.Error(), nil)
		return
	}
	if st.UID == "" {
		info := "GetShoppingCart failed,UID = nil "
		JsLogger.Error(info)
		s.Forward("2", info, nil)
		return
	}
	data := ShoppingCart{}
	if err := JsRedis.Redis_hget(constant.SHOPPINGCART, st.UID, &data); err != nil {
		JsLogger.Error(err.Error())
		s.Forward("3,", err.Error(), nil)
		return
	}
	s.Forward("0", "GetShoppingCart success", data)
}

//增加一个产品
func AddToShoppingCart(s *JsHttp.Session) {
	type Para struct {
		UID     string //用户信息
		Product Goods  //商品
	}
	st := &Para{}
	//判断数据完整性
	if e := s.GetPara(st); e != nil {
		JsLogger.Error(e.Error())
		s.Forward("1", e.Error(), nil)
		return
	}
	//逻辑正确性
	if st.Product.ProID == "" || (st.Product.ProFormat.Format != "" && st.Product.ProFormat.Price <= 0) || st.Product.Nums == 0 || st.UID == "" {
		info := fmt.Sprintf("AddToShoppingCart failed,ProID=%s,Nums=%d,UID=%s,Format=%s,Price=%d\n",
			st.Product.ProID, st.Product.Nums, st.UID, st.Product.ProFormat.Format, st.Product.ProFormat.Price)
		JsLogger.Error(info)
		s.Forward("2", info, nil)
		return
	}
	//判断是否存在用户（用户第一次往购物车里面加商品）
	data := ShoppingCart{}
	if err := JsRedis.Redis_hget(constant.SHOPPINGCART, st.UID, &data); err != nil {
		data.UID = st.UID
		data.CreatTime = util.CurTime()
	}
	//判断是否同款
	index := -1
	for i, v := range data.Data {
		if v.ProID == st.Product.ProID && v.ProFormat.Format == st.Product.ProFormat.Format {
			index = i
			break
		}
	}
	if index == -1 {
		st.Product.CreatTime = util.CurTime()
		data.Data = append(data.Data, st.Product) //不是同款追加
	} else {
		//追加数量
		data.Data[index].Nums += st.Product.Nums
	}
	data.UID = st.UID
	if err := JsRedis.Redis_hset(constant.SHOPPINGCART, st.UID, &data); err != nil {
		info := "Redis_hset error,:" + err.Error()
		JsLogger.Error(info)
		s.Forward("5", info, nil)
		return
	}
	s.Forward("0", "success", data)
}

//移除一个产品
func RemoveFromShoppingCart(s *JsHttp.Session) {
	type Para struct {
		UID    string //用户id
		ProID  string //产品id
		Format string //商品规格
	}
	st := &Para{}
	if e := s.GetPara(st); e != nil {
		JsLogger.Error(e.Error())
		s.Forward("1", e.Error(), nil)
		return
	}
	if st.UID == "" {
		info := "RemoveFromShoppingCart : UID is empty\n"
		JsLogger.Error(info)
		s.Forward("2", info, nil)
		return
	}
	if st.ProID == "" {
		info := "RemoveFromShoppingCart : ProID is empty\n"
		JsLogger.Error(info)
		s.Forward("2", info, nil)
		return
	}
	data := ShoppingCart{}
	JsRedis.Redis_hget(constant.SHOPPINGCART, st.UID, &data)
	index := -1
	for i, v := range data.Data {
		str := fmt.Sprintf("ProFormat.Format =%s ,Format =%s", v.ProFormat.Format, st.Format)
		JsLogger.Error(str)
		if v.ProID == st.ProID && v.ProFormat.Format == st.Format {
			index = i
			break
		}
	}
	if index != -1 {
		info := "index != -1aaaaaaaaaaaaaaaaaaaaa"
		JsLogger.Error(info)
		data.Data = append(data.Data[:index], data.Data[index+1:]...)
	}
	if err := JsRedis.Redis_hset(constant.SHOPPINGCART, st.UID, &data); err != nil {
		info := "set error, try again" + err.Error()
		JsLogger.Error(info)
		s.Forward("5", info, nil)
		return
	}
	s.Forward("0", "success", nil)
}

//修改产品的数量
func ModifyGoodNumWithShoppingCart(s *JsHttp.Session) {
	type Para struct {
		UID    string //用户id
		ProID  string //产品id
		Num    int    //数量
		Format string //商品规格
	}
	st := &Para{}
	if e := s.GetPara(st); e != nil {
		JsLogger.Error(e.Error())
		s.Forward("1", e.Error(), nil)
		return
	}
	//逻辑正确性
	if st.UID == "" || st.ProID == "" || st.Num < 1 {
		info := fmt.Sprintf("ModifyGoodNumWithShoppingCart failed:UID=%s,ProID=%s,Num=%d\n", st.UID, st.ProID, st.Num)
		JsLogger.Error(info)
		s.Forward("2", info, nil)
		return
	}
	data := ShoppingCart{}
	if e := JsRedis.Redis_hget(constant.SHOPPINGCART, st.UID, &data); e != nil {
		JsLogger.Error(e.Error())
		s.Forward("3", e.Error(), nil)
	}
	index := -1
	for i, v := range data.Data {
		if v.ProID == st.ProID && v.ProFormat.Format == st.Format {
			index = i
			break
		}
	}
	if index != -1 {
		data.Data[index].Nums = st.Num
	}
	if err := JsRedis.Redis_hset(constant.SHOPPINGCART, st.UID, &data); err != nil {
		info := err.Error()
		JsLogger.Error(info)
		s.Forward("5", info, nil)
		return
	}
	s.Forward("0", "success", data)
}

//清空购物车
func ClearUserShoppingCart(s *JsHttp.Session) {
	type Para struct {
		UID string //用户id
	}
	st := &Para{}
	if e := s.GetPara(st); e != nil {
		JsLogger.Error(e.Error())
		s.Forward("1", e.Error(), nil)
		return
	}
	if st.UID == "" {
		JsLogger.Error("ClearUserShoppingCart failed,UID i empty\n")
		s.Forward("2", "ClearUserShoppingCart failed,UID i empty\n", nil)
		return
	}
	data := ShoppingCart{}
	if e := JsRedis.Redis_hset(constant.SHOPPINGCART, st.UID, &data); e != nil {
		JsLogger.Error(e.Error())
		s.Forward("4", e.Error(), nil)
		return
	}
	s.Forward("0", "ClearUserShoppingCart success", data)
}

//从购物车删除多个商品（从购物车支付后删除购物商品）
func RemoveMoreFromShoppingCart(s *JsHttp.Session) {
	type id_form struct {
		ProID  string //产品id
		Format string //商品规格
	}
	type Para struct {
		UID    string    //用户id
		IDform []id_form //产品ID和规格
	}

	st := &Para{}
	if e := s.GetPara(st); e != nil {
		JsLogger.Error(e.Error())
		s.Forward("1", e.Error(), nil)
		return
	}
	if st.UID == "" {
		JsLogger.Error("ClearUserShoppingCart failed,UID i empty\n")
		s.Forward("2", "ClearUserShoppingCart failed,UID i empty\n", nil)
		return
	}
	for _, v := range st.IDform {
		if v.ProID == "" || v.Format == "" {
			info := "RemoveFromShoppingCart : ProID or Format = nil."
			JsLogger.Error(info)
			s.Forward("2", info, nil)
			return
		}
	}
	data := ShoppingCart{}
	JsRedis.Redis_hget(constant.SHOPPINGCART, st.UID, &data)
	l := len(data.Data)
	newdata := make([]Goods, l)
	for i, v := range data.Data {
		for _, w := range st.IDform {
			if v.ProID == w.ProID && v.ProFormat.Format == w.Format {
			} else {
				newdata = append(newdata, data.Data[i])
			}
		}
	}
	data.Data = newdata
	//index := -1
	//for _, w := range st.IDform {
	//	for i, v := range data.Data {
	//		if v.ProID == w.ProID && v.ProFormat.Format == w.Format {
	//			index = i
	//			break
	//		}
	//	}
	//	if index != -1 {
	//		data.Data = append(data.Data[:index], data.Data[index+1:]...)
	//	}
	//}
	if err := JsRedis.Redis_hset(constant.SHOPPINGCART, st.UID, &data); err != nil {
		info := "set error, try again" + err.Error()
		JsLogger.Error(info)
		s.Forward("5", info, nil)
		return
	}
	s.Forward("0", "success", nil)
}

//选择商品

//小计商品数量

//商品金额

//计算

//总价

//总计产品数量
