package com

import (
	"JsGo/JsBench/JsProduct"
	"JsGo/JsHttp"
	"JsGo/JsLogger"
	"JsGo/JsStore/JsRedis"

	"JunSie/util"
	"fmt"

	"JunSie/constant"

	"sync"

	"github.com/chanxuehong/wechat/json"
)

type PriceRange struct {
	EnableRange bool //启用价格区间
	MinPrice    int  //最小价格
	MaxPrice    int  //最大价格
}

//产品结构
type XM_Product struct {
	JsProduct.Product //通用产品
	PriceRange        //价格区间
}

var Inventory_Lower int = 100 //库存下限
var lowerMutex sync.Mutex     //锁

func loadInventory() {
	lowerMutex.Lock()
	Inventory_Lower = getInventory_Lower()
	lowerMutex.Unlock()
}

func Init_product() {
	JsHttp.WhiteHttps("/newproduct", NewProduct)                     //新建产品
	JsHttp.WhiteHttps("/modproduct", ModifyProduct)                  //修改产品信息
	JsHttp.WhiteHttps("/modifyproductformat", ModifyProductFormat)   //修改产品规格
	JsHttp.WhiteHttps("/delproductdb", DelProductDB)                 //数据库永久删除产品
	JsHttp.WhiteHttps("/delproductmark", DelProductMark)             //标记删除产品
	JsHttp.WhiteHttps("/updownproduct", Updownproduct)               //上下架产品
	JsHttp.WhiteHttps("/newproductdraft", NewProductDraft)           //新建产品草稿
	JsHttp.WhiteHttps("/removeproductdraft", RemoveProductDraft)     //删除产品草稿
	JsHttp.WhiteHttps("/getproductdraft", GetProductDraft)           //获取产品草稿
	JsHttp.WhiteHttps("/removeproducttag", RemoveProductTag)         //移除产品的某一个标签
	JsHttp.WhiteHttp("/queryproduct", QueryProduct)                  //查询单个产品
	JsHttp.WhiteHttp("/querymoreproducts", QueryMoreProducts)        //查询多个产品
	JsHttp.WhiteHttp("/getpageproducts", GetPageProducts)            //获取分页产品
	JsHttp.WhiteHttp("/getproductnums", GetProductNums)              //获取不同状态产品的数量
	JsHttp.WhiteHttp("/getproducttags", GetProductTags)              //获取产品的所有标签
	JsHttp.WhiteHttp("/getproductstatics", GetProductStatics)        //获取产品统计
	JsHttp.WhiteHttps("/newproductvisit", NewProVisit)               //增加产品访问
	JsHttp.WhiteHttps("/newproductpraise", NewProPraise)             //增加产品点赞
	JsHttp.WhiteHttps("/newproductattention", NewProAttention)       //增加产品关注
	JsHttp.WhiteHttps("/cancelproductpraise", RemoveProPraise)       //产品取消点赞
	JsHttp.WhiteHttps("/cancelproductattention", RemoveProAttention) //产品取消关注
	JsHttp.WhiteHttps("/cancelproductvisit", RemoveProVisit)         //删除产品访问（+）

	JsHttp.WhiteHttp("/getinventorylower", GetInventoryLower) //获取库存下限
	JsHttp.WhiteHttp("/setinventorylower", SetInventoryLower) //设置库存下限
	loadInventory()
}
func Init_productMall() {
	JsHttp.WhiteHttps("/modproduct", ModifyProduct)                  //修改产品信息
	JsHttp.WhiteHttps("/modifyproductformat", ModifyProductFormat)   //修改产品规格
	JsHttp.WhiteHttps("/delproductdb", DelProductDB)                 //数据库永久删除产品
	JsHttp.WhiteHttps("/delproductmark", DelProductMark)             //标记删除产品
	JsHttp.WhiteHttps("/updownproduct", Updownproduct)               //上下架产品
	JsHttp.WhiteHttps("/newproductdraft", NewProductDraft)           //新建产品草稿
	JsHttp.WhiteHttps("/removeproductdraft", RemoveProductDraft)     //删除产品草稿
	JsHttp.WhiteHttps("/getproductdraft", GetProductDraft)           //获取产品草稿
	JsHttp.WhiteHttps("/removeproducttag", RemoveProductTag)         //移除产品的某一个标签
	JsHttp.WhiteHttps("/queryproduct", QueryProduct)                 //查询单个产品
	JsHttp.WhiteHttps("/querymoreproducts", QueryMoreProducts)       //查询多个产品
	JsHttp.WhiteHttps("/getpageproducts", GetPageProducts)           //获取分页产品
	JsHttp.WhiteHttps("/getproductnums", GetProductNums)             //获取不同状态产品的数量
	JsHttp.WhiteHttps("/getproducttags", GetProductTags)             //获取产品的所有标签
	JsHttp.WhiteHttps("/getproductstatics", GetProductStatics)       //获取产品统计
	JsHttp.WhiteHttps("/newproductvisit", NewProVisit)               //增加产品访问
	JsHttp.WhiteHttps("/newproductpraise", NewProPraise)             //增加产品点赞
	JsHttp.WhiteHttps("/newproductattention", NewProAttention)       //增加产品关注
	JsHttp.WhiteHttps("/cancelproductpraise", RemoveProPraise)       //产品取消点赞
	JsHttp.WhiteHttps("/cancelproductattention", RemoveProAttention) //产品取消关注
	JsHttp.WhiteHttps("/cancelproductvisit", RemoveProVisit)         //删除产品访问（+）

	JsHttp.WhiteHttp("/getinventorylower", GetInventoryLower) //获取库存下限
	JsHttp.WhiteHttp("/setinventorylower", SetInventoryLower) //设置库存下限
	loadInventory()
}

//新建产品
func NewProduct(s *JsHttp.Session) {
	type Para struct {
		Pro  XM_Product //产品结构
		Tags []string   //标签列表
	}
	st := &Para{}
	if err := s.GetPara(st); err != nil {
		info := "newProduct:" + err.Error()
		JsLogger.Error(info)
		s.Forward("1", info, nil)
		return
	}
	if err := checkPro(st.Pro); err != nil {
		s.Forward("1", err.Error(), nil)
		return
	}
	if st.Pro.Status == "" {
		st.Pro.Status = constant.Status_ON
	}
	st.Pro.ID = util.IDer(constant.DB_Product)
	st.Pro.CreatDate = util.CurTime()

	if st.Pro.UID == "" {
		info := "NewProduct,UID is empty\n"
		JsLogger.Error(info)
		s.Forward("1", info, nil)
		return
	}

	if err := JsRedis.Redis_hset(constant.DB_Product, st.Pro.ID, st.Pro); err != nil {
		info := "newProduct,Redis_hset:" + err.Error()
		JsLogger.Error(info)
		s.Forward("1", info, nil)
		return
	}

	//建立产品搜索索引

	go creatProSearchIndex(&st.Pro)

	go appendToGlobalProducts(st.Pro.ID)
	if len(st.Tags) > 0 {
		go pro2Tag(st.Pro.ID, st.Tags)
	}
	TagLinkP(st.Pro.ID, st.Tags)
	// if erro := tage.TagLinkP(st.Pro.ID, st.Tags); erro != nil {
	// 	info := "Tage contect fail:" + erro.Error()
	// 	JsLogger.Error(info)
	// 	s.Forward("1", info, nil)
	// 	return
	// }
	s.Forward("0", "success", st)
}

//查询产品
func QueryProduct(s *JsHttp.Session) {
	type Info struct {
		ID string //产品id
	}
	st := &Info{}
	if err := s.GetPara(st); err != nil {
		info := "queryProduct,GetPara:" + err.Error()
		JsLogger.Error(info)
		s.Forward("1", info, nil)
		return
	}
	if st.ID == "" {
		info := "queryProduct param ID is empty\n"
		JsLogger.Error(info)
		s.Forward("1", info, nil)
		return
	}
	data := &XM_Product{}
	data, err := GetProductInfo(st.ID)
	if err != nil {
		info := "queryProduct Redis_get:" + err.Error()
		JsLogger.Error(info)
		s.Forward("1", info, data)
		return
	}
	s.Forward("0", "success\n", data)
}

//查询多个产品（Abs）
func QueryMoreProducts(s *JsHttp.Session) {
	type Info struct {
		IDs []string ////产品id列表
	}
	st := &Info{}
	if err := s.GetPara(st); err != nil {
		info := "queryMoreProducts,GetPara:" + err.Error()
		JsLogger.Error(info)
		s.Forward("1", info, nil)
		return
	}
	if len(st.IDs) < 0 {
		info := "queryMoreProducts param IDs is empty\n"
		JsLogger.Error(info)
		s.Forward("1", info, nil)
		return
	}
	s.Forward("0", "success", GetMoreProductInfo(st.IDs))
}

func ModifyProductFormat(s *JsHttp.Session) {
	type Para struct {
		ProID     string                    //产品创建者
		ProFormat []JsProduct.ProductFormat //产品规格
	}
	st := &Para{}
	if err := s.GetPara(st); err != nil {
		info := "modifyProduct,GetPara:" + err.Error()
		JsLogger.Error(info)
		s.Forward("1", info, nil)
		return
	}
	if st.ProID == "" {
		JsLogger.Error("ModifyProductFormat ProID is empty\n")
		s.Forward("1", "ModifyProductFormat ProID is empty\n", nil)
		return
	}
	if len(st.ProFormat) > 0 {
		for _, v := range st.ProFormat {
			if v.Format == "" || v.Price <= 0 {
				info := fmt.Sprintf("ModifyProductFormat,Format =%s , Price =%d:", v.Format, v.Price)
				JsLogger.Error(info)
				s.Forward("1", info, nil)
				return
			}
		}
	}

	data := &XM_Product{}
	if err := JsRedis.Redis_hget(constant.DB_Product, st.ProID, data); err != nil {
		info := "ModifyProductFormat Redis_get:" + err.Error()
		JsLogger.Error(info)
		s.Forward("1", info, data)
		return
	}
	data.ProFormat = st.ProFormat
	if err := JsRedis.Redis_hset(constant.DB_Product, st.ProID, data); err != nil {
		info := err.Error()
		JsLogger.Error(info)
		s.Forward("1", info, nil)
		return
	}
	s.Forward("0", "ModifyProductFormat success \n", data)
}

//修改产品
func ModifyProduct(s *JsHttp.Session) {
	type Para struct {
		Pro    XM_Product //产品结构
		Tags   []string   //标签列表
		IsTags bool       //是否修改标签
	}
	st := &Para{}
	if err := s.GetPara(st); err != nil {
		info := "modifyProduct,GetPara:" + err.Error()
		JsLogger.Error(info)
		s.Forward("1", info, nil)
		return
	}
	data := &XM_Product{}
	if err := JsRedis.Redis_hget(constant.DB_Product, st.Pro.ID, data); err != nil {
		info := "modifyProduct Redis_get:" + err.Error()
		JsLogger.Error(info)
		s.Forward("1", info, data)
		return
	}
	d, e := json.Marshal(st.Pro)
	if e != nil {
		info := e.Error()
		JsLogger.Error(info)
		s.Forward("1", info, nil)
		return
	}
	if err := json.Unmarshal(d, data); err != nil {
		info := err.Error()
		JsLogger.Error(info)
		s.Forward("1", info, nil)
		return
	}
	if err := JsRedis.Redis_hset(constant.DB_Product, data.ID, data); err != nil {
		info := err.Error()
		JsLogger.Error(info)
		s.Forward("1", info, nil)
		return
	}
	if st.IsTags {
		go changeProTag(data.ID, st.Tags)
	}
	s.Forward("0", "success", data)
}

//数据库删除产品
func DelProductDB(s *JsHttp.Session) {
	type Info struct {
		ID string //产品id
	}
	st := &Info{}
	if err := s.GetPara(st); err != nil {
		info := "delProductDB,GetPara:" + err.Error()
		JsLogger.Error(info)
		s.Forward("1", info, nil)
		return
	}
	if st.ID == "" {
		info := "delProductDB param ID is empty\n"
		JsLogger.Error(info)
		s.Forward("1", info, nil)
		return
	}
	if ok, e := JsRedis.Redis_hexists(constant.DB_Product, st.ID); e == nil && ok {
		if err := JsRedis.Redis_hdel(constant.DB_Product, st.ID); err != nil {
			info := "delProductDB" + err.Error()
			JsLogger.Error(info)
			s.Forward("1", info, nil)
			return
		}
	}
	go delFromGlobalProducts(st.ID)
	s.Forward("0", "success", nil)
}

//标记删除
func DelProductMark(s *JsHttp.Session) {
	type Info struct {
		ID string
	}
	st := &Info{}
	if err := s.GetPara(st); err != nil {
		info := "delProductMark,GetPara:" + err.Error()
		JsLogger.Error(info)
		s.Forward("1", info, nil)
		return
	}
	if st.ID == "" {
		info := "delProductMark param ID is empty\n"
		JsLogger.Error(info)
		s.Forward("1", info, nil)
		return
	}
	data := &XM_Product{}
	if err := JsRedis.Redis_hget(constant.DB_Product, st.ID, data); err != nil {
		info := "delProductMark,Redis_hset:" + err.Error()
		JsLogger.Error(info)
		s.Forward("1", info, nil)
		return
	}
	data.DelTag = true
	if err := JsRedis.Redis_hset(constant.DB_Product, st.ID, data); err != nil {
		info := "delProductMark,Redis_hset:" + err.Error()
		JsLogger.Error(info)
		s.Forward("1", info, nil)
		return
	}
	go delFromGlobalProducts(st.ID)
	s.Forward("0", "success", data)
}

//上下架产品
func Updownproduct(s *JsHttp.Session) {
	type Info struct {
		ID   string //id
		IsUp bool   //是否上架
	}
	st := &Info{}
	if err := s.GetPara(st); err != nil {
		info := "updownproduct,GetPara:" + err.Error()
		JsLogger.Error(info)
		s.Forward("1", info, nil)
		return
	}
	if st.ID == "" {
		info := "updownproduct param ID is empty\n"
		JsLogger.Error(info)
		s.Forward("1", info, nil)
		return
	}
	data := &XM_Product{}
	if err := JsRedis.Redis_hget(constant.DB_Product, st.ID, data); err != nil {
		info := "updownproduct,Redis_get:" + err.Error()
		JsLogger.Error(info)
		s.Forward("1", info, nil)
		return
	}
	if st.IsUp {
		data.Status = constant.Status_ON
	} else {
		data.Status = constant.Status_OFF
	}
	if err := JsRedis.Redis_hset(constant.DB_Product, st.ID, data); err != nil {
		info := "updownproduct,Redis_get:" + err.Error()
		JsLogger.Error(info)
		s.Forward("1", info, nil)
		return
	}
	s.Forward("0", "success", data)
}

//获取多个产品详情（摘要）
func GetMoreProductInfo(ids []string) []*JsProduct.ProductAbs {
	data := []*JsProduct.ProductAbs{}
	//fmt.Printf("id = %v\n", ids)
	for _, v := range ids {
		retData := &JsProduct.ProductAbs{}
		err := JsRedis.Redis_hget(constant.DB_Product, v, retData)
		if err != nil {
			JsLogger.Error(err.Error())
			continue
		}

		data = append(data, retData)
	}
	return data
}

//获取单个产品详情
func GetProductInfo(id string) (*XM_Product, error) {
	data := &XM_Product{}
	err := JsRedis.Redis_hget(constant.DB_Product, id, data)
	return data, err
}

func ReSaveProduct(pro *XM_Product) error {
	return JsRedis.Redis_hset(constant.DB_Product, pro.ID, pro)
}

func checkPro(pro XM_Product) error {

	if len(pro.ProFormat) > 0 {
		for _, v := range pro.ProFormat {
			if v.Price <= 0 || v.Format == "" || v.Inventory < 0 {
				return JsLogger.ErrorLog("Product ProFormat failed \n")
			}
		}
	}

	if len(pro.Images) < 0 {
		return JsLogger.ErrorLog("Product Images is empty \n")
	}
	if pro.EnableRange {
		if pro.MinPrice < 0 || pro.MaxPrice < 0 || pro.MaxPrice <= pro.MinPrice {
			return JsLogger.ErrorLog("Product Price Range error,MinPrice=%d,MaxPrice=%d\n", pro.MinPrice, pro.MaxPrice)
		}
	} else {
		if pro.OriPrice <= 0 || pro.NowPrice <= 0 || pro.NowPrice > pro.OriPrice {
			return JsLogger.ErrorLog("Product price error,OriPrice=%d,NowPrice=%d\n", pro.OriPrice, pro.NowPrice)
		}
	}
	return nil
}

//获取库存下限
func GetInventoryLower(s *JsHttp.Session) {
	loadInventory()
	s.Forward("0", "GetInventoryLower success\n", Inventory_Lower)
}

//设置库存下限
func SetInventoryLower(s *JsHttp.Session) {
	type Para struct {
		Inventory int //库存下限
	}
	st := &Para{}
	if err := s.GetPara(st); err != nil {
		JsLogger.Error(err.Error())
		s.Forward("1", err.Error(), nil)
		return
	}
	if st.Inventory <= 0 {
		str := fmt.Sprintf("SetInventoryLower failed,Inventory=%d\n", st.Inventory)
		JsLogger.Error(str)
		s.Forward("1", str, nil)
		return
	}
	if err := setInventory_Lower(st.Inventory); err != nil {
		JsLogger.Error(err.Error())
		s.Forward("1", err.Error(), nil)
		return
	}
	lowerMutex.Lock()
	Inventory_Lower = st.Inventory
	lowerMutex.Unlock()
	s.Forward("0", "success\n", st.Inventory)
}

//获取库存下限
func getInventory_Lower() int {
	Inventory := 100
	if err := JsRedis.Redis_get("InventoryLower", &Inventory); err != nil {
		go setInventory_Lower(Inventory)
	}
	return Inventory
}

//设置库存下限
func setInventory_Lower(Inventory int) error {
	return JsRedis.Redis_set("InventoryLower", &Inventory)
}
