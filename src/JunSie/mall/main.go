package main

import (
	"JsGo/JsExit"
	"JsGo/JsHttp"
	"JsGo/JsMobile"
	"JsGo/JsQiniu"
	"JsGo/JsWeChat/JsAppAuth"
	"JsGo/JsWeChat/JsMiniP"
	"JsGo/JsWeChat/JsWechatAuth"
	"JsGo/JsWeChat/JsWechatJsApi/JsJdk"
	"JunSie/com"
	"JunSie/config"
)

func exit() int {
	JsHttp.Close()
	return 0
}

func main() {
	JsExit.RegisterExitCb(exit)
	config.InitConfig(func() { //配置(微信相关)必须先加载
		JsWechatAuth.WxauthInit(com.WxNewUser) //公众号认证
		JsAppAuth.AppInit(com.WxNewUser)       //App认证
		JsMiniP.MiniPInit(com.WxNewUser)       //公小程序认证
		JsJdk.JsJdkInit()                      //wxjsjdk
		com.InitPay()                          //支付
	})

	JsQiniu.QiniuInit()
	// JsMobile.AlidayuInit()
	JsMobile.NewAlidayuInit()
	com.Init_productMall()
	com.Init_contentMall()
	com.Init_commentMall()
	com.InitUserMall()
	com.InitHome()
	com.InitModuleMall()
	com.Peoplerecord()
	com.InitShow()
	com.InitShoppingCart()
	com.InitOrderMall()
	com.Init_MessageMall()
	com.InitSearchMall()
	com.InitBusinessMall()
	com.InitWebConfig()
	JsHttp.Run()
}
