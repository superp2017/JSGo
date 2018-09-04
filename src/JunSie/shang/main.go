package main

import (
	"JsGo/JsExit"
	"JsGo/JsHttp"
	"JsGo/JsMobile"
	"JsGo/JsQiniu"
	"JunSie/com"
)

func exit() int {
	JsHttp.Close()
	return 0
}

func main() {
	JsExit.RegisterExitCb(exit) //退出
	JsQiniu.QiniuInit()
	JsMobile.AlidayuInit()
	com.InitBusiness()
	com.InitWebConfig()
	JsHttp.Run()
}
