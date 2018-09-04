package JsMiniP

import (
	"JsGo/JsBench/JsUser"
	"JsGo/JsHttp"
	"JsGo/JsLogger"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"JsGo/JsConfig"

	"github.com/chanxuehong/wechat/mp/user/oauth2"
)

var (
	weChatAppId  string
	weChatSecret string
)

type WxMiniPAuthCb func(user *oauth2.UserInfo, session *JsHttp.Session, userChannel string) string

var g_minipAuthCb WxMiniPAuthCb
var g_minipClient oauth2.Client

type SessionKeys struct {
	OpenID     string `json:"openid"`      // 用户的唯一标识
	SessionKey string `json:"session_key"` // 会话密钥
	UnionID    string `json:"unionid"`     // unionID
}

func MiniPInit(cb WxMiniPAuthCb) {
	if cb == nil {
		log.Fatalln("cb == nil")
	}
	g_minipAuthCb = cb

	config := JsConfig.GetConfig()
	if config == nil {
		log.Fatalln("JsConfig.GetConfig is nil")
		return
	}
	weChatAppId = config.WxMiniP.WxAppId
	weChatSecret = config.WxMiniP.WxSecret

	//var e error
	//weChatAppId, e = JsConfig.GetConfigString([]string{"WxMiniP", "WxAppId"})
	//if e != nil {
	//	log.Fatalln(e.Error())
	//}
	//weChatSecret, e = JsConfig.GetConfigString([]string{"WxMiniP", "WxSecret"})
	//if e != nil {
	//	log.Fatalln(e.Error())
	//}

	g_minipClient = oauth2.Client{
		Config: oauth2.NewOAuth2Config(weChatAppId, weChatSecret, "", "snsapi_userinfo"),
	}

	JsHttp.WhiteHttps("/code2session", CodeGetSessionKey)
	JsHttp.WhiteHttps("/minpUserInfo", minPUserInfo)

}

func minPUserInfo(session *JsHttp.Session) {

	type Para struct {
		Code  string
		State string
	}

	para := &Para{}

	e := session.GetPara(para)
	if e != nil {
		JsLogger.Error(e.Error())
		session.Forward("1", e.Error(), nil)
		return
	}

	_, e = g_minipClient.Exchange(para.Code)

	if e != nil {
		JsLogger.Error(e.Error())
		session.Forward("2", e.Error(), nil)
		return
	}

	userinfo, e := g_minipClient.UserInfo(oauth2.Language_zh_CN)
	if e != nil {
		JsLogger.Error(e.Error())
		session.Forward("2", e.Error(), nil)
		return
	}
	g_minipAuthCb(userinfo, session, "minip")
	session.Forward("0", "success", userinfo)
}

///小程序和App  Code获取openid和unionid
func CodeGetSessionKey(session *JsHttp.Session) {
	type INFo struct {
		Code string
	}
	st := &INFo{}
	if err := session.GetPara(st); err != nil {
		JsLogger.Error(err.Error())
		session.Forward("1", err.Error(), nil)
		return
	}
	if st.Code == "" {
		JsLogger.Error("CodeGetSessionKey Code isEmpty\n")
		session.Forward("1", "CodeGetSessionKey Code isEmpty\n", nil)
		return
	}
	url := "https://api.weixin.qq.com/sns/jscode2session?appid="
	url += weChatAppId
	url += "&secret="
	url += weChatSecret
	url += "&js_code="
	url += st.Code
	url += "&grant_type=authorization_code"

	JsLogger.Error(url)
	response, e := http.Get(url)

	if e != nil {
		JsLogger.Error(e.Error())
		session.Forward("1", e.Error(), nil)
		return
	}
	defer response.Body.Close()
	b, err := ioutil.ReadAll(response.Body)
	JsLogger.Error(string(b))
	if err != nil {
		JsLogger.Error(err.Error())
		session.Forward("1", e.Error(), nil)
		return
	}
	data := &SessionKeys{}
	if err := json.Unmarshal(b, data); err != nil {
		JsLogger.Error(err.Error())
		session.Forward("1", err.Error(), nil)
		return
	}

	JsLogger.Error("%v", data)
	user := &oauth2.UserInfo{
		OpenId:  data.OpenID,
		UnionId: data.UnionID,
	}
	uid := g_minipAuthCb(user, session, "minip")
	JsLogger.Error("uid=%s", uid)
	U := &JsUser.User{}
	U.ID = uid
	U.Openid = user.OpenId
	U.Unionid = user.UnionId
	session.Forward("0", "success", U)
}
