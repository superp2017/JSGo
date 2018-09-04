package com

import (
	"JsGo/JsConfig"
	"JsGo/JsHttp"
	"JsGo/JsLogger"
	"encoding/json"
	"log"
	"net/rpc"
	"net/rpc/jsonrpc"
)

var RPC_Client *rpc.Client = nil

type Index struct {
	ID   string
	Type string
}
type SearchData struct {
	Pro     []*XM_Product  //多个产品
	Content []*XM_Contents //多个内容
}

func (this *SearchData) appendPro(p *XM_Product) {
	exist := false
	for _, v := range this.Pro {
		if v.ID == p.ID {
			exist = true
			break
		}
	}
	if !exist {
		this.Pro = append(this.Pro, p)
	}
}

func (this *SearchData) appendContent(p *XM_Contents) {
	exist := false
	for _, v := range this.Content {
		if v.ID == p.ID {
			exist = true
			break
		}
	}
	if !exist {
		this.Content = append(this.Content, p)
	}
}

func InitSearch() {
	JsHttp.WhiteHttp("/serchinfo", SerchInfo) //搜索内容
	/////JsHttp.WhiteHttp("/indexsearch", indexSearch) //初始化内容
}
func InitSearchMall() {
	JsHttp.WhiteHttps("/serchinfo", SerchInfo) //搜索内容
}

func init() {
	if err := initRPC(); err != nil {
		log.Println("Error dail search rpc")
	}
}

func initRPC() error {

	host, e := JsConfig.GetConfigString([]string{"Searcher", "host"})
	if e != nil {
		return e
	}
	RPC_Client, e = jsonrpc.Dial("tcp", host)
	if e != nil {
		return e
	}
	return nil
}

func SerchInfo(session *JsHttp.Session) {
	type Para struct {
		Content string //搜索的内容
	}
	st := &Para{}
	if err := session.GetPara(st); err != nil {
		JsLogger.Error(err.Error())
		session.Forward("1", err.Error(), nil)
		return
	}
	if st.Content == "" {
		JsLogger.Error("搜索的内容不能为空!\n")
		session.Forward("1", "搜索的内容不能为空!\n", nil)
		return
	}

	data, err := search(st.Content)
	if err != nil {
		JsLogger.Error(err.Error())
		session.Forward("1", err.Error(), nil)
		return
	}
	serarch := &SearchData{}
	for _, v := range data {
		if v.Type == "product" {
			pro, err := GetProductInfo(v.ID)
			if err == nil {
				serarch.appendPro(pro)
			}
		}
		if v.Type == "content" {
			con, err := Getcontent(v.ID)
			if err == nil {
				serarch.appendContent(con)
			}
		}
	}
	if len(serarch.Pro) == 0 && len(serarch.Content) == 0 {
		session.Forward("0", "search success\n", nil)
		return
	}
	session.Forward("0", "search success\n", serarch)
}

func appendSearchContent(content map[string]string) error { //0.01
	if RPC_Client == nil {
		if err := initRPC(); err != nil {
			return err
		}
	}

	var reply map[string]interface{}
	e := RPC_Client.Call("JSHandler.Index", content, &reply)
	if e != nil {
		log.Println("Error call rpc method:", e)
		return e
	}
	return nil
}

func search(content string) ([]Index, error) {
	if RPC_Client == nil {
		if err := initRPC(); err != nil {
			return nil, err
		}
	}
	sText := make(map[string]string)
	sText["1"] = content

	var reply map[string]interface{}
	e := RPC_Client.Call("JSHandler.Query", sText, &reply)
	if e != nil {
		log.Println("Error call rpc method:", e)
		return nil, e
	}
	ret := []Index{}
	R1, ok := reply["1"]
	if ok {
		X1, ok := R1.(map[string]interface{})
		if ok {
			for _, v := range X1 {
				b, ok := v.(string)
				if ok {
					in := Index{}
					json.Unmarshal([]byte(b), &in)
					ret = append(ret, in)
				}
			}
		}
	}
	return ret, nil
}

func creatProSearchIndex(product *XM_Product) { //建立产品搜索索引0.00
	key := make(map[string]string)
	key["ID"] = product.ID
	key["Type"] = "product"
	text := ""
	if product.Title != "" {
		text += product.Title
	}
	if product.SubTitle != "" {
		text += product.SubTitle
	}
	if product.Desc != "" {
		text += product.Desc
	}
	if text != "" {
		key["Text"] = text
		go appendSearchContent(key)
	}
}

func creatContentSearchIndex(contents *XM_Contents) {
	key := make(map[string]string)
	key["ID"] = contents.ID
	key["Type"] = "content"
	text := ""
	if contents.Title != "" {
		text += contents.Title
	}
	if contents.SubTitle != "" {
		text += contents.Title
	}
	if contents.Brief != "" {
		text += contents.Brief
	}
	if contents.Author != "" {
		text += contents.Author
	}
	for _, v1 := range contents.Content {
		text += v1.Text
	}
	if text != "" {
		key["Text"] = text
		go appendSearchContent(key)
	}
}

func indexSearch(session *JsHttp.Session) {
	pro, err := getGlobalProducts()
	if err == nil {
		for _, v := range pro {
			p, e := GetProductInfo(v)
			if e == nil {
				go creatProSearchIndex(p)
			}
		}
	}

	content, err := getGlobalContent()
	if err == nil {
		for _, v := range content {
			con, e := Getcontent(v)
			if e == nil {
				go creatContentSearchIndex(con)
			}
		}
	}
	session.Forward("0", "indexSearch success\n", nil)
}
