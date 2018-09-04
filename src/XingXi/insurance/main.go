package main

import (
	"JsGo/JsExit"
	"JsGo/JsHttp"
	"JsGo/JsMobile"
	"JsGo/JsQiniu"
	"JsGo/JsWeChat/JsAppAuth"
	"JsGo/JsWeChat/JsWechatAuth"
	"JsGo/JsWeChat/JsWechatJsApi/JsJdk"
	"JunSie/user"
)

func exit() int {
	JsHttp.Close()
	return 0
}

func initInsurance() {
	JsHttp.WhiteHttp("/queryinsurance", QueryInsurance)     //查询单个保单
	JsHttp.WhiteHttp("/creatinsurance", CreatInsurance)     //创建保单
	JsHttp.WhiteHttp("/getuserinsurance", GetUserInsurance) //查询用户所有的保单
}
func order_init() {
	JsHttp.WhiteHttp("/submitorder", newOrder)              //提交订单
	JsHttp.WhiteHttp("/paysuccesscb", paySuccessCb)         //微信支付回调
	JsHttp.WhiteHttp("/cancelorder", orderCancle)           //取消订单
	JsHttp.WhiteHttp("/getuserorderlist", getUserOrderList) //获取用户所有订单
	JsHttp.WhiteHttp("/getorderinfo", getOrderInfo)         //获取单个表单
}

func main() {
	JsExit.RegisterExitCb(exit)
	JsWechatAuth.WxauthInit(user.WxNewUser)
	JsAppAuth.AppInit(user.WxNewUser)
	JsJdk.JsJdkInit_Unsafe()
	JsQiniu.QiniuInit()
	JsMobile.AlidayuInit()
	user.InitUser()
	initInsurance()
	order_init()
	JsHttp.Run()
}
