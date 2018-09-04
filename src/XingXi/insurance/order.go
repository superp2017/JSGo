package main

import (
	"JsGo/JsHttp"
	"JsGo/JsLogger"
	"JsGo/JsStore/JsRedis"
	"JsGo/JsWeChat/JsWechatPay/jswxpay"
	"JunSie/util"
	"encoding/xml"
	"io/ioutil"
	"log"
	"strconv"
	"strings"
)

const (
	Opreat_Order_Submit     = "Opreat_Order_Submit"     //提交订单
	Opreat_Order_UserPaid   = "Opreat_Order_UserPaid"   //用户已支付
	Opreat_Order_UserCancle = "Opreat_Order_UserCancle" //用户取消  ??cancel 取消
	Opreat_Order_Refund     = "Opreat_Order_Refund"     //用户已经退款
	Insurance_Order         = "Insurance_Order"         //保单
)

var payHandler *jswxpay.PayHandler

func init() {
	payHandler = jswxpay.NewPayHandler("PubPay") //handler 管理者，处理者
	if payHandler == nil {
		log.Fatalln("NewPayHandler failed!")
		return
	}
}

//获取用户所有订单
func getUserOrderList(session *JsHttp.Session) {
	type para struct {
		UID string
	}
	st := &para{}
	if err := session.GetPara(st); err != nil {
		session.Forward("1", err.Error(), nil)
		return
	}
	list, _ := getUserOrder(st.UID)
	data := []*jswxpay.StOrder{}
	for _, v := range list {
		d, e := getOrder(v)
		if e == nil {
			data = append(data, d)
		}
	}
	session.Forward("0", "success", data)
}

//提交订单
func newOrder(session *JsHttp.Session) {
	para := &jswxpay.StOrder{}
	e := session.GetPara(para)
	if e != nil {
		JsLogger.Error(e.Error())
		session.Forward("2", e.Error(), nil)
		return
	}
	insurID, ok := para.ExData["InsuranceID"]
	if !ok || insurID == "" {
		session.Forward("1", "ExData[InsuranceID] is empty\n", nil)
		return
	}
	addr := session.Req.RemoteAddr
	i := strings.Index(addr, ":")
	para.TerminalIp = addr[:i]
	para.TimeStamp = util.CurStamp()
	para.OrderId = util.IDer(Insurance_Order)
	para.OrderId = para.OrderId[:32]
	para.Trade_type = "JSAPI"
	ch, e := payHandler.Wx_pay(para)
	if e != nil {
		JsLogger.Error(e.Error())
		session.Forward("4", e.Error(), nil)
		return
	}
	para.Charge = ch

	if err := submitOrder(para); err != nil {
		session.Forward("1", err.Error(), nil)
		return
	}
	session.Forward("0", "success", para)
}

//支付回调
func paySuccessCb(session *JsHttp.Session) {
	body, e := ioutil.ReadAll(session.Req.Body)
	if e != nil {
		return
	}
	paycb := &jswxpay.ST_PayCb{}
	err := xml.Unmarshal(body, paycb)
	xml := ""
	if err != nil {
		xml = `<xml>
 				<return_code><![CDATA[FAIL]]></return_code>
 				<return_msg><![CDATA[` + e.Error() + `]]></return_msg>
			   </xml>`
	} else {
		if _, err := orderUserPaid(paycb.Out_trade_no, paycb); err != nil {
			JsLogger.ErrorLog(err.Error())
		}
		xml = `<xml>
 				<return_code><![CDATA[SUCCESS]]></return_code>
 				<return_msg><![CDATA[OK]]></return_msg>
			</xml>`
	}
	session.WriteString(xml)
}

//提交订单
func submitOrder(order *jswxpay.StOrder) error {
	JsLogger.Info("Enter submint order")
	if order.Uid == "" || order.OrderId == "" {
		return JsLogger.ErrorLog("SubmitOrder param is empty,UID=%s,OrderID=%s,\n", order.Uid, order.OrderId)
	}
	order.Status = Opreat_Order_Submit
	order.CreateDate = util.CurTime()
	if err := JsRedis.Redis_hset(Insurance_Order, order.OrderId, order); err != nil {
		return nil
	}
	go appendUserOrder(order.Uid, order.OrderId)
	insur, ok := order.ExData["InsuranceID"]
	if ok {
		go UpdateInsuranceOrder(insur, order.OrderId, order.Status, order.TimeStamp)
	}

	return nil
}

//取消订单
func orderCancle(session *JsHttp.Session) {
	type Para struct {
		OrderID string
	}
	para := &Para{}
	if e := session.GetPara(para); e != nil {
		JsLogger.Error(e.Error())
		session.Forward("1", e.Error(), nil)
		return
	}
	data := &jswxpay.StOrder{}
	if err := JsRedis.Redis_hget(Insurance_Order, para.OrderID, data); err != nil {
		session.Forward("1", err.Error(), nil)
		return
	}

	if data.Status == Opreat_Order_UserPaid {
		if err := orderRefund(data); err != nil {
			session.Forward("1", err.Error(), nil)
			return
		}
	} else {
		if err := orderCancel(data); err != nil {
			session.Forward("1", err.Error(), nil)
			return
		}
	}

	insur, ok := data.ExData["InsuranceID"]
	if ok {
		go UpdateInsuranceOrder(insur, data.OrderId, data.Status, util.CurStamp())
	}
	session.Forward("0", "success", nil)
}

//支付成功
func orderUserPaid(OrderID string, cb *jswxpay.ST_PayCb) (*jswxpay.StOrder, error) {
	JsLogger.Error("Enter OrderUserPaid .............................")
	if OrderID == "" {
		return nil, JsLogger.ErrorLog("OrderUserPaid faild,OrderID=%s\n", OrderID)
	}
	if cb == nil {
		return nil, JsLogger.ErrorLog("OrderUserPaid,WxST_PayCb is nil\n")
	}
	data := &jswxpay.StOrder{}
	if err := JsRedis.Redis_hget(Insurance_Order, OrderID, data); err != nil {
		return nil, JsLogger.ErrorLog(err.Error())
	}
	data.PaidTimeStamp = util.CurStamp()
	data.PayCb = cb //微信支付回调
	m, e := strconv.Atoi(cb.Cash_fee)
	if e == nil {
		data.RefundFee = m ///支付价格
	}
	data.Status = Opreat_Order_UserPaid
	if err := JsRedis.Redis_hset(Insurance_Order, OrderID, data); err != nil {
		return nil, JsLogger.ErrorLog(err.Error())
	}

	insur, ok := data.ExData["InsuranceID"]
	if ok {
		go UpdateInsuranceOrder(insur, data.OrderId, data.Status, data.PaidTimeStamp)
	}
	JsLogger.Error("Leave OrderUserPaid .............................")
	return data, nil
}

//取消退款
func orderRefund(order *jswxpay.StOrder) error {
	cb, e := payHandler.Wx_refund(order)
	if e != nil {
		JsLogger.Error(e.Error())
		return e
	}
	JsLogger.Info("result_code=%v\n", cb)
	if cb["result_code"] == "FAIL" {
		return JsLogger.ErrorLog("微信退款失败,return_msg:%s,err_code_des:%s\n", cb["return_msg"], cb["err_code_des"])
	}
	order.RefundCb = cb
	order.Status = Opreat_Order_Refund
	return JsRedis.Redis_hset(Insurance_Order, order.OrderId, order)
}

//取消不退款
func orderCancel(order *jswxpay.StOrder) error {
	order.Status = Opreat_Order_UserCancle
	return JsRedis.Redis_hset(Insurance_Order, order.OrderId, order)
}

//查询单个订单、、
func getOrder(orderID string) (*jswxpay.StOrder, error) {
	data := &jswxpay.StOrder{}
	err := JsRedis.Redis_hget(Insurance_Order, orderID, data)
	return data, err
}

//添加订单到用户
func appendUserOrder(UID, ID string) error {
	return JsRedis.Redis_Sset(UID, ID)
}

//获取用户的订单id列表
func getUserOrder(uid string) ([]string, error) {
	data := []string{}
	d, err := JsRedis.Redis_Sget(uid)
	for _, v := range d {
		data = append(data, string(v.([]byte)))
	}
	return data, err
}

//获取单个订单信息（详情）
func getOrderInfo(session *JsHttp.Session) {
	type para struct {
		OrderID string
	}
	st := &para{}
	if err := session.GetPara(st); err != nil {
		session.Forward("1", err.Error(), nil)
		return
	}
	d, e := getOrder(st.OrderID)
	if e != nil {
		session.Forward("2", e.Error(), nil)
		return
	}
	session.Forward("0", "success", d)
}
