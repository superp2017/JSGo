package main

import (
	"JsGo/JsConfig"
	"fmt"
)

func main() {

	keys := make([]string, 2)

	keys[0] = "Hello"
	keys[1] = "Meng"

	ret, err := JsConfig.GetConfigString(keys)
	if err == nil {
		fmt.Println(ret)
	} else {
		fmt.Println(err)
	}

	keys = make([]string, 3)

	keys[0] = "xxx"
	keys[1] = "yyy"
	keys[2] = "zzz"

	rmap, err := JsConfig.GetConfigString(keys)
	if err == nil {
		fmt.Println(rmap)
	} else {
		fmt.Println(err)
	}

	keys[0] = "a"
	keys[1] = "yyy"
	keys[2] = "zzz"

	av, err := JsConfig.GetConfigInteger(keys)
	if err == nil {
		fmt.Println(av)
	} else {
		fmt.Println(err)
	}

	keys[0] = "b"
	keys[1] = "yyy"
	keys[2] = "zzz"

	bv, err := JsConfig.GetConfigFloat(keys)
	if err == nil {
		fmt.Println(bv)
	} else {
		fmt.Println(err)
	}

	keys = make([]string, 2)
	keys[0] = "c"
	keys[1] = "yyy"

	cv, err := JsConfig.GetConfigMap(keys)
	if err == nil {
		fmt.Printf("cv = %v\n", cv)
	} else {
		fmt.Println(err)
	}
}
