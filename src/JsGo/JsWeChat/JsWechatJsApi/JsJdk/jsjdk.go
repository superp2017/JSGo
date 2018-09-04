package JsJdk

import (
	"JsGo/JsConfig"
	"JsGo/JsHttp"
	"JsGo/JsLogger"
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/astaxie/beego"
)

type ST_WeChat_AccessToken struct {
	Access_token string `json:"access_token"`
	Expires_in   int    `json:"expires_in"`
}

type ST_WeChat_Jsapi_Ticket struct {
	Errcode    int    `json:"errcode"`
	Errmsg     string `json:"errmsg"`
	Ticket     string `json:"ticket"`
	Expires_in int    `json:"expires_in"`
}

type ST_WeChatJsapiController struct {
	beego.Controller
}

type ST_Jsapi_Interface struct {
	AppId     string   `json:"appid"`
	Timestamp string   `json:"timestamp"`
	NonceStr  string   `json:"nonceStr"`
	Signature string   `json:"signature"`
	JsApiList []string `json:"jsApiList"`
}

type ST_JsApiRet struct {
	ST_Jsapi_Interface
	Token string
}

type ST_ParaUrl struct {
	Url string `json:"url"`
}

var g_wechat_token string = ""
var g_wechat_jsapi_ticket string = ""
var g_jsConfig ST_Jsapi_Interface
var g_lock sync.Mutex
var g_jString string = ""

func JsJdkInit() {

	config := JsConfig.GetConfig()
	if config == nil {
		log.Fatalln("JsConfig.GetConfig is nil")
		return
	}
	accessPath := config.WxJsApi.WeChatAccessToken
	ticketPath := config.WxJsApi.WeChatJsapiTicket
	g_jsConfig.AppId = config.WxJsApi.WeChatAppId
	apiList := config.WxJsApi.WeChatJsapiList

	//accessPath, err := JsConfig.GetConfigString([]string{"WxJsApi", "WeChatAccessToken"})
	//if err != nil {
	//	log.Fatalln(err)
	//}
	//
	//ticketPath, err := JsConfig.GetConfigString([]string{"WxJsApi", "WeChatJsapiTicket"})
	//if err != nil {
	//	log.Fatalln(err)
	//}
	//
	//g_jsConfig.AppId, err = JsConfig.GetConfigString([]string{"WxJsApi", "WeChatAppId"})
	//if err != nil {
	//	log.Fatalln(err)
	//}
	//
	//apiList, err := JsConfig.GetConfigString([]string{"WxJsApi", "WeChatJsapiList"})
	//if err != nil {
	//	log.Fatalln(err)
	//}
	g_jsConfig.JsApiList = strings.Split(apiList, ",")

	go wechat_token_coolie(accessPath, ticketPath)

	JsHttp.WhiteHttps("/wxjsapi", doWxJsapi)
}

func JsJdkInit_Unsafe() {

	config := JsConfig.GetConfig()
	if config == nil {
		log.Fatalln("JsConfig.GetConfig is nil")
		return
	}
	accessPath := config.WxJsApi.WeChatAccessToken
	ticketPath := config.WxJsApi.WeChatJsapiTicket
	g_jsConfig.AppId = config.WxJsApi.WeChatAppId
	apiList := config.WxJsApi.WeChatJsapiList

	//accessPath, err := JsConfig.GetConfigString([]string{"WxJsApi", "WeChatAccessToken"})
	//if err != nil {
	//	log.Fatalln(err)
	//}
	//
	//ticketPath, err := JsConfig.GetConfigString([]string{"WxJsApi", "WeChatJsapiTicket"})
	//if err != nil {
	//	log.Fatalln(err)
	//}
	//
	//g_jsConfig.AppId, err = JsConfig.GetConfigString([]string{"WxJsApi", "WeChatAppId"})
	//if err != nil {
	//	log.Fatalln(err)
	//}
	//
	//apiList, err := JsConfig.GetConfigString([]string{"WxJsApi", "WeChatJsapiList"})
	//if err != nil {
	//	log.Fatalln(err)
	//}
	g_jsConfig.JsApiList = strings.Split(apiList, ",")

	go wechat_token_coolie(accessPath, ticketPath)

	JsHttp.WhiteHttp("/wxjsapi", doWxJsapi)
}

func wechat_token_coolie(accessPath, ticketPath string) {
	for {
		response, e := http.Get(accessPath)
		if e != nil {
			JsLogger.Error(e.Error())
			if response != nil {
				response.Body.Close()
			}
		} else {
			b := make([]byte, 2048)
			n, _ := response.Body.Read(b)

			var token ST_WeChat_AccessToken
			json.Unmarshal(b[:n], &token)
			g_wechat_token = token.Access_token

			// fmt.Printf("token=%s", g_wechat_token)
			JsLogger.Info("token=%s", g_wechat_token)

			ticket_path := ticketPath + "?access_token=" + g_wechat_token + "&type=jsapi"

			update_jsapi_ticket(g_wechat_token, ticket_path)
		}

		time.Sleep(time.Hour)
	}
}

func update_jsapi_ticket(token, ticketPath string) {
	response, e := http.Get(ticketPath)

	if e != nil {
		JsLogger.Error(e.Error())
		response.Body.Close()

	} else {
		b := make([]byte, 2048)
		n, _ := response.Body.Read(b)

		var ticket ST_WeChat_Jsapi_Ticket
		json.Unmarshal(b[:n], &ticket)
		g_wechat_jsapi_ticket = ticket.Ticket

	}
}

func buildSignature(url string) {

	config := JsConfig.GetConfig()
	if config == nil {
		return
	}
	g_jsConfig.NonceStr = config.WxJsApi.WeChatNoncestr

	//var e error
	//g_jsConfig.NonceStr, e = JsConfig.GetConfigString([]string{"WxJsApi", "WeChatNoncestr"})
	//if e != nil {
	//	Error(e.Error())
	//}

	now := time.Now().Nanosecond()
	timestamp := strconv.Itoa(now)
	g_jsConfig.Timestamp = timestamp

	str := "jsapi_ticket="
	str += g_wechat_jsapi_ticket
	str += "&noncestr="
	str += g_jsConfig.NonceStr
	str += "&timestamp="
	str += timestamp
	str += "&url="
	str += url

	//产生一个散列值得方式是 sha1.New()，sha1.Write(bytes)，然后 sha1.Sum([]byte{})。这里我们从一个新的散列开始。
	h := sha1.New()
	//写入要处理的字节。如果是一个字符串，需要使用[]byte(s) 来强制转换成字节数组。
	h.Write([]byte(str))
	//这个用来得到最终的散列值的字符切片。Sum 的参数可以用来dui现有的字符切片追加额外的字节切片：一般不需要要。
	g_jsConfig.Signature = fmt.Sprintf("%x", string(h.Sum(nil)))

}

func RegisterJsApi(url string) {
	JsHttp.Http(url, doWxJsapi)
}

func doWxJsapi(session *JsHttp.Session) {

	var req_url ST_ParaUrl
	e := session.GetPara(&req_url)
	if e != nil {
		JsLogger.Error(e.Error())
		session.Forward("1", e.Error(), nil)
		return
	}

	buildSignature(req_url.Url)

	ret := &ST_JsApiRet{}
	ret.AppId = g_jsConfig.AppId
	ret.JsApiList = g_jsConfig.JsApiList
	ret.NonceStr = g_jsConfig.NonceStr
	ret.Signature = g_jsConfig.Signature
	ret.Timestamp = g_jsConfig.Timestamp
	ret.Token = g_wechat_token

	session.Forward("0", "success", ret)
}
