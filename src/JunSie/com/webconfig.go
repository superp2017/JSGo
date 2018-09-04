package com

import (
	"JsGo/JsHttp"
	"JsGo/JsLogger"
	"JsGo/JsStore/JsRedis"
)

type WebConfig struct {
	D_host_tls               string   `json:"d_host_tls"`
	S_host                   string   `json:"s_host"`
	Res_host                 string   `json:res_host`
	Menubk_string_js         string   `json:menubk_string_js`
	Minebk_string_js         string   `json:minebk_string_js`
	Startimg_string_js       []string `json:startimg_string_js`
	Review_status_js_boolean bool     `json:review_status_js_boolean`
	Share_name_js_string     string   `json:share_name_js_string`
	Share_desc_js_string     string   `json:share_desc_js_string`
	Share_url_js_string      string   `json:share_url_js_string`
	Default_header_string    string   `json:default_header_string`
	Share_logo_js_string     string   `json:share_logo_js_string`
	Test_uid                 string   `json:test_uid`
	Wxjsapi                  string   `json:Wxjsapi`
}

func InitWebConfig() {
	JsHttp.WhiteHttps("/getwebconfig", GetWebConfig)   //获取webconfig
	JsHttp.WhiteHttps("/savewebconfig", SaveWebConfig) //保存webconfig
}

func GetWebConfig(s *JsHttp.Session) {
	type Param struct {
		Name string
	}
	st := &Param{}
	if err := s.GetPara(st); err != nil {
		s.Forward("1", err.Error(), nil)
		return
	}
	if st.Name == "" {
		JsLogger.Error("GetWebConfig:Param Name is empty！")
		s.Forward("1", "GetWebConfig:Param Name is empty！", nil)
		return
	}
	data := &WebConfig{}

	hand := JsRedis.GetRedis("Redis_JS")
	if hand == nil {
		JsLogger.Error("请求中央后台失败！")
		s.Forward("1", "请求中央后台失败！", nil)
		return
	}
	if err := hand.Redis_hget("H_CONFIG", st.Name, data); err != nil {
		JsLogger.Error("1", err.Error(), nil)
		s.Forward("1", err.Error(), nil)
		return
	}
	s.Forward("0", "success", data)
}

func SaveWebConfig(s *JsHttp.Session) {
	type Para struct {
		Name   string
		Config WebConfig
	}

	st := &Para{}
	if err := s.GetPara(st); err != nil {
		s.Forward("1", err.Error(), nil)
		return
	}
	if st.Name == "" {
		JsLogger.Error("SaveWebConfig:Param Name is empty！")
		s.Forward("1", "SaveWebConfig:Param Name is empty！", nil)
		return
	}

	hand := JsRedis.GetRedis("Redis_JS")
	if hand == nil {
		JsLogger.Error("请求中央后台失败！")
		s.Forward("1", "请求中央后台失败！", nil)
		return
	}
	if err := hand.Redis_hset("H_CONFIG", st.Name, st.Config); err != nil {
		JsLogger.Error("1", err.Error(), nil)
		s.Forward("1", err.Error(), nil)
		return
	}
	s.Forward("0", "success！", nil)
}
