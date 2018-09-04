package JsWechatAuth

import (
	"JsGo/JsHttp"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"sync"
	"time"

	"JsGo/JsLogger"

	"JsGo/JsConfig"

	"github.com/chanxuehong/wechat/mp/user/oauth2"
)

type WxAuthCb func(user *oauth2.UserInfo, session *JsHttp.Session, userChannel string) string

type StAuth struct {
	cb      WxAuthCb
	url     string
	channel string
	// session *JsHttp.Session
}

var g_rand rand.Source
var g_cbmap map[string]*StAuth
var g_mutex sync.Mutex
var g_exec bool = false

var (
	weChatAppId        string
	weChatSecret       string
	weChatOAuth2       string
	weChatOriHome      string
	weChatRedirectHome string
)

var (
	g_auth   string = "/wxauth"
	g_authcb string = "/wxauthcb"
)

var g_wxcb WxAuthCb = nil
var g_channel = ""

func WxauthInit(cb WxAuthCb) {
	g_mutex.Lock()
	defer g_mutex.Unlock()
	if g_exec {
		return
	}
	g_wxcb = cb
	g_rand = rand.NewSource(int64(time.Now().Nanosecond()))
	g_cbmap = make(map[string]*StAuth)

	config := JsConfig.GetConfig()
	if config == nil {
		log.Fatalln("JsConfig.GetConfig is nil")
		return
	}
	weChatAppId = config.WxJsApi.WeChatAppId
	weChatSecret = config.WxJsApi.WeChatSecret
	weChatOAuth2 = config.WxJsApi.WeChatOAuth2
	weChatOriHome = config.WxJsApi.WeChatOriHome
	weChatRedirectHome = config.WxJsApi.WeChatRedirectHome
	g_authcb = config.WxJsApi.WeChatOAuth2Path

	//var e error
	//weChatAppId, e = JsConfig.GetConfigString([]string{"WxJsApi", "WeChatAppId"})
	//if e != nil {
	//	log.Fatalln(e.Error())
	//}
	//weChatSecret, e = JsConfig.GetConfigString([]string{"WxJsApi", "WeChatSecret"})
	//if e != nil {
	//	log.Fatalln(e.Error())
	//}
	//weChatOAuth2, e = JsConfig.GetConfigString([]string{"WxJsApi", "WeChatOAuth2"})
	//if e != nil {
	//	log.Fatalln(e.Error())
	//}
	//weChatOriHome, e = JsConfig.GetConfigString([]string{"WxJsApi", "WeChatOriHome"})
	//if e != nil {
	//	log.Fatalln(e.Error())
	//}
	//weChatRedirectHome, e = JsConfig.GetConfigString([]string{"WxJsApi", "WeChatRedirectHome"})
	//if e != nil {
	//	log.Fatalln(e.Error())
	//}
	//g_authcb, e = JsConfig.GetConfigString([]string{"WxJsApi", "WeChatOAuth2Path"})
	//if e != nil {
	//	log.Fatalln(e.Error())
	//}

	JsHttp.WhiteHttp(g_auth, auth)
	JsHttp.WhiteHttp(g_authcb, authCb)

	JsHttp.WhiteHttps(g_auth, auth)
	JsHttp.WhiteHttps(g_authcb, authCb)
	ex_init()
	g_exec = true
}

func auth(session *JsHttp.Session) {
	UID := session.Req.FormValue("UID")
	//channel := session.Req.FormValue("channel")
	type home_ret struct {
		Ret string
		Msg string
	}
	ret := &home_ret{}
	if len(UID) == 0 {
		oauth(session, g_wxcb)
		return
	}
	if len(UID) > 0 {
		rs := strings.Split(session.Req.URL.String(), "?")
		if len(rs) == 2 {
			//if channel != "" {
			//	redirect, _ := JsConfig.GetConfigString([]string{"WxJsApi", channel})
			//	if redirect != "" {
			//		http.Redirect(session.Rsp, session.Req, redirect+"?"+rs[1], 302)
			//	}
			//} else {
			http.Redirect(session.Rsp, session.Req, weChatRedirectHome+"?"+rs[1], 302)
			//	}
		} else {
			ret.Ret = "2"
			ret.Msg = "uid == nil"
			session.Forward("2", "uid == nil", nil)
		}
	} else {
		ret.Ret = "1"
		ret.Msg = "uid == nil"
		session.Forward("1", "uid == nil", nil)
	}
}

//////////
//      //
//微信  //
//////////

func oauth(session *JsHttp.Session, cb WxAuthCb) {
	g_mutex.Lock()
	defer g_mutex.Unlock()
	sid := fmt.Sprintf("%d", g_rand.Int63())
	auth := &StAuth{cb, session.Req.RequestURI, session.Req.FormValue("channel")}
	g_cbmap[sid] = auth
	oauth2Config := oauth2.NewOAuth2Config(weChatAppId, weChatSecret, weChatOAuth2+g_authcb, "snsapi_userinfo")
	AuthCodeURL := oauth2Config.AuthCodeURL(sid, nil)
	http.Redirect(session.Rsp, session.Req, AuthCodeURL, 302)
}

func authCb(session *JsHttp.Session) {
	code := session.Req.FormValue("code")
	sid := session.Req.FormValue("state")

	if code == "" {
		session.Forward("-1", "code is nil", nil)
		return
	}
	oauth2Client := oauth2.Client{
		Config: oauth2.NewOAuth2Config(weChatAppId, weChatSecret, weChatOAuth2+g_authcb, "snsapi_userinfo"),
	}
	_, err := oauth2Client.Exchange(code)
	if err != nil {
		session.Forward("-1", err.Error(), nil)
		return
	}
	userinfo, err := oauth2Client.UserInfo(oauth2.Language_zh_CN)
	if err != nil {
		session.Forward("-1", err.Error(), nil)
		return
	}
	g_mutex.Lock()
	defer g_mutex.Unlock()
	auth, ok := g_cbmap[sid]
	UID := ""
	if ok {
		if auth.cb != nil {
			UID = auth.cb(userinfo, session, "web")
		}
		JsLogger.Error("UID:", UID)
		rs := ""
		if strings.Index(auth.url, "?") == -1 {
			rs = "?UID=" + UID
		} else {
			rs = "&UID=" + UID
		}
		//rs += "&channel="
		//rs += auth.channel

		JsLogger.Error(weChatOriHome + auth.url + rs)
		http.Redirect(session.Rsp, session.Req, weChatOriHome+auth.url+rs, 302)
		delete(g_cbmap, sid)
	}
}
