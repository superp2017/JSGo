package main

import (
	"JsGo/JsExit"
	"JsGo/JsHttp"
	"JsGo/JsQiniu"
	"JunSie/com"
	"JunSie/config"
)

func exit() int {
	JsHttp.Close()
	return 0
}

func main() {
	JsExit.RegisterExitCb(exit)
	config.InitConfig(nil)
	InitLogin()
	JsQiniu.QiniuInit()
	com.InitTag()
	com.Init_product()
	com.Init_content()
	com.Init_comment()
	com.InitUser()
	com.Find()
	com.InitHome()
	com.InitModule()
	com.Peoplerecord()
	com.InitShow()
	com.InitOrder()
	com.Init_Message()
	com.InitSearch()
	com.InitBusiness()
	com.Init_operation()
	com.InitWebConfig()
	JsHttp.Run()
}
