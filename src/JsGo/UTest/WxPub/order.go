package main

import (
	"JsGo/JsHttp"
	. "JsGo/JsLogger"
	"JsGo/JsWeChat/JsWechatPay/jswxpay"
	"encoding/base64"
	"encoding/xml"
	"io/ioutil"
	"strconv"
	"strings"
	"time"
)

const (
	base64Table = "123QRSTUabcdVWXYZHijKLAWDCABDstEFGuvwxyzGHIJklmnopqr234560178912"
)

var coder = base64.NewEncoding(base64Table)

var g_order map[string]*jswxpay.StOrder

var pubPayHandler *jswxpay.PayHandler = nil
var appPayHandler *jswxpay.PayHandler = nil

func InitOrder() {
	pubPayHandler = jswxpay.NewPayHandler("PubPay")
	appPayHandler = jswxpay.NewPayHandler("AppPay")

	JsHttp.Https("/pubneworder", pubNewOrder)
	JsHttp.Https("/pubpaysuccess", pubWxPayCb)
	JsHttp.Https("/pubrefund", pubRefund)

	JsHttp.Https("/appneworder", appNewOrder)
	JsHttp.Https("/apppaysuccess", appWxPayCb)
	JsHttp.Https("/apprefund", appRefund)

	g_order = make(map[string]*jswxpay.StOrder)
}

func pubNewOrder(session *JsHttp.Session) {
	type Para struct {
		Desc    string
		OrderId string
	}
	para := &Para{}

	e := session.GetPara(para)
	if e != nil {
		Error(e.Error())
		session.Forward("1", e.Error(), nil)
		return
	}

	addrs := strings.Split(session.Req.RemoteAddr, ":")

	order := &jswxpay.StOrder{}
	order.Desc = para.Desc
	order.OpenId = para.OrderId
	order.ProjectId = ""
	order.TerminalIp = addrs[0]
	order.Amount = 10
	order.OrderId = strconv.Itoa(time.Now().Nanosecond())
	order.Trade_type = "JSAPI"

	m, e := pubPayHandler.Wx_pay(order)
	if e != nil {
		Error(e.Error())
		session.Forward("2", e.Error(), nil)
		return
	}

	g_order[order.OrderId] = order
	m["OrderId"] = order.OrderId

	session.Forward("0", "success", m)

}

func pubWxPayCb(session *JsHttp.Session) {

	str := ""
	body, e := ioutil.ReadAll(session.Req.Body)
	if e != nil {
		str = `<xml>
  				<return_code><![CDATA[FAIL]]></return_code>
  				<return_msg><![CDATA[` + e.Error() + `]]></return_msg>
			   </xml>`
		session.WriteString(str)
		return
	}

	paycb := &jswxpay.ST_PayCb{}
	e = xml.Unmarshal(body, paycb)

	if e != nil {
		str = `<xml>
  				<return_code><![CDATA[FAIL]]></return_code>
  				<return_msg><![CDATA[` + e.Error() + `]]></return_msg>
			   </xml>`

	} else {
		order, ok := g_order[paycb.Out_trade_no]
		if ok {
			str = `<xml>
  				<return_code><![CDATA[SUCCESS]]></return_code>
  				<return_msg><![CDATA[OK]]></return_msg>
			  </xml>`

			order.PayCb = paycb
		} else {
			str = `<xml>
  				<return_code><![CDATA[FAIL]]></return_code>
  				<return_msg><![CDATA[` + "no order!" + `]]></return_msg>
			   </xml>`
		}

	}

	session.WriteString(str)
}

func pubRefund(session *JsHttp.Session) {
	type Para struct {
		OrderId string
	}

	para := &Para{}
	e := session.GetPara(para)
	if e != nil {
		Error(e.Error())
		session.Forward("1", e.Error(), nil)
		return
	}

	order, ok := g_order[para.OrderId]
	if !ok {
		Error("OrderId = %v no match order", para.OrderId)
		session.Forward("2", "no order", nil)
		return
	}

	order.RefundId = strconv.Itoa(time.Now().Nanosecond())

	order.RefundFee = 10

	m, e := pubPayHandler.Wx_refund(order)
	if e != nil {
		Error(e.Error())
		session.Forward("1", e.Error(), nil)
		return
	}

	session.Forward("0", "success", m)
}

//////////////////////////////////////////////////////////////////////////////////////////
//app支付                                                                               //
/////////////////////////////////////////////////////////////////////////////////////////

func appNewOrder(session *JsHttp.Session) {
	type Para struct {
		Desc    string
		OrderId string
	}
	para := &Para{}

	e := session.GetPara(para)
	if e != nil {
		Error(e.Error())
		session.Forward("1", e.Error(), nil)
		return
	}

	addrs := strings.Split(session.Req.RemoteAddr, ":")

	order := &jswxpay.StOrder{}
	order.Desc = para.Desc
	order.OpenId = para.OrderId
	order.ProjectId = ""
	order.TerminalIp = addrs[0]
	order.Amount = 10
	order.OrderId = strconv.Itoa(time.Now().Nanosecond())
	order.Trade_type = "APP"

	m, e := appPayHandler.Wx_pay(order)
	if e != nil {
		Error(e.Error())
		session.Forward("2", e.Error(), nil)
		return
	}

	g_order[order.OrderId] = order
	m["OrderId"] = order.OrderId

	session.Forward("0", "success", m)

}

func appWxPayCb(session *JsHttp.Session) {

	str := ""
	body, e := ioutil.ReadAll(session.Req.Body)
	if e != nil {
		str = `<xml>
  				<return_code><![CDATA[FAIL]]></return_code>
  				<return_msg><![CDATA[` + e.Error() + `]]></return_msg>
			   </xml>`
		session.WriteString(str)
		return
	}

	paycb := &jswxpay.ST_PayCb{}
	e = xml.Unmarshal(body, paycb)

	if e != nil {
		str = `<xml>
  				<return_code><![CDATA[FAIL]]></return_code>
  				<return_msg><![CDATA[` + e.Error() + `]]></return_msg>
			   </xml>`

	} else {
		order, ok := g_order[paycb.Out_trade_no]
		if ok {
			str = `<xml>
  				<return_code><![CDATA[SUCCESS]]></return_code>
  				<return_msg><![CDATA[OK]]></return_msg>
			  </xml>`

			order.PayCb = paycb
		} else {
			str = `<xml>
  				<return_code><![CDATA[FAIL]]></return_code>
  				<return_msg><![CDATA[` + "no order!" + `]]></return_msg>
			   </xml>`
		}

	}

	session.WriteString(str)
}

func appRefund(session *JsHttp.Session) {
	type Para struct {
		OrderId string
	}

	para := &Para{}
	e := session.GetPara(para)
	if e != nil {
		Error(e.Error())
		session.Forward("1", e.Error(), nil)
		return
	}

	order, ok := g_order[para.OrderId]
	if !ok {
		Error("OrderId = %v no match order", para.OrderId)
		session.Forward("2", "no order", nil)
		return
	}

	order.RefundId = strconv.Itoa(time.Now().Nanosecond())

	order.RefundFee = 10

	m, e := appPayHandler.Wx_refund(order)
	if e != nil {
		Error(e.Error())
		session.Forward("1", e.Error(), nil)
		return
	}

	session.Forward("0", "success", m)
}
