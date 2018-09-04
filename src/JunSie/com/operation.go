package com

import (
	"JsGo/JsBench/JsProduct"
	"JsGo/JsHttp"
	"JsGo/JsLogger"
	"JsGo/JsStore/JsRedis"
	"JunSie/constant"
	"JunSie/util"
	"fmt"
	"strconv"
	"time"
)

func Init_operation() {
	JsHttp.WhiteHttp("/getoperationnum", GetOperationNum) //获取运营版面数字
	JsHttp.WhiteHttp("/getprorepertory", GetProRepertory) //获取获取库存不足产品
	JsHttp.WhiteHttp("/gethotprolist", GetHotProList)     //获取热销排行榜
	JsHttp.WhiteHttp("/getsalestrend", GetSalesTrend)     //获取销售趋势sales trend
}

type operationNum struct {
	VisitNum     int //用户访问量（总访问量）
	WaitOrd      int //待处理订单数量
	WaitMesage   int //待处理消息数量
	ProRepertory int //库存紧张产品数量
}

//获取运营版面数字
func GetOperationNum(s *JsHttp.Session) {
	data := &operationNum{}
	data.VisitNum = GetVisitNum()
	type paratA struct {
		Num int
	}
	dataA := &paratA{}
	if err := JsRedis.Redis_hget(constant.H_Order, constant.OrderStatus_Paid, dataA); err == nil {
		data.WaitOrd = dataA.Num
	} else {
		data.WaitOrd = 0
	}
	data.WaitMesage = GetMessageNum() //待处理消息数量

	data.ProRepertory = len(GetProRepertoryf()) ////库存紧张产品数量

	s.Forward("0", "success", data)
}

//商家管理后台 运营面板需要的统计数据：
//1.待处理订单数量     一、
//2.待处理消息数量     二、OK
//3.库存紧张产品   (支付之后，库存减少同时判断库存数量如果库存在小于库存值之后，就添加到库存紧张列表中)
//4.热门购买产品（前五或前十）
//5.用户访问量（总访问量）点击次数，访问人数，
//6.订单周销售趋势或则月销售趋势

//11111.进入待处理列表状态
//（待处理表放订单ID数组，和订单ID的数量）
//用户下单商家未发货                        支付完成
//进入列表条件
// 进入已支付列表 ，

type WaitSt struct {
	WaitNum int      //等待处理数量
	OrdID   []string //订单ID
}

//33333.库存紧张
//减少库存条件 //用户支付订单
type WaitStProductFormat struct {
	ProID     string //产品ID
	Title     string //标题
	Format    string //规格说明
	Pic       string //规格对应的图片
	Price     int    //规格对应的价格
	Inventory int    //规格库存
}

//添加到库存不足表中
func RepertoryLadd(ProID, Format, Title, Pic string, Price, Inventory, Nums int) { //RepertoryListAdd 添加到库存表中
	para := []WaitStProductFormat{}
	JsRedis.Redis_get(constant.RepertoryLess, &para)
	exist := false
	for _, v := range para {
		if v.ProID == ProID && v.Format == Format { //判读是否存在
			v.Inventory = Inventory - Nums
			//订单中的库存数量是没有购买成功的数量，如果购买成功数量需要减去，才回与当前库存数量一致。
			exist = true
			break
		}
	}
	if !exist { //不存在
		para = append(para, WaitStProductFormat{
			ProID:     ProID,
			Title:     Title,
			Format:    Format,
			Pic:       Pic,
			Price:     Price,
			Inventory: Inventory - Nums,
		})
	}

	JsRedis.Redis_set(constant.RepertoryLess, &para)
	return
}

func RepertoryLRemove(ProID, Format string) { //RepertoryListAdd 添加到库存表中
	para := []WaitStProductFormat{}
	JsRedis.Redis_get(constant.RepertoryLess, &para)
	i := 0
	for site, v := range para {
		if v.ProID == ProID && v.Format == Format { //判读是否存在
			i = site
			para = append(para[:i], para[i+1:]...)
			JsRedis.Redis_set(constant.RepertoryLess, &para)
			break
		}
	}
	return
}

//获取库存不足产品规格列表
func GetProRepertory(s *JsHttp.Session) {
	//获取库存数量是动态更新的
	dataA := GetProRepertoryf()

	s.Forward("0", "success", dataA)
}
func GetProRepertoryf() (v []WaitStProductFormat) {
	//获取库存数量是动态更新的
	type Format struct {
		ProFormat []JsProduct.ProductFormat //产品规格
	}
	dataT := Format{}
	data := []WaitStProductFormat{}
	//wait :=WaitStProductFormat{}
	dataA := []WaitStProductFormat{}
	if err := JsRedis.Redis_get(constant.RepertoryLess, &data); err != nil {
		JsLogger.Error(err.Error())
		return nil
	}
	for _, v := range data {
		JsRedis.Redis_hget(constant.DB_Product, v.ProID, &dataT)
		for _, vb := range dataT.ProFormat { //产品规格中的库存数量是最新的
			if v.Format == vb.Format {
				v.Inventory = vb.Inventory
				break
			}
		}
		if v.Inventory < Inventory_Lower {
			dataA = append(dataA, v)
		}
	}

	return dataA

}

//获取库存不足数量
func GetProRepertoryNum() (num int) {
	data := []WaitStProductFormat{}
	if err := JsRedis.Redis_get(constant.RepertoryLess, &data); err != nil {
		JsLogger.Error(err.Error())
		num = 0
	} else {
		num = len(data)
	}
	return num
}

//增加库存//商家修改库存数量
//修改库存数量同样会调接口检查库存不足表中是否存在，存在的话就需要删除，
//设置库存为多少时是库存紧张
//默认值设置表，

//4444444444.热门购买产品（前五或前十）        map根据v值排序
// 做排行榜(一个热销榜单只保存前n个) K是所需要排的产品，  V是产品购买量，根据v的值给K排序。
//销售量 添加到产品统计信息中 。。。。。。。。。。。。。。。。。。。。。。。。。。。。。。。。（支付后）
//排行是所有，
//第一种，时时更新，
//那个一个产品，
//增加销售量SalesVolume

//结构
type Sales struct {
	ProID string //产品ID
	Num   int    //销售数量（各种规格的总数）
}

//最大的放在最前面
//更新热销排行榜如果有新的
func UpdateHotList(ProID string) { //
	type Sales struct {
		ProID string //产品ID
		Num   int    //销售数量（各种规格的总数）
	}
	type parao struct {
		SalesVolume int //销量（各种规格的总数）
	}
	Patao := &parao{}
	if err := JsRedis.Redis_hget(constant.ProStatistics, ProID, Patao); err != nil {
		JsLogger.Info(err.Error())
		return
	}
	Max := 10 //排行榜数量
	Para := make([]Sales, 0, Max+1)
	JsRedis.Redis_get(constant.HotProList, &Para)
	if len(Para) >= Max && Patao.SalesVolume < Para[Max].Num && Para[Max].ProID != ProID {
		return //如果不在排行榜范围内（排行榜长度最大，并且比最后一个还小，最后一个不是它）
	}
	Data := make([]Sales, 0, len(Para)+1)

	sa := Sales{ProID: ProID, Num: Patao.SalesVolume}
	for i, v := range Para {
		if v.ProID == ProID {
			Para = append(Para[:i], Para[i+1:]...) //如果存在取出来删除生成一个新队列
		} //取出后添加的都变成排行榜中没有的
	}
	exist := false
	for j, w := range Para {
		if Patao.SalesVolume >= w.Num { //如果新的产品销售数量大于排行
			exist = true
			temp := append([]Sales{}, Para[j:]...) //向中间追加
			Para = append(Para[:j], sa)
			Data = append(Para, temp...)
			//插入到产品所在排行位置
			break
		}
	}
	if !exist {
		Data = append(Data, sa)
	}
	if len(Data) > Max {
		Data = Data[:Max]
	}
	JsRedis.Redis_set(constant.HotProList, &Data)
}

//获取产品热销排行榜
func GetHotProList(s *JsHttp.Session) {

	type statistics struct {
		VisitNum       int     // 访问量
		PraiseNum      int     // 点赞数量
		AttentionNums  int     // 关注数量
		CompositeScore float64 // 综合评分（所有评分的平均值）
		CommentNum     int     // 评论的人数
		SalesVolume    int     //销量
	}
	type ProListAbs struct { //abstract摘要
		ID        string     //产品ID
		Title     string     //标题
		Thumbnail string     //缩略图
		Images    []string   //海报
		Tongji    statistics //统计
	}
	hotpro := ProListAbs{}
	Hotpros := []ProListAbs{}
	data := []Sales{}
	if err := JsRedis.Redis_get(constant.HotProList, &data); err != nil {
		JsLogger.Error(err.Error())
		s.Forward("3", err.Error(), nil)
		return
	}
	for _, v := range data {
		err := JsRedis.Redis_hget(constant.DB_Product, v.ProID, &hotpro)
		if err != nil {
			JsLogger.Error(err.Error())
			continue
		}
		err1 := JsRedis.Redis_hget(constant.ProStatistics, v.ProID, &hotpro.Tongji)
		if err1 != nil {
			JsLogger.Error(err1.Error())
			continue
		}
		Hotpros = append(Hotpros, hotpro)
	}

	s.Forward("0", "success", &Hotpros)

}

//5.用户访问量
//用户访问量（总访问量）点击次数，访问人数，

type PageView struct {
	PV int // "PageView" 页面浏览量或点击量
	//UV int //"Unique_visitor"独立访客数(在同一天内访问者A再次访问该网站则不计数)
}

//存一个表保存当日浏览用户的UID
//UIDs []string //用户ID切片
func NewVisit() { //增加用户访问量
	para := &PageView{}
	JsRedis.Redis_get(constant.VisitNumber, para)
	para.PV++
	JsRedis.Redis_set(constant.VisitNumber, para)
}

//获取总访问量
func GetVisitNum() (PV int) {
	para := &PageView{}
	err := JsRedis.Redis_get(constant.VisitNumber, para)
	if err == nil {
		data := para.PV
		return data
	} else {
		JsLogger.Error(err.Error())
		data := 0
		return data
	}

}

//666666666.订单周销售趋势或则月销售趋势

//获取销售趋势sales trend
func GetSalesTrend(s *JsHttp.Session) { //
	type Para struct {
		StartDate string //开始时间格式2006-01-02
		EndDate   string //结束时间格式2006-01-02
	}
	para := &Para{}
	if err := s.GetPara(para); err != nil {
		JsLogger.Error(err.Error())
		s.Forward("1", err.Error(), nil)
		return
	}
	if para.StartDate == "" || para.EndDate == "" {
		str := fmt.Sprintf("GetSalesTrend failed,StartDate=%s,EndDate=%s\n", para.StartDate, para.EndDate)
		JsLogger.Error(str)
		s.Forward("2", str, nil)
		return
	}
	//时间转换
	startTime, err1 := time.ParseInLocation("2006-01-02", para.StartDate, time.Local)
	endTime, err2 := time.ParseInLocation("2006-01-02", para.EndDate, time.Local)
	if err1 != nil {
		JsLogger.Error(err1.Error())
		s.Forward("2", err1.Error(), nil)
		return
	}
	if err2 != nil {
		JsLogger.Error(err2.Error())
		s.Forward("2", err2.Error(), nil)
		return
	}
	timeNowUnix := time.Now().Unix()
	startTimeUnix := startTime.Unix()
	endTimeUnix := endTime.Unix()

	if endTimeUnix < startTimeUnix || startTimeUnix > timeNowUnix || endTimeUnix > timeNowUnix {
		str := fmt.Sprintf("GetSalesTrend failed,StartDate=%s,EndDate=%s\n", para.StartDate, para.EndDate)
		JsLogger.Error(str)
		s.Forward("2", str, nil)
		return
	}

	curtime := startTime //循环时间
	type Info struct {
		Date string
		Sale int
	}
	DATA := []Info{}
	lastTime := curtime //记录上一天的时间
	data := make(map[string]int)
	for {
		if curtime.Unix() > endTimeUnix || curtime.Unix() > timeNowUnix {
			fmt.Println(endTime.Format("2006-01-02"), time.Now().Format("2006-01-02"), curtime.Format("2006-01-02"))
			break
		}
		date := curtime.Format("2006-01-02")
		key := strconv.Itoa(curtime.Year())

		if lastTime.Year() < curtime.Year() || lastTime.Day() == curtime.Day() { //如果是第一次获取，循环时间等于上次循环时间

			fmt.Print(lastTime.Year(), curtime.Year(), lastTime.Day(), curtime.Day())
			//如果是上次循环的年比本次循环的年小，说明跨年了，key发生了变化。
			if err := JsRedis.Redis_hget(constant.SalesTrend, key, &data); err != nil {
				JsLogger.Warn(err.Error())
			}
		}
		k, ok := data[date] //获取k值，并判断是否存在map的K值。
		info := Info{}
		info.Date = date
		if ok {
			info.Sale = k
		} else {
			info.Sale = 0
		}
		DATA = append(DATA, info)
		lastTime = curtime
		curtime = curtime.AddDate(0, 0, 1)
	}
	s.Forward("0", "success", DATA)
}

func updataOrderNum() { //保存更新当日订单数量
	key := strconv.Itoa(time.Now().Year())
	data := make(map[string]int)
	if err := JsRedis.Redis_hget(constant.SalesTrend, key, &data); err != nil {
		JsLogger.Error(err.Error())
	}
	cur_date := util.CurDate()
	if k, ok := data[cur_date]; ok { //判断K是否存在
		data[cur_date] = k + 1
	} else {
		data[cur_date] = 1
	}
	JsRedis.Redis_hset(constant.SalesTrend, key, &data)
	return
}
