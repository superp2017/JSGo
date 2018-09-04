package main

import (
	"JsGo/JsHttp"
)

func content(session *JsHttp.Session) {
	session.WriteString("Hello, I am https!")
}

func main() {
	JsHttp.Http("/getContent", content)
	JsHttp.Run()
}
