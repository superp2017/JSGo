package main

import (
	"JsGo/JsHttp"
	"fmt"
)

func Hello(session *JsHttp.Session) {
	session.Forward("0", "success", "")
}

func main() {
	// JsHttp.Http("/hello", Hello)
	// JsHttp.Https("/hello", Hello)

	fmt.Printf("%08d", 1)
	// JsHttp.EnableSession()
	// JsHttp.EnableGet()
	// JsHttp.Run()
}
