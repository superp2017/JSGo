package main

import (
	"JsGo/JsExit"
	"JsGo/JsHttp"
	"JsGo/JsQiniu"
)

func exit() int {
	JsHttp.Close()
	return 0
}

func main() {
	JsExit.RegisterExitCb(exit)
	JsQiniu.QiniuInit()
	JsHttp.Run()
}
