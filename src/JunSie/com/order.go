package com

import (
	"JsGo/JsBench/JsUser"
	"JsGo/JsConfig"
	"JsGo/JsHttp"
	"JsGo/JsLogger"
	"JsGo/JsMobile"
	"JsGo/JsStore/JsRedis"
	"JunSie/constant"
	"JunSie/util"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"strconv"
	"strings"
	"time"
	"sync"
)

type St_PayCb struct {
	AppId                string `xml:"appid"`
	Mch_id               string `xml:"mch_id"`
	Device_info          string `xml:"device_info"`
	Nonce_str            string `xml:"nonce_str"`
	Sign                 string `xml:"sign"`
	Sign_type            string `xml:"sign_type"`
	Result_code          string `xml:"result_code"`
	Err_code             string `xml:"err_code"`
	Err_code_des         string `xml:"err_code_des"`
	Openid               string `xml:"openid"`
	Is_subscribe         string `xml:"is_subscribe"`
	Trade_type           string `xml:"trade_type"`
	Bank_type            string `xml:"bank_type"`
	Total_fee            string `xml:"total_fee"`
	Settlement_total_fee string `xml:"settlement_total_fee"`
	Fee_type             string `xml:"fee_type"`
	Cash_fee             string `xml:"cash_fee"`
	Cash_fee_type        string `xml:"cash_fee_type"`
	Transaction_id       string `xml:"transaction_id"`
	Out_trade_no         string `xml:"out_trade_no"`
	Attach               string `xml:"attach"`
	Time_end             string `xml:"time_end"`
}

type PayFlow struct {
	OpreatTime string //操作时间
	Status     string //操作之后的状态
	TimeStamp  int64  //操作时间戳
	OpreatID   string //操作人id
	OpreatName string //操作人名字
}

type PayInfo struct { //0
	TerminalIp   string //创建者ip
	Trade_type   string //支付方式。公众号还是app
	ProjectId    string
	PayWay       string            //0:微信 1:支付宝
	Freight      int               //运费
	OrderChannel string            //订单渠道，"web"、"app"、"minip"
	ThirdOrderID string            //第三方订单号
	RefundId     string            //退款ID
	TransferId   string            //打款id
	RefundFee    int               //退款金额
	Charge       map[string]string //票据
	PayCb        *St_PayCb         //支付
	RefundCb     map[string]string //退款回调
	Amount       int               //金钱
	AppId        string
	Mch_id       string
	Nonce_str    string
}

type OrderUserinfo struct { //1
	UID      string         //用户id
	OpenID   string         //微信opendID
	UnionID  string         //微信unionID
	UserName string         //用户姓名
	UserHead string         //用户头像
	Addr     JsUser.RecAddr //收货地址
}

type Order struct { //2
	OrderID       string            //订单id
	OrderUserinfo                   //订单的用户信息
	PayInfo                         //支付信息
	Product       []Goods           //商品
	Desc          string            //订单描述
	HasExpress    bool              //是否物流
	Express       *util.ExpressInfo //物流信息
	Flow          []*PayFlow        //支付流程
	Current       *PayFlow          //当前流程
	ExData        map[string]string //扩展的数据（备用）
	CreatTime     string            //创建时间
}

type OrderABS struct { //2
	OrderID  string //订单id
	UserName string //用户姓名
	//UserHead string         //用户头像
	Addr       JsUser.RecAddr    //收货地址
	Product    []Goods           //商品
	Desc       string            //订单描述
	HasExpress bool              //是否物流
	Express    *util.ExpressInfo //物流信息
	Current    *PayFlow          //当前流程
	PayWay     string            //0:微信 1:支付宝
	Amount     int               //金钱
	CreatTime  string            //创建时间
}

func (this *Order) appendFlow(ID, Name, Status string) {
	flow := &PayFlow{
		OpreatTime: util.CurTime(),
		Status:     Status,
		TimeStamp:  util.CurStamp(),
		OpreatID:   ID,
		OpreatName: Name,
	}
	this.Current = flow
	this.Flow = append(this.Flow, flow) //这个可以去掉
}

var pubPayHandler *PayHandler = nil
var smallPayHandler *PayHandler = nil
var appPayHandler *PayHandler = nil
var createMutex sync.Mutex     //锁
func InitPay() {
	pubPayHandler = NewPayHandler("PubPay") //handler 管理者，处理者
	if pubPayHandler == nil {
		log.Fatalln("Order pubPayHandler NewPayHandler failed!")
		return
	}
	smallPayHandler = NewPayHandler("WxMiniP") //handler 管理者，处理者
	if smallPayHandler == nil {
		log.Fatalln("Order smallPayHandler NewPayHandler failed!")
		return
	}
	appPayHandler = NewPayHandler("AppPay") //handler 管理者，处理者
	if appPayHandler == nil {
		log.Fatalln("Order appPayHandler NewPayHandler failed!")
		return
	}
}

func InitOrder() {
	JsHttp.WhiteHttp("/getglobalorder", GetGlobalOrder)          //获取分页的全局订单
	JsHttp.WhiteHttp("/getglobalordernums", GetGlobalOrderPages) //获取全局订单的不同状态的页数
	JsHttp.WhiteHttp("/queryorderinfo", QueryOrder)              //查询单个订单信息
	JsHttp.WhiteHttp("/querymoreorder", QueryMoreOrder)          //查询多个订单信息
	JsHttp.WhiteHttp("/getuserorderlist", GetUserOrderList)      //获取用户订单列表
	//JsHttp.WhiteHttp("/submitorder", SubmitOrder)                //提交订单到第三方支付
	//JsHttp.WhiteHttp("/paidordersuccesscb", PaidOrderSuccessCb)  //支付成功回调
	JsHttp.WhiteHttp("/ordersending", OrderSending)           //订单发货
	JsHttp.WhiteHttp("/modifyorderaddr", ModifyOrderAddr)     //修改订单收货地址
	JsHttp.WhiteHttp("/queryorderexpress", QueryOrderExpress) //查询订单的快递信息
	JsHttp.WhiteHttp("/queryexpress", queryExpress)           //查询快递信息
	JsHttp.WhiteHttp("/getoderlist", GetOderList)             //获取分页的全局订单GetOderListABS

}
func InitOrderMall() {
	//JsHttp.WhiteHttps("/getglobalorder", GetGlobalOrder)          //获取分页的全局订单
	//JsHttp.WhiteHttps("/getglobalordernums", GetGlobalOrderPages) //获取全局订单的不同状态的页数
	JsHttp.WhiteHttps("/queryorderinfo", QueryOrder)                           //查询单个订单信息
	JsHttp.WhiteHttps("/querymoreorder", QueryMoreOrder)                       //查询多个订单信息
	JsHttp.WhiteHttps("/getuserorderlist", GetUserOrderList)                   //获取用户订单列表
	JsHttp.WhiteHttps("/creatorder", CreatOrder)                               //创建订单
	JsHttp.WhiteHttps("/cancelorder", CancelOrder)                             //取消订单
	JsHttp.WhiteHttps("/submitorder", SubmitOrder)                             //提交订单到第三方支付
	JsHttp.WhiteHttps("/paidordersuccesscb", PaidOrderSuccessCb)               //支付成功回调
	JsHttp.WhiteHttps("/apprrderpaidsuccessnotify", AppOrderPaidSuccessNotify) //app支付成功后通知
	JsHttp.WhiteHttps("/ordersending", OrderSending)                           //订单发货
	JsHttp.WhiteHttps("/orderreception", OrderReception)                       //用户主动收货
	JsHttp.WhiteHttps("/modifyorderaddr", ModifyOrderAddr)                     //修改订单收货地址
	JsHttp.WhiteHttps("/queryorderexpress", QueryOrderExpress)                 //查询订单的快递信息
	JsHttp.WhiteHttps("/queryexpress", queryExpress)                           //查询快递信息
	JsHttp.WhiteHttps("/ordersuccess", OrderSuccess)                           //用户完成评价

}

//获取分页的全局订单
func GetGlobalOrder(session *JsHttp.Session) {
	type Para struct {
		Status string //状态(为空时,表示不区分状态)
		SIndex int    //启始索引
		Size   int    //个数
	}
	st := &Para{}
	if err := session.GetPara(st); err != nil {
		info := "getPageProducts,GetPara:" + err.Error()
		JsLogger.Error(info)
		session.Forward("1", info, nil)
		return
	}
	if !checkStatus(st.Status) || st.SIndex < 0 || st.Size <= 0 {
		info := fmt.Sprintf("GetGlobalOrder param error,Status=%s,SIndex=%d,Size=%d\n", st.Status, st.SIndex, st.Size)
		JsLogger.Error(info)
		session.Forward("1", info, nil)
		return
	}
	data := []*Order{}
	list, err := JsRedis.Redis_hkeys(constant.H_Order)
	if err != nil {
		session.Forward("0", "GetGlobalOrder: 获取全局订单号失败\n", data)
		return
	}
	for _, v := range list {
		d, err := getOrder(v)
		if err == nil {
			if st.Status == "" || st.Status == d.Current.Status {
				data = append(data, d)
			}
		}
	}

	if st.SIndex < len(data) {
		if st.Size < len(data) && (st.SIndex+st.Size) < len(data) {
			session.Forward("0", "GetGlobalOrder success", data[st.SIndex:st.SIndex+st.Size])
			return
		}
		session.Forward("0", "GetGlobalOrder success", data[st.SIndex:])
		return
	}
	session.Forward("0", "GetGlobalOrder: 没有满足条件的订单", data)
}

//获取分页的全局订单
func GetGlobalOrderPages(s *JsHttp.Session) {
	type para struct {
		Status string //状态(为空时,表示不区分状态)
	}
	st := &para{}
	if e := s.GetPara(st); e != nil {
		JsLogger.Error(e.Error())
		s.Forward("1", e.Error(), nil)
		return
	}

	type paratA struct {
		Num int
	}
	data := paratA{}

	switch st.Status {
	case "OrderStatus_Paid":
		JsRedis.Redis_hget(constant.H_Order, constant.OrderStatus_Paid, &data)
	case "OrderStatus_Send":
		JsRedis.Redis_hget(constant.H_Order, constant.OrderStatus_Send, &data)
	case "OrderStatus_Receive":
		JsRedis.Redis_hget(constant.H_Order, constant.OrderStatus_Receive, &data)
	default:
		info := "Status error"
		JsLogger.Error(info)
		s.Forward("2", info, nil)
		return
	}
	s.Forward("0", "success", data.Num)
}

//type Para struct {
//	Status string //状态(为空时,表示不区分状态)
//}
//st := &Para{}
//if err := session.GetPara(st); err != nil {
//	info := "GetGlobalOrderPages,GetPara:" + err.Error()
//	JsLogger.Error(info)
//	session.Forward("1", info, nil)
//	return
//}
//if !checkStatus(st.Status) {
//	info := fmt.Sprintf("GetGlobalOrderPages param error,Status=%s\n", st.Status)
//	JsLogger.Error(info)
//	session.Forward("1", info, 0)
//	return
//}
//data := []*Order{}
//list, err := JsRedis.Redis_hkeys(constant.H_Order)
//if err != nil {
//	session.Forward("0", "GetGlobalOrderPages: 获取全局订单号页数失败\n", 0)
//	return
//}
//for _, v := range list {
//	d, err := getOrder(v)
//	if err == nil {
//		if st.Status == "" || st.Status == d.Current.Status {
//			data = append(data, d)
//		}
//	}
//}
//session.Forward("0", "GetGlobalOrderPages: success\n", len(data))
//}

//查询单个订单信息
func QueryOrder(session *JsHttp.Session) {
	type Para struct {
		OrderID string //订单id
	}
	st := &Para{}
	if err := session.GetPara(st); err != nil {
		JsLogger.Error(err.Error())
		session.Forward("1", err.Error(), nil)
		return
	}
	if st.OrderID == "" {
		JsLogger.Error("QueryOrder Failed,OrderID is empty\n")
		session.Forward("1", "QueryOrder Failed,OrderID is empty\n", nil)
		return
	}
	data, err := getOrder(st.OrderID)
	if err != nil {
		JsLogger.Error(err.Error())
		session.Forward("1", err.Error(), nil)
		return
	}
	session.Forward("0", "QueryOrder SUCCESS\n", data)
}

//查询多个订单信息
func QueryMoreOrder(session *JsHttp.Session) {
	type Para struct {
		IDs []string //订单id列表
	}
	st := &Para{}
	if err := session.GetPara(st); err != nil {
		JsLogger.Error(err.Error())
		session.Forward("1", err.Error(), nil)
		return
	}
	data := []*Order{}
	for _, v := range st.IDs {
		d, err := getOrder(v)
		if err == nil {
			data = append(data, d)
		}
	}
	session.Forward("0", "QueryMoreOrder success\n", data)
}

//获取用户订单列表
func GetUserOrderList(session *JsHttp.Session) {
	type Para struct {
		UID string //用户id
	}
	st := &Para{}
	if err := session.GetPara(st); err != nil {
		JsLogger.Error(err.Error())
		session.Forward("1", err.Error(), nil)
		return
	}
	if st.UID == "" {
		JsLogger.Error("GetUserOrderList failed,UID is empty\n")
		session.Forward("1", "GetUserOrderList failed,UID is empty\n", nil)
		return
	}
	list := []string{}
	if err := JsRedis.Redis_hget(constant.H_UserOrder, st.UID, &list); err != nil {
		JsLogger.Error(err.Error())
		session.Forward("1", err.Error(), nil)
		return
	}
	/*用户端需要展示
	  1、待付款订单（提交订单后，提示付款，未付款进入待付款）（已创建，已提交）
	2、待收货（已支付，已发货，用户确认收货，查看订单物流状态）
	3、待评价（已收货，待评价）
	4、已完成订单()(????????????????????????)
	*/
	type Statuslist struct {
		WaitPay      []*Order //待支付
		WaitReceive  []*Order //待收货
		WaitEvaluate []*Order //待评价
		Success      []*Order //订单完成
	}

	dataList := &Statuslist{}
	goodList := list //存放正常订单
	index := -1
	for i, v := range list {
		d, err := getOrder(v)
		if err != nil {
			_, errt := getOrder(v)
			//两次判断
			if err == errt {
				index = i
				goodList = append(goodList[:index], goodList[index+1:]...)
			}
			continue
		}
		//session.Forward("3", err.Error(), nil)
		//return //如果有找不到的会导致一个都拉不到。
		//判断状态，存到对应的状态列表中，
		if d.Current.Status == constant.OrderStatus_Creat || d.Current.Status == constant.OrderStatus_Submit {
			if len(d.Flow) > 0 {
				if d.Flow[0].Status == constant.OrderStatus_Creat && time.Now().Unix()-d.Flow[0].TimeStamp > 3600 {
					//	delList = append(delList, v)
					continue
				}
			}
			dataList.WaitPay = append(dataList.WaitPay, d)
		}
		if d.Current.Status == constant.OrderStatus_Paid {
			dataList.WaitReceive = append(dataList.WaitReceive, d)
			continue
		}
		if d.Current.Status == constant.OrderStatus_Send {
			dataList.WaitReceive = append(dataList.WaitReceive, d)
			continue
		}

		if d.Current.Status == constant.OrderStatus_Receive {
			dataList.WaitEvaluate = append(dataList.WaitEvaluate, d)
			continue
		}
		if d.Current.Status == constant.OrderStatus_Success {
			dataList.Success = append(dataList.Success, d)
			continue
		}
	}

	go JsRedis.Redis_hset(constant.H_UserOrder, st.UID, &goodList)
	session.Forward("0", "GetUserOrderList success\n", dataList)
}

//创建订单
func CreatOrder(session *JsHttp.Session) {
createMutex.Lock()
	st := &Order{}
	fmt.Print("oodd1111111111", st)
	if err := session.GetPara(st); err != nil {
		JsLogger.Error(err.Error())
		session.Forward("1", err.Error(), nil)
		return
	}
	if st.UID == "" || len(st.Product) == 0 {
		str := fmt.Sprintf("CreatOrder failed: UID =%s ,Len(Product) =%d", st.UID, len(st.Product))
		JsLogger.Error(str)
		session.Forward("1", str, nil)
		return
	}
	if st.UserName == "" || st.UserHead == "" {
		user := JsUser.User{}
		if e := JsRedis.Redis_hget(constant.USER, st.UID, &user); e == nil {
			st.UserName = user.Nickname
			st.UserHead = user.Header
			st.OpenID = user.Openid
			st.UnionID = user.Unionid
		}
	}

	////////////??????????????检查用户信息？？？？？？？？？？？、、、、//////////

	/////////////////查询更新库存//////////////////////////////////////////
	if err := updateInventory(st, false); err != nil {
		session.Forward("1", err.Error(), nil)
		return
	}
	st.OrderID = util.IDer(constant.H_Order)
	st.OrderID = st.OrderID[:32]
	st.CreatTime = util.CurTime()
	st.appendFlow(st.UID, st.UserName, constant.OrderStatus_Creat)
	st.Amount = 0
	for _, v := range st.Product {
		st.Amount += v.ProFormat.Price * v.Nums
	}
	st.Freight = 0
	if st.Addr.Province != "" {
		st.Freight = queryExpressPrice(st.Addr.Province)
	}
	st.Amount += st.Freight

	payDebug, e := JsConfig.GetConfigString([]string{"PayDebug"})

	if e != nil {

	} else if payDebug == "true" {
		JsLogger.Info("Pay Debug is set! All price is ￥0.01")
		st.Amount = 1
	}

	//st.Amount = 1
	if err := JsRedis.Redis_hset(constant.H_Order, st.OrderID, st); err != nil {
		JsLogger.Error(err.Error())
		session.Forward("1", err.Error(), nil)
		go updateInventory(st, true)
		return
	}
	//添加用户订单
	go appendUserOrder(st.UID, st.OrderID)
	go appendWaitPayOrder(st)
	session.Forward("0", "CreatOrder success\n", st)
	createMutex.Unlock()
}

//取消订单
func CancelOrder(session *JsHttp.Session) {
	type Para struct {
		OrderID string //订单id
		UID     string //用户UID
	}
	st := &Para{}
	if err := session.GetPara(st); err != nil {
		JsLogger.Error(err.Error())
		session.Forward("1", err.Error(), nil)
		return
	}
	if st.OrderID == "" || st.UID == "" {
		str := fmt.Sprintf("CancelOrder failed: OrderID = %s ,UID = %s", st.OrderID, st.UID)
		JsLogger.Error(str)
		session.Forward("1", str, nil)
		return
	}
	data := &Order{}
	if err := JsRedis.Redis_hget(constant.H_Order, st.OrderID, data); err != nil {
		JsLogger.Error(err.Error())
		session.Forward("1", err.Error(), nil)
		return
	}
	if data.UID != st.UID {
		str := fmt.Sprintf("CancelOrder failed: OrderID(%s) != UID(%s)", st.OrderID, st.UID)
		JsLogger.Error(str)
		session.Forward("1", str, nil)
		return
	}
	if data.Current.Status == constant.OrderStatus_Success {
		str := fmt.Sprintf("CancelOrder failed: Current Status is %s", data.Current.Status)
		JsLogger.Error(str)
		session.Forward("1", str, nil)
		return
	}

	if data.Amount >= 100 && (data.Current.Status == constant.OrderStatus_Paid ||
		data.Current.Status == constant.OrderStatus_Send ||
		data.Current.Status == constant.OrderStatus_Receive) {
		if data.PayWay == "" {
			cb := make(map[string]string)
			var payErr error = nil
			if data.OrderChannel == "web" {
				data.Trade_type = "JSAPI"
				cb, payErr = pubPayHandler.Wx_refund(data)
			} else if data.OrderChannel == "app" {
				data.Trade_type = "APP"
				cb, payErr = appPayHandler.Wx_refund(data)
			} else if data.OrderChannel == "minip" {
				data.Trade_type = "JSAPI"
				cb, payErr = smallPayHandler.Wx_refund(data)
			} else {

			}
			if payErr != nil {
				JsLogger.Error(payErr.Error())
				session.Forward("1", payErr.Error(), nil)
				return
			}
			JsLogger.Error("result_code=%v\n", cb)
			if cb["result_code"] == "FAIL" {
				str := fmt.Sprintf("微信退款失败,return_msg:%s,err_code_des:%s\n", cb["return_msg"], cb["err_code_des"])
				JsLogger.Error(str)
				session.Forward("1", str, nil)
				return
			}
			data.RefundCb = cb
		}

		data.appendFlow(data.UID, data.UserName, constant.OrderStatus_Cancel)
	} else {
		data.appendFlow(data.UID, data.UserName, constant.OrderStatus_Cancel)
	}
	/////////////////恢复库存///////////////////////////
	go updateInventory(data, true)

	if err := JsRedis.Redis_hset(constant.H_Order, st.OrderID, data); err != nil {
		JsLogger.Error(err.Error())
		session.Forward("1", err.Error(), nil)
		return
	}
	//从待支付的列表中删除
	go removeWaitPayOrder(st.OrderID)
	session.Forward("0", "CancelOrder success\n", data)
}

//提交订单到第三方支付
func SubmitOrder(session *JsHttp.Session) {
	type Para struct {
		OrderID      string //订单id
		UID          string //用户UID
		OrderChannel string //订单渠道，"web"、"app"、"minip"
		PayWay       string //0:微信 1:支付宝
	}
	st := &Para{}
	if err := session.GetPara(st); err != nil {
		JsLogger.Error(err.Error())
		session.Forward("1", err.Error(), nil)
		return
	}
	if st.OrderID == "" || st.UID == "" ||
		(st.OrderChannel == "" || (st.OrderChannel != "web" && st.OrderChannel != "app" && st.OrderChannel != "minip")) ||
		(st.PayWay != "0" && st.PayWay != "1") {
		str := fmt.Sprintf("SubmitOrder failed: OrderID = %s ,UID = %s,OrderChannel=%s,PayWay =%s",
			st.OrderID, st.UID, st.OrderChannel, st.PayWay)
		JsLogger.Error(str)
		session.Forward("1", str, nil)
		return
	}
	data := &Order{}
	if err := JsRedis.Redis_hget(constant.H_Order, st.OrderID, data); err != nil {
		JsLogger.Error(err.Error())
		session.Forward("1", err.Error(), nil)
		return
	}

	addr := session.Req.RemoteAddr
	i := strings.Index(addr, ":")
	data.TerminalIp = addr[:i]

	data.PayWay = st.PayWay
	data.OrderChannel = st.OrderChannel

	ch := make(map[string]string)
	var payErr error = nil
	if data.PayWay == "0" { //微信支付
		user := &JsUser.User{}
		if err := JsRedis.Redis_hget(constant.USER, st.UID, user); err != nil {
			JsLogger.Error(err.Error())
			session.Forward("1", err.Error(), nil)
			return
		}
		if st.OrderChannel == "web" {
			data.OpenID = user.Openid
		}
		if st.OrderChannel == "app" {
			data.OpenID = user.Openid_app
		}
		if st.OrderChannel == "minip" {
			data.OpenID = user.Openid_small
		}

		if st.OrderChannel == "web" {
			data.Trade_type = "JSAPI"
			ch, payErr = pubPayHandler.Wx_pay(data)
		} else if st.OrderChannel == "app" {
			data.Trade_type = "APP"
			ch, payErr = appPayHandler.Wx_pay(data)
		} else if st.OrderChannel == "minip" {
			data.Trade_type = "JSAPI"
			ch, payErr = smallPayHandler.Wx_pay(data)
		} else {

		}
		if payErr != nil {
			JsLogger.Error(payErr.Error())
			session.Forward("4", payErr.Error(), nil)
			return
		}
		data.Charge = ch
	} else { //支付宝支付

	}
	data.appendFlow(st.OrderID, data.UID, constant.OrderStatus_Submit)
	if err := JsRedis.Redis_hset(constant.H_Order, st.OrderID, data); err != nil {
		JsLogger.Error(err.Error())
		session.Forward("1", err.Error(), nil)
		return
	}
	session.Forward("0", "SubmitOrder success\n", data)
}

//app端支付回调
func AppOrderPaidSuccessNotify(session *JsHttp.Session) {
	type Para struct {
		OrderID string
	}
	st := &Para{}
	if err := session.GetPara(st); err != nil {
		session.Forward("1", err.Error(), nil)
		return
	}
	if err := orderPaid(st.OrderID, nil); err != nil {
		session.Forward("1", err.Error(), nil)
		return
	}

	session.Forward("0", "success!\n", nil)
}

//支付成功回调
func PaidOrderSuccessCb(session *JsHttp.Session) {
	body, e := ioutil.ReadAll(session.Req.Body)
	if e != nil {
		return
	}
	paycb := &St_PayCb{}
	err := xml.Unmarshal(body, paycb)
	xml := ""
	if err != nil {
		xml = `<xml>
 				<return_code><![CDATA[FAIL]]></return_code>
 				<return_msg><![CDATA[` + err.Error() + `]]></return_msg>
			   </xml>`
	} else {
		if err := orderPaid(paycb.Out_trade_no, paycb); err != nil {
			JsLogger.ErrorLog(err.Error())
		}
		xml = `<xml>
 				<return_code><![CDATA[SUCCESS]]></return_code>
 				<return_msg><![CDATA[OK]]></return_msg>
			</xml>`
	}
	session.WriteString(xml)
}

func orderPaid(OrderID string, Cb *St_PayCb) error {
	JsLogger.Error("Enter orderPaid ...................................................\n")
	if OrderID == "" {
		return JsLogger.ErrorLog("orderPaid faild,OrderID=%s\n", OrderID)
	}
	if Cb == nil {
		JsLogger.Info("orderPaid,St_PayCb is nil\n")
	}
	data := &Order{}
	if err := JsRedis.Redis_hget(constant.H_Order, OrderID, data); err != nil {
		return JsLogger.ErrorLog(err.Error())
	}

	if data.Current.Status != constant.OrderStatus_Submit {
		str := fmt.Sprintf("orderPaid failed: 当前订单%s状态为%s,不能进行支付", OrderID, data.Current.Status)
		return JsLogger.ErrorLog(str)
	}
	if Cb != nil {
		data.PayCb = Cb //微信支付回调
		data.ThirdOrderID = Cb.Out_trade_no
		m, e := strconv.Atoi(Cb.Cash_fee)
		if e == nil {
			data.RefundFee = m ///支付价格
		}
	} else {
		data.PayCb = nil
		data.RefundFee = data.Amount
	}
	data.appendFlow(OrderID, data.UID, constant.OrderStatus_Paid)
	if err := JsRedis.Redis_hset(constant.H_Order, OrderID, data); err != nil {
		return JsLogger.ErrorLog(err.Error())
	}

	//更新单个产品销量(库存不足表)
	for _, v := range data.Product {
		if v.ProFormat.Inventory < Inventory_Lower {
			go RepertoryLadd(v.ProID, v.ProFormat.Format, v.ProName, v.ProFormat.Pic, v.ProFormat.Price, v.ProFormat.Inventory, v.Nums)
		} else {
			go RepertoryLRemove(v.ProID, v.ProFormat.Format)
		}
		//热销排行榜
		go NewProductSaleNum(v.ProID, v.Nums)

	}
	//更新当日订单数
	go updataOrderNum()
	go AddOrderSend(OrderID)
	///发送信息
	sms_order_paid(data)
	//从待支付的列表中删除
	go removeWaitPayOrder(OrderID)
	JsLogger.Error("Leave orderPaid ...........................................\n")
	return nil
}

func sms_order_paid(order *Order) {

	type SMS struct {
		SignName string
		SmsCode  string
		Mobile   string
	}
	sms := &SMS{}
	JsRedis.Redis_hget(constant.ADMIN, "SMS", sms)

	para := make(map[string]string)
	//新订单通知，订单号：${order}，订单日期：${date}，产品：${product}，价格：${price}，寄送地址：${address}，收货人：${name}，联系方式：${cell}
	para["order"] = order.OrderID[:8]
	para["date"] = order.Current.OpreatTime
	if len(order.Product) == 1 {
		str := fmt.Sprintf("%s[%s]", order.Product[0].ProName, order.Product[0].ProFormat.Format)
		ary := []rune{}
		for _, v := range str {
			if v == '】' || v == '【' {

			} else {
				ary = append(ary, v)
			}
		}
		para["product"] = string(ary)
	} else if len(order.Product) > 1 {
		str := ""
		for i, v := range order.Product {
			str += fmt.Sprintf("产品%d：%s【%s】；", i, v.ProName, v.ProFormat.Format)
		}
		para["product"] = str
	}
	f := (float64)(order.Amount / 100.0)

	para["price"] = fmt.Sprintf("%.2f", f)
	provice := ""
	if order.Addr.Province == "上海市" || order.Addr.Province == "北京市" || order.Addr.Province == "天津市" || order.Addr.Province == "重庆市" {

	} else {
		provice = order.Addr.Province
	}
	para["address"] = provice + order.Addr.City + order.Addr.Area + order.Addr.Addr
	para["name"] = order.Addr.Name
	para["cell"] = order.Addr.Cell
	para["rand"] = strconv.Itoa(rand.Int() % 1000000)

	JsMobile.ComJsMobileVerify(sms.SignName, sms.Mobile, sms.SmsCode, "c", 300, para)
}

//订单发货
func OrderSending(session *JsHttp.Session) {
	type Para struct {
		OrderID       string //订单id
		OpreatID      string //操作者ID
		HasExpress    bool   //是否物流
		ExpressName   string //快递名称
		ExpressNumber string //快递单号
	}
	st := &Para{}
	if err := session.GetPara(st); err != nil {
		JsLogger.Error(err.Error())
		session.Forward("1", err.Error(), nil)
		return
	}
	if st.OrderID == "" || st.OpreatID == "" {
		str := fmt.Sprintf("OrderSending failed: OrderID = %s ,UID = %s", st.OrderID, st.OpreatID)
		JsLogger.Error(str)
		session.Forward("1", str, nil)
		return
	}
	if st.HasExpress {
		if st.ExpressName == "" || st.ExpressNumber == "" {
			str := fmt.Sprintf("OrderSending failed: ExpressName = %s ,ExpressNumber = %s", st.ExpressName, st.ExpressNumber)
			JsLogger.Error(str)
			session.Forward("1", str, nil)
			return
		}
	}

	data := &Order{}
	if err := JsRedis.Redis_hget(constant.H_Order, st.OrderID, data); err != nil {
		JsLogger.Error(err.Error())
		session.Forward("1", err.Error(), nil)
		return
	}
	if data.Current.Status != constant.OrderStatus_Paid {
		str := fmt.Sprintf("OrderSending failed: 当前订单%s状态为%s,不能进行发货", st.OrderID, data.Current.Status)
		JsLogger.Error(str)
		session.Forward("1", str, nil)
		return
	}
	/////////////////添加物流信息//////////////////////////////
	data.HasExpress = st.HasExpress
	data.Express = &util.ExpressInfo{
		ExpressNumber: st.ExpressNumber,
		ExpressName:   st.ExpressName,
	}

	data.appendFlow(st.OrderID, "", constant.OrderStatus_Send)

	if err := JsRedis.Redis_hset(constant.H_Order, st.OrderID, data); err != nil {
		JsLogger.Error(err.Error())
		session.Forward("1", err.Error(), nil)
		return
	}
	go AddOrderReceive(st.OrderID)
	session.Forward("0", "OrderSending success\n", data)
}

//查询订单的物流信息
func QueryOrderExpress(session *JsHttp.Session) {
	type Para struct {
		OrderID string //订单id
	}
	st := &Para{}
	if err := session.GetPara(st); err != nil {
		JsLogger.Error(err.Error())
		session.Forward("1", err.Error(), nil)
		return
	}
	if st.OrderID == "" {
		str := fmt.Sprintf("QueryOrderExpress failed: OrderID = %s", st.OrderID)
		JsLogger.Error(str)
		session.Forward("1", str, nil)
		return
	}
	data := &Order{}
	if err := JsRedis.Redis_hget(constant.H_Order, st.OrderID, data); err != nil {
		JsLogger.Error(err.Error())
		session.Forward("1", err.Error(), nil)
		return
	}
	ex := []util.ExpressFlow{
		{
			Time:    util.CurTime(),
			Ftime:   util.CurTime(),
			Context: "暂无物流信息",
		},
	}
	if data.Current.Status != constant.OrderStatus_Send &&
		data.Current.Status != constant.OrderStatus_Receive &&
		data.Current.Status != constant.OrderStatus_Success {
		str := fmt.Sprintf("QueryOrderExpress failed: OrderID = %sStatus=%s,订单还未发货 ", st.OrderID, data.Current.Status)
		JsLogger.Error(str)
		session.Forward("0", "暂无物流信息", ex)
		return
	}

	express, err := util.QueryExpress(data.Express.ExpressName, data.Express.ExpressNumber)
	if err != nil {
		session.Forward("0", "暂无物流信息", ex)
		return
	}
	//express, err := util.QueryExpress("yunda", "3940237440749")
	//if err != nil {
	//	session.Forward("0", "暂无物流信息", ex)
	//	return
	//}
	session.Forward("0", "query success\n", express)
}

func queryExpress(session *JsHttp.Session) {
	type Para struct {
		ExpressName   string //物流公司名称队名的编号
		ExpressNumber string //快递单号
	}
	st := &Para{}
	if err := session.GetPara(st); err != nil {
		JsLogger.Error(err.Error())
		session.Forward("1", err.Error(), nil)
		return
	}
	if st.ExpressNumber == "" || st.ExpressName == "" {
		str := fmt.Sprintf("queryExpress failed:ExpressNumber=%s,ExpressName=%s\n", st.ExpressNumber, st.ExpressName)
		JsLogger.Error(str)
		session.Forward("1", str, nil)
		return
	}
	express, err := util.QueryExpress(st.ExpressName, st.ExpressNumber)
	if err != nil {
		session.Forward("1", "暂无物流信息", nil)
		return
	}
	session.Forward("0", "success", express)
}

//用户主动收货
func OrderReception(session *JsHttp.Session) {
	type Para struct {
		OrderID string //订单id
		UID     string //用户UID
	}
	st := &Para{}
	if err := session.GetPara(st); err != nil {
		JsLogger.Error(err.Error())
		session.Forward("1", err.Error(), nil)
		return
	}
	if st.OrderID == "" || st.UID == "" {
		str := fmt.Sprintf("OrderReception failed: OrderID = %s ,UID = %s", st.OrderID, st.UID)
		JsLogger.Error(str)
		session.Forward("1", str, nil)
		return
	}
	data := &Order{}
	if err := JsRedis.Redis_hget(constant.H_Order, st.OrderID, data); err != nil {
		JsLogger.Error(err.Error())
		session.Forward("1", err.Error(), nil)
		return
	}

	if data.UID != st.UID {
		str := fmt.Sprintf("OrderReception failed: 不是同一个用户 data.UID = %s ,st.UID = %s", data.UID, st.UID)
		JsLogger.Error(str)
		session.Forward("1", str, nil)
		return
	}
	if data.Current.Status != constant.OrderStatus_Send {
		str := fmt.Sprintf("OrderReception failed: 当前订单%s状态为%s,不能进确认收货", st.OrderID, data.Current.Status)
		JsLogger.Error(str)
		session.Forward("1", str, nil)
		return
	}
	data.appendFlow(st.OrderID, data.UserName, constant.OrderStatus_Receive)

	if err := JsRedis.Redis_hset(constant.H_Order, st.OrderID, data); err != nil {
		JsLogger.Error(err.Error())
		session.Forward("1", err.Error(), nil)
		return
	}
	go AddOrderSuccess(st.OrderID)
	session.Forward("0", "OrderReception success\n", data)
}

//用户完成评价
func OrderSuccess(session *JsHttp.Session) {
	type Para struct {
		OrderID string //订单id
		UID     string //用户UID
	}
	st := &Para{}
	if err := session.GetPara(st); err != nil {
		JsLogger.Error(err.Error())
		session.Forward("1", err.Error(), nil)
		return
	}
	if st.OrderID == "" || st.UID == "" {
		str := fmt.Sprintf("OrderReception failed: OrderID = %s ,UID = %s", st.OrderID, st.UID)
		JsLogger.Error(str)
		session.Forward("1", str, nil)
		return
	}
	data := &Order{}
	if err := JsRedis.Redis_hget(constant.H_Order, st.OrderID, data); err != nil {
		JsLogger.Error(err.Error())
		session.Forward("3", err.Error(), nil)
		return
	}

	if data.UID != st.UID {
		str := fmt.Sprintf("OrderReception failed: 不是同一个用户 data.UID = %s ,st.UID = %s", data.UID, st.UID)
		JsLogger.Error(str)
		session.Forward("2", str, nil)
		return
	}
	if data.Current.Status != constant.OrderStatus_Receive {
		str := fmt.Sprintf("OrderReception failed: 当前订单%s状态为%s,不能完成订单", st.OrderID, data.Current.Status)
		JsLogger.Error(str)
		session.Forward("2", str, nil)
		return
	}
	data.appendFlow(st.OrderID, data.UserName, constant.OrderStatus_Success)

	if err := JsRedis.Redis_hset(constant.H_Order, st.OrderID, data); err != nil {
		JsLogger.Error(err.Error())
		session.Forward("4", err.Error(), nil)
		return
	}
	session.Forward("0", "OrderSuccess success\n", data)
}

//修改订单收货地址
//在未发货前修改
func ModifyOrderAddr(session *JsHttp.Session) {
	type Para struct {
		OrderID string         //订单ID
		UID     string         //用户id
		Addr    JsUser.RecAddr //收货地址
	}
	st := &Para{}
	if err := session.GetPara(st); err != nil {
		JsLogger.Error(err.Error())
		session.Forward("1", err.Error(), nil)
		return
	}
	if st.OrderID == "" || st.UID == "" || st.Addr.Name == "" || st.Addr.Cell == "" || st.Addr.Addr == "" {
		str := fmt.Sprintf("ModifyOrderAddr failed: OrderID = %s ,UID = %s,Name=%s,Cell=%s,Addr=%s",
			st.OrderID, st.UID, st.Addr.Name, st.Addr.Cell, st.Addr.Addr)
		JsLogger.Error(str)
		session.Forward("1", str, nil)
		return
	}
	data := &Order{}
	if err := JsRedis.Redis_hget(constant.H_Order, st.OrderID, data); err != nil {
		JsLogger.Error(err.Error())
		session.Forward("1", err.Error(), nil)
		return
	}
	if data.UID != st.UID {
		str := fmt.Sprintf("ModifyOrderAddr failed: 不是同一个用户 data.UID = %s ,st.UID = %s", data.UID, st.UID)
		JsLogger.Error(str)
		session.Forward("1", str, nil)
		return
	}
	if data.Current.Status != constant.OrderStatus_Creat &&
		data.Current.Status != constant.OrderStatus_Submit &&
		data.Current.Status != constant.OrderStatus_Paid {
		str := fmt.Sprintf("ModifyOrderAddr failed: 当前状态%s不能修改地址", data.Current.Status)
		JsLogger.Error(str)
		session.Forward("1", str, nil)
		return
	}

	data.Addr = st.Addr
	if err := JsRedis.Redis_hset(constant.H_Order, st.OrderID, data); err != nil {
		JsLogger.Error(err.Error())
		session.Forward("1", err.Error(), nil)
		return
	}
	session.Forward("0", "ModifyOrderAddr success\n", data)
}

//获取单个订单信息
func getOrder(OrderID string) (*Order, error) {
	data := &Order{}
	err := JsRedis.Redis_hget(constant.H_Order, OrderID, data)
	return data, err
}

//获取不同状态的订单列表
func getStatusOrder(Status string) ([]*Order, error) {
	data := []*Order{}
	list, err := JsRedis.Redis_hkeys(constant.H_Order)
	if err != nil {
		return nil, err
	}
	for _, v := range list {
		d, err := getOrder(v)
		if err == nil {
			if Status == d.Current.Status {
				data = append(data, d)
			}
		}
	}
	return data, nil
}

//添加用户订单
func appendUserOrder(UID, OrderID string) error {
	list := []string{}
	JsRedis.Redis_hget(constant.H_UserOrder, UID, &list)
	exist := false
	for _, v := range list {
		if v == OrderID {
			exist = true
			break
		}
	}
	if !exist {
		list = append(list, OrderID)
	}
	return JsRedis.Redis_hset(constant.H_UserOrder, UID, &list)
}

//更新库存
func updateInventory(order *Order, isAdd bool) error {
	/////////////////恢复库存//////////////////////////////////////////
	if len(order.Product) > 0 {
		for _, v := range order.Product {
			pro, err := GetProductInfo(v.ProID)
			if err == nil {
				for i, f := range pro.ProFormat {
					if f.Format == v.ProFormat.Format {
						if isAdd {
							pro.ProFormat[i].Inventory += v.Nums
						} else {
							if f.Inventory < v.Nums {
								str := fmt.Sprintf("下单失败: ProID=%s ,Format=%s,Inventory=%d,OrderNums=%d,库存不足\n",
									v.ProID, f.Format, f.Inventory, v.Nums)
								return JsLogger.ErrorLog(str)
							}
							pro.ProFormat[i].Inventory -= v.Nums
						}
						go JsRedis.Redis_hset(constant.DB_Product, v.ProID, pro)
						break
					}
				}
			}
		}
	}
	return nil
}

//校验订单状态
func checkStatus(status string) bool {
	if status != "" &&
		status != constant.OrderStatus_Creat &&
		status != constant.OrderStatus_Submit &&
		status != constant.OrderStatus_Cancel &&
		status != constant.OrderStatus_Paid &&
		status != constant.OrderStatus_Send &&
		status != constant.OrderStatus_Receive &&
		status != constant.OrderStatus_Success {
		return false
	}
	return true
}
