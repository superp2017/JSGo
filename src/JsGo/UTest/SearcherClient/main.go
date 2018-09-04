package main

import (
	"log"
	"net/rpc"
	"net/rpc/jsonrpc"
)

func index(rpcClient *rpc.Client) {
	doc1 := map[string]interface{}{}
	doc1["ID"] = "11111"
	doc1["Type"] = "1"
	doc1["Text"] = "Google Is Experimenting With Virtual Reality Advertising"

	doc2 := map[string]interface{}{}
	doc2["ID"] = "2222"
	doc2["Type"] = "2"
	doc2["Text"] = "Google accidentally pushed Bluetooth update for Home speaker early"

	doc3 := map[string]interface{}{}
	doc3["ID"] = "33333"
	doc3["Type"] = "1"
	doc3["Text"] = "Google Is Experimenting With Virtual Reality Advertising"

	var reply map[string]interface{}
	e := rpcClient.Call("JSHandler.Index", doc1, &reply)
	if e != nil {
		log.Println("Error call rpc method:", e)
		return
	}

	e = rpcClient.Call("JSHandler.Index", doc2, &reply)
	if e != nil {
		log.Println("Error call rpc method:", e)
		return
	}

	e = rpcClient.Call("JSHandler.Index", doc3, &reply)
	if e != nil {
		log.Println("Error call rpc method:", e)
		return
	}
}

func search(rpcClient *rpc.Client) {
	sText := make(map[string]string)
	sText["1"] = "Home"
	sText["2"] = "google"

	var reply map[string]interface{}
	e := rpcClient.Call("JSHandler.Query", sText, &reply)
	if e != nil {
		log.Println("Error call rpc method:", e)
		return
	}
	R1, ok := reply["1"]
	log.Println(R1)

	if ok {
		X1, ok := R1.(map[string]interface{})
		if ok {
			for k, v := range X1 {
				log.Println(k)

				b, ok := v.(string)
				if ok {
					log.Println(b)
				}

			}
		}
		// for k, _ := range R1 {
		// 	log.Println(k)
		// }

		//log.Println(R1)
	}
	// R2, ok := reply["2"]
	// if ok {
	// 	log.Println(R2)
	// }
}

func main() {
	// RPC calls.

	rpcClient, e := jsonrpc.Dial("tcp", "127.0.0.1:5135")
	if e != nil {
		log.Println("Error dail rpc server:", e)

		return
	}

	// index(rpcClient)
	//
	search(rpcClient)
}
