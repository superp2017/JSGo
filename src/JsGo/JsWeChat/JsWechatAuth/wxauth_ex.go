package JsWechatAuth

import (
	"JsGo/JsHttp"
	"JsGo/JsLogger"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/chanxuehong/wechat/mp/user/oauth2"
)

// { "access_token":"ACCESS_TOKEN",
// 	"expires_in":7200,
// 	"refresh_token":"REFRESH_TOKEN",
// 	"openid":"OPENID",
// 	"scope":"SCOPE" }

type AToken struct {
	Access_token  string `json:access_token`
	Refresh_token string `json:refresh_token`
	Openid        string `json:openid`
	Scope         string `json:scope`
}

var (
	g_auth_ex      string               = "/wxauthex"
	g_authcb_ex    string               = "/wxauthcbex"
	g_oauth2Config *oauth2.OAuth2Config = nil
)

func ex_init() {
	fmt.Printf("weChatOAuth2 = %s\n", weChatOAuth2)
	g_oauth2Config = oauth2.NewOAuth2Config(weChatAppId, weChatSecret, weChatOAuth2+g_authcb_ex, "snsapi_userinfo")
	JsHttp.WhiteHttp(g_auth_ex, auth_ex)
	JsHttp.WhiteHttp(g_authcb_ex, authCb_ex)
}

func auth_ex(session *JsHttp.Session) {
	UID := session.Req.FormValue("UID")
	//channel := session.Req.FormValue("channel")
	type home_ret struct {
		Ret string
		Msg string
	}
	ret := &home_ret{}
	if len(UID) == 0 {
		oauth_ex(session, g_wxcb)
		return
	}
	if len(UID) > 0 {
		rs := strings.Split(session.Req.URL.String(), "?")
		if len(rs) == 2 {
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

func oauth_ex(session *JsHttp.Session, cb WxAuthCb) {
	g_mutex.Lock()
	defer g_mutex.Unlock()
	sid := fmt.Sprintf("%d", g_rand.Int63())
	auth := &StAuth{cb, session.Req.RequestURI, session.Req.FormValue("channel")}
	g_cbmap[sid] = auth

	AuthCodeURL := g_oauth2Config.AuthCodeURL(sid, nil)

	fmt.Printf("AuthCodeURL = %s\n", AuthCodeURL)
	http.Redirect(session.Rsp, session.Req, AuthCodeURL, 302)
}

func authCb_ex(session *JsHttp.Session) {
	code := session.Req.FormValue("code")
	sid := session.Req.FormValue("state")

	if code == "" {
		session.Forward("-1", "code is nil", nil)
		return
	}

	url := `https://api.weixin.qq.com/sns/oauth2/access_token?appid=` + weChatAppId + `&secret=` + weChatSecret + `&code=` + code + `&grant_type=authorization_code`

	resp, e := http.Get(url)
	if e != nil {
		fmt.Println(e.Error())
	}

	buffer := make([]byte, 2048)
	n, _ := resp.Body.Read(buffer)

	fmt.Println(string(buffer[:n]))

	token := &AToken{}
	e = json.Unmarshal(buffer[:n], token)
	if e != nil {
		JsLogger.Error(e.Error())
		return
	}

	url = `https://api.weixin.qq.com/sns/userinfo?access_token=` + token.Access_token + `&openid=` + token.Openid + `&lang=zh_CN`
	fmt.Println(url)
	resp, e = http.Get(url)
	defer resp.Body.Close()
	if e != nil {
		fmt.Println(e.Error())
	}

	buff := make([]byte, 2048)
	n, _ = resp.Body.Read(buff)

	fmt.Println(string(buff[:n]))

	userinfo := &oauth2.UserInfo{}
	e = json.Unmarshal(buff[:n], userinfo)
	if e != nil {
		JsLogger.Error(e.Error())
		return
	}
	fmt.Printf("unionid = %s, openid = %s\n", userinfo.UnionId, userinfo.OpenId)

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

		JsLogger.Error(weChatOriHome + auth.url + rs)
		http.Redirect(session.Rsp, session.Req, weChatOriHome+auth.url+rs, 302)
		delete(g_cbmap, sid)
	}
}
