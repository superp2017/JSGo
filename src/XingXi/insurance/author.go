package main

import (
	"JsGo/JsLogger"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"
)

type ResData struct {
	Open_id         string `xml:"open_id"`
	Expires_in      int    `xml:"expires_in"`
	Refresh_token   string `xml:"refresh_token"`
	Access_tokencvb string `xml:"access_token"`
}
type ResMsg struct {
	Message string  `xml:"message"`
	Code    string  `xml:"code"`
	Data    ResData `xml:"data"`
	State   string  `xml:"state"`
}
type AuthorRespond struct {
	ResponseMessage ResMsg `xml:"response-message"`
}

var g_token ResData
var g_mutex sync.Mutex

//鉴权认证
func author() error {
	authorUrl := "http://183.60.22.143/sns/oauth2/authorize?app_id=87FE4376AA9ACF82909611380B8028B5&app_secret=FD3F2E8231DB556842926D781D9FC7E0&code=1anOpcode&grant_type=authorization_code"
	response, e := http.Get(authorUrl)
	defer response.Body.Close()
	if e != nil {
		return e
	}
	b, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return e
	}
	v := AuthorRespond{}
	err = json.Unmarshal(b, &v)
	if err != nil {
		fmt.Printf("error: %v", err)
		return err
	}
	g_mutex.Lock()
	g_token = v.ResponseMessage.Data
	g_mutex.Unlock()
	return nil
}

func refreshToken() error {
	if err := author(); err != nil {
		return err
	}
	appid, ref := g_token.Open_id, g_token.Refresh_token
	url := "http://183.60.22.143/sns/oauth2/token_refresh?app_id=" + appid + "&app_secret=FD3F2E8231DB556842926D781D9FC7E0&grant_type=refresh_token&refresh_token=" + ref
	response, e := http.Get(url)
	defer response.Body.Close()
	if e != nil {
		return e
	}
	b, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return e
	}
	data := &ResMsg{}
	if err := json.Unmarshal(b, data); err != nil {
		return e
	}
	if data.State == "1" {
		g_mutex.Lock()
		g_token = data.Data
		g_mutex.Unlock()
		return nil
	}
	return nil
}

func creatOrder() error {
	if err := author(); err != nil {
		return err
	}
	token, appid := g_token.Access_token, g_token.Open_id
	url := "http://183.60.22.143/sns/dy-order-service-dev? access_token=" + token + "&open_id=" + appid
	response, e := http.Get(url)
	defer response.Body.Close()
	if e != nil {
		return e
	}
	b, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return e
	}
	data := &ResMsg{}
	if err := json.Unmarshal(b, data); err != nil {
		return e
	}
	if data.State == "1" {
		return nil
	}
	return JsLogger.ErrorLog("创建失败,err=" + data.Message)
}

func cancleOrder() error {
	token, appid := g_token.Access_token, g_token.Open_id
	url := "http://183.60.22.143/sns/cancel-service-sit? access_token=" + token + "&open_id=" + appid
	response, e := http.Get(url)
	defer response.Body.Close()
	if e != nil {
		return e
	}
	b, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return e
	}
	data := &ResMsg{}
	if err := json.Unmarshal(b, data); err != nil {
		return e
	}
	if data.State == "1" {
		return nil
	}
	return JsLogger.ErrorLog("退保失败,err=" + data.Message)
}
