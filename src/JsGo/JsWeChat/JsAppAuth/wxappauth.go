package JsAppAuth

import (
	"JsGo/JsHttp"
	. "JsGo/JsLogger"
	"log"

	"encoding/json"
	"io/ioutil"
	"net/http"

	"JsGo/JsConfig"

	"github.com/chanxuehong/wechat/mp/user/oauth2"
)

type AccessToken struct {
	Access_token  string `json:access_token`
	Expires_in    int    `json:expires_in`
	Refresh_token string `json:refresh_token`
	Openid        string `json:openid`
	Scope         string `json:scope`
	Unionid       string `json:unionid`
}

type WxAppAuthCb func(user *oauth2.UserInfo, session *JsHttp.Session, userChannel string) string

var (
	wxAppId        string
	wxSecret       string
	wxOAuth2       string
	wxOriHome      string
	wxRedirectHome string
)

var g_oauth2Client oauth2.Client
var g_appAuthCb WxAppAuthCb

func AppInit(cb WxAppAuthCb) {
	if cb == nil {
		log.Fatalln("cb == nil")
	}
	g_appAuthCb = cb

	config := JsConfig.GetConfig()
	if config == nil {
		log.Fatalln("JsConfig.GetConfig is nil")
		return
	}
	wxAppId = config.AppPay.WxAppId
	wxSecret = config.AppPay.WxSecret

	//var e error
	//wxAppId, e = JsConfig.GetConfigString([]string{"AppPay", "WxAppId"})
	//if e != nil {
	//	log.Fatalln(e.Error())
	//}
	//wxSecret, e = JsConfig.GetConfigString([]string{"AppPay", "WxSecret"})
	//if e != nil {
	//	log.Fatalln(e.Error())
	//}

	g_oauth2Client = oauth2.Client{
		Config: oauth2.NewOAuth2Config(wxAppId, wxSecret, "", "snsapi_userinfo"),
	}

	JsHttp.WhiteHttps("/appuserinfo", appUserInfo)
	JsHttp.WhiteHttps("/appcode2accesstoken", AppCode2AccessToken)

}

func appUserInfo(session *JsHttp.Session) {

	type Para struct {
		Code  string
		State string
	}

	para := &Para{}

	e := session.GetPara(para)
	if e != nil {
		Error(e.Error())
		session.Forward("1", e.Error(), nil)
		return
	}

	_, e = g_oauth2Client.Exchange(para.Code)

	if e != nil {
		Error(e.Error())
		session.Forward("2", e.Error(), nil)
		return
	}

	userinfo, e := g_oauth2Client.UserInfo(oauth2.Language_zh_CN)
	if e != nil {
		Error(e.Error())
		session.Forward("2", e.Error(), nil)
		return
	}
	UID := g_appAuthCb(userinfo, session, "app")
	session.Forward("0", "success", UID)
}

func AppCode2AccessToken(session *JsHttp.Session) {
	type Para struct {
		Code string
	}
	st := &Para{}
	if err := session.GetPara(st); err != nil {
		session.Forward("1", err.Error(), nil)
		return
	}
	if st.Code == "" {
		Error("AppCode2AccessToken field,code is empty\n")
		session.Forward("1", "AppCode2AccessToken field,code is empty\n", nil)
		return
	}
	url := "https://api.weixin.qq.com/sns/oauth2/access_token?appid="
	url += wxAppId
	url += "&secret="
	url += wxSecret
	url += "&code="
	url += st.Code
	url += "&grant_type=authorization_code"

	response, e := http.Get(url)
	defer response.Body.Close()
	if e != nil {
		Error(e.Error())
		session.Forward("1", e.Error(), nil)
		return
	}
	b, err := ioutil.ReadAll(response.Body)
	Error(string(b))
	if err != nil {
		Error(err.Error())
		session.Forward("1", e.Error(), nil)
		return
	}
	data := &AccessToken{}
	if err := json.Unmarshal(b, data); err != nil {
		Error(err.Error())
		session.Forward("1", err.Error(), nil)
		return
	}
	session.Forward("0", "success", data)
}
