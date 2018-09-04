package main

import (
	"JsGo/JsConfig"
	"log"
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"

	"JsGo/JsLogger"
)

// type SearchDoc struct {
// 	ID   int
// 	Type string
// 	DID  string
// 	Text string
// }

type JSHandler struct {
}

func parsePara(args *map[string]interface{}) map[string]string {
	para := map[string]string{}

	for k, v := range *args {
		s, ok := v.(string)
		if ok {
			para[k] = s
		}
	}
	return para
}

func (h *JSHandler) Index(args *map[string]interface{}, reply *map[string]interface{}) error {
	defer JsLogger.TraceException()

	*reply = index(*args)
	return nil
}

func (h *JSHandler) Query(args *map[string]interface{}, reply *map[string]interface{}) error {
	defer JsLogger.TraceException()
	*reply = query(parsePara(args))

	return nil
}

var running bool = true
var ln *net.TCPListener

func exit() int {
	running = false
	ln.Close()
	searcher.Close()
	return 0
}

func main() {
	InitSearcher()

	jsHandler := new(JSHandler)
	rpc.Register(jsHandler)

	port, e := JsConfig.GetConfigString([]string{"Searcher", "Port"})
	if e != nil {
		log.Fatal(e.Error())
	}
	addr, _ := net.ResolveTCPAddr("tcp", ":"+port)
	ln, e = net.ListenTCP("tcp", addr)
	if e != nil {
		panic(e)
	}
	for {

		conn, e := ln.Accept()
		if e != nil {
			continue
		}
		if !running {
			break
		}

		go jsonrpc.ServeConn(conn)
	}
}
