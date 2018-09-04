package jswxpay

import (
	"JsGo/JsConfig"
	. "JsGo/JsLogger"
	"JsGo/JsUuid"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"
)

// const (
// 	// appId  = CFG.DirectPay.AppId     // 微信公众平台应用ID
// 	// mchId  = CFG.DirectPay.MchId     // 微信支付商户平台商户号
// 	// apiKey = CFG.DirectPay.SecretKey // 微信支付商户平台API密钥

// 	// 微信支付商户平台证书路径
// 	certFile   = "cert/apiclient_cert.pem"
// 	keyFile    = "cert/apiclient_key.pem"
// 	rootcaFile = "cert/rootca.pem"
// )

// var (
// 	appId          string // 微信公众平台应用ID
// 	mchId          string // 微信支付商户平台商户号
// 	apiKey         string // 微信支付商户平台API密钥
// 	wxPubPayCb     string //回调URL
// 	wxPubPayUrl    string //付款URL
// 	wxPubRefundUrl string //退款URL
// )

const (
	base64Table = "123QRSTUabcdVWXYZHijKLAWDCABDstEFGuvwxyzGHIJklmnopqr234560178912"
)

type PayHandler struct {
	WxAppId          string // 微信公众平台应用ID
	WxSecret         string
	WxMchId          string // 微信支付商户平台商户号
	WxSecretKey      string // 微信支付商户平台API密钥
	WxPubPayCb       string //回调URL
	WxPubPayUrl      string //付款URL
	WxPubRefundUrl   string //退款URL
	WxPubTransferUrl string
	WxPubSendredpack string
	WxSpbillCreateIp string

	CertFile   string
	KeyFile    string
	RootcaFile string

	C      *Client
	Coder  *base64.Encoding
	EnvKey string
}

func NewPayHandler(env string) *PayHandler {
	handler := &PayHandler{}
	handler.Init_wx_pay(env)
	return handler
}

func (pay *PayHandler) Init_wx_pay(env string) {
	pay.EnvKey = env
	pay.Coder = base64.NewEncoding(base64Table)

	config := JsConfig.GetConfig()
	if config == nil {
		return
	}
	var PayType JsConfig.St_Pay
	if pay.EnvKey == "PubPay" {
		PayType = config.PubPay
	}
	if pay.EnvKey == "WxMiniP" {
		PayType = config.WxMiniP
	}
	if pay.EnvKey == "AppPay" {
		PayType = config.AppPay
	}
	pay.WxAppId = PayType.WxAppId
	pay.WxSecret = PayType.WxSecret
	pay.WxMchId = PayType.WxMchId
	pay.WxSecretKey = PayType.WxSecretKey
	pay.WxPubPayCb = PayType.WxPubPayCb
	pay.WxPubPayUrl = PayType.WxAppId
	pay.WxPubRefundUrl = PayType.WxPubPayUrl
	pay.WxPubTransferUrl = PayType.WxPubRefundUrl
	pay.WxPubSendredpack = PayType.WxPubSendredpack
	pay.WxSpbillCreateIp = PayType.WxSpbillCreateIp
	pay.CertFile = PayType.CertFile
	pay.KeyFile = PayType.KeyFile

	//var e error
	//pay.WxAppId, e = JsConfig.GetConfigString([]string{pay.EnvKey, "WxAppId"})
	//if e != nil {
	//	log.Fatalln(e.Error())
	//}
	//
	//pay.WxSecret, e = JsConfig.GetConfigString([]string{pay.EnvKey, "WxSecret"})
	//if e != nil {
	//	log.Fatalln(e.Error())
	//}
	//
	//pay.WxMchId, e = JsConfig.GetConfigString([]string{pay.EnvKey, "WxMchId"})
	//if e != nil {
	//	log.Fatalln(e.Error())
	//}
	//pay.WxSecretKey, e = JsConfig.GetConfigString([]string{pay.EnvKey, "WxSecretKey"})
	//if e != nil {
	//	log.Fatalln(e.Error())
	//}
	//pay.WxPubPayCb, e = JsConfig.GetConfigString([]string{pay.EnvKey, "WxPubPayCb"})
	//if e != nil {
	//	log.Fatalln(e.Error())
	//}
	//pay.WxPubPayUrl, e = JsConfig.GetConfigString([]string{pay.EnvKey, "WxPubPayUrl"})
	//if e != nil {
	//	log.Fatalln(e.Error())
	//}
	//pay.WxPubRefundUrl, e = JsConfig.GetConfigString([]string{pay.EnvKey, "WxPubRefundUrl"})
	//if e != nil {
	//	log.Fatalln(e.Error())
	//}
	//pay.WxPubTransferUrl, e = JsConfig.GetConfigString([]string{pay.EnvKey, "WxPubTransferUrl"})
	//if e != nil {
	//	log.Fatalln(e.Error())
	//}
	//pay.WxPubSendredpack, e = JsConfig.GetConfigString([]string{pay.EnvKey, "WxPubSendredpack"})
	//if e != nil {
	//	log.Fatalln(e.Error())
	//}
	//pay.WxSpbillCreateIp, e = JsConfig.GetConfigString([]string{pay.EnvKey, "WxSpbillCreateIp"})
	//if e != nil {
	//	log.Fatalln(e.Error())
	//}
	//
	//pay.CertFile, e = JsConfig.GetConfigString([]string{pay.EnvKey, "CertFile"})
	//if e != nil {
	//	log.Fatalln(e.Error())
	//}
	//
	//pay.KeyFile, e = JsConfig.GetConfigString([]string{pay.EnvKey, "KeyFile"})
	//if e != nil {
	//	log.Fatalln(e.Error())
	//}

	// pay.RootcaFile, e = JsConfig.GetConfigString([]string{pay.EnvKey, "RootcaFile"})
	// if e != nil {
	// 	log.Fatalln(e.Error())
	// }

	pay.C = NewClient(pay.WxAppId, pay.WxMchId, pay.WxSecretKey)

	// 附着商户证书
	err := pay.C.WithCert(pay.CertFile, pay.KeyFile)
	if err != nil {
		log.Fatal(err)
	}
}

func (pay *PayHandler) Wx_pay(order *StOrder) (map[string]string, error) {

	b := make([]byte, 8)
	rand.Read(b)
	nonce_str := strings.ToUpper(pay.Coder.EncodeToString(b))

	order.AppId = pay.C.AppId
	order.Mch_id = pay.C.MchId
	order.Nonce_str = nonce_str

	params := make(Params)
	// 查询企业付款接口请求参数
	params.SetString("appid", pay.C.AppId)
	params.SetString("mch_id", pay.C.MchId)
	params.SetString("device_info", "WEB")
	params.SetString("nonce_str", nonce_str) // 随机字符串
	params.SetString("body", order.Desc)
	params.SetString("out_trade_no", order.OrderId) // 商户订单号
	params.SetString("total_fee", strconv.Itoa(order.Amount))
	params.SetString("spbill_create_ip", order.TerminalIp)
	params.SetString("notify_url", pay.WxPubPayCb)
	params.SetString("trade_type", order.Trade_type)
	params.SetString("attach", order.ProjectId)
	params.SetString("openid", order.OpenId)

	params.SetString("sign", pay.C.Sign(params)) // 签名

	//url := "https://api.mch.weixin.qq.com/pay/unifiedorder"

	ret, err := pay.C.Post(pay.WxPubPayUrl, params, true)
	if err != nil {
		Error(err.Error())
		return nil, err
	}

	//fmt.Printf("ret = %v\n", ret)

	if ret["return_code"] == "FAIL" {
		return nil, ErrorLog("微信提交订单失败,return_msg:%s,err_code_des:%s\n", ret["return_msg"], ret["err_code_des"])
	}

	charge := make(map[string]string)

	if order.Trade_type == "JSAPI" {
		charge["appId"] = pay.C.AppId
		charge["timeStamp"] = fmt.Sprintf("%d", time.Now().Unix())
		charge["nonceStr"] = ret["nonce_str"]
		charge["package"] = "prepay_id=" + ret["prepay_id"]
		charge["signType"] = "MD5"
		charge["paySign"] = pay.C.Sign(charge)
	} else if order.Trade_type == "APP" {
		charge["appid"] = pay.C.AppId
		charge["timestamp"] = fmt.Sprintf("%d", time.Now().Unix())
		charge["noncestr"] = ret["nonce_str"]
		charge["partnerid"] = pay.C.MchId
		charge["package"] = "Sign=WXPay"
		charge["prepayid"] = ret["prepay_id"]
		charge["sign"] = pay.C.Sign(charge)
	}

	return charge, nil
}

func (pay *PayHandler) Wx_refund(order *StOrder) (map[string]string, error) {

	b := make([]byte, 8)
	rand.Read(b)
	nonce_str := strings.ToUpper(pay.Coder.EncodeToString(b))

	order.RefundId = JsUuid.NewV4().String()

	params := make(Params)
	// 查询企业付款接口请求参数
	params.SetString("appid", pay.C.AppId)
	params.SetString("mch_id", pay.C.MchId)
	params.SetString("device_info", "WEB")
	params.SetString("nonce_str", nonce_str) // 随机字符串
	params.SetString("transaction_id", order.PayCb.Transaction_id)
	params.SetString("out_trade_no", order.OrderId) // 商户订单号
	params.SetString("out_refund_no", order.RefundId)
	params.SetString("total_fee", strconv.Itoa(order.Amount))
	params.SetString("refund_fee", strconv.Itoa(order.RefundFee))
	params.SetString("op_user_id", pay.C.MchId)

	params.SetString("sign", pay.C.Sign(params)) // 签名

	//url := "https://api.mch.weixin.qq.com/secapi/pay/refund"

	ret, err := pay.C.Post(pay.WxPubRefundUrl, params, true)
	if err != nil {
		Error(err.Error())
		return nil, err
	}

	return ret, nil
}
