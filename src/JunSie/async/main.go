package main

import "JsGo/JsHttp"

func exit() int {

	JsHttp.Close()
	return 0
}

func main() {
	JsHttp.Run()
}
