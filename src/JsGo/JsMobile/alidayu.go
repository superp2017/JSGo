package JsMobile

import (
	"JsGo/JsConfig"
	"JsGo/JsHttp"
	. "JsGo/JsLogger"
	"JsGo/JsNet"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/coocood/freecache"
)

var g_smscache *freecache.Cache

var g_rand_chan chan int

var g_sms_cfg map[string]string

//兼容老接口
func AlidayuInit() {

	g_smscache = freecache.NewCache(32 * 1024 * 1024) // 32MB

	g_rand_chan = make(chan int)

	var err error
	g_sms_cfg, err = JsConfig.GetConfigMap([]string{"MobileVerify"})
	if err != nil {
		log.Fatalln(err)
	}

	go randCoolie()

	JsNet.WhiteHttp("/alidayu", alidayu_old)
}

func alidayu_old(session *JsNet.StSession) {
	type Para struct {
		SignName string
		Mobile   string
		SmsCode  string
		Expire   int
	}

	para := &Para{}
	e := session.GetPara(para)
	if e != nil {
		session.Forward("1", e.Error(), nil)
		return
	}

	ComJsMobileVerify(para.SignName, para.Mobile, para.SmsCode, "a", para.Expire, nil)
	session.Forward("0", "success", "")
}

func NewAlidayuInit() {

	g_smscache = freecache.NewCache(32 * 1024 * 1024) // 32MB

	g_rand_chan = make(chan int)

	var err error
	g_sms_cfg, err = JsConfig.GetConfigMap([]string{"MobileVerify"})
	if err != nil {
		log.Fatalln(err)
	}

	go randCoolie()

	JsHttp.WhiteHttps("/alidayu", alidayu)
}

func alidayu(session *JsHttp.Session) {
	type Para struct {
		SignName string
		Mobile   string
		SmsCode  string
		Expire   int
	}

	para := &Para{}
	e := session.GetPara(para)
	if e != nil {
		session.Forward("1", e.Error(), nil)
		return
	}
	Error("alidayu=%v", para)

	ComJsMobileVerify(para.SignName, para.Mobile, para.SmsCode, "a", para.Expire, nil)
	session.Forward("0", "success", "")
}

func randCoolie() {
	rand_gen := rand.New(rand.NewSource(int64(time.Now().Nanosecond())))
	for {
		g_rand_chan <- rand_gen.Int()
	}
}

func getCode() string {
	code := <-g_rand_chan
	ret := fmt.Sprintf("%06d", code%1000000)
	return ret
}

func ComJsMobileVerify(signName, mobile, smscode, t string, expire int, par map[string]string) {

	code := ""
	if par != nil {
		code = par["code"]
	}

	if code == "" {
		code = getCode()

		g_smscache.Set([]byte(mobile), []byte(code), expire)
	}

	if signName == "" {
		signName = g_sms_cfg["SignName"]
	}

	para := "?appkey="
	para += g_sms_cfg["AppKey"]
	para += "&secretkey="
	para += g_sms_cfg["SecretKey"]
	para += "&signname="
	para += signName
	para += "&mobile="
	para += mobile
	para += "&smscode="
	para += smscode
	para += "&type="
	para += t

	for k, v := range par {
		para += "&"
		para += k
		para += "="
		para += v
	}

	if par == nil || par["code"] == "" {
		para += "&code="
		para += code
	}

	Info(g_sms_cfg["VUrl"] + para)
	response, e := http.Get(g_sms_cfg["VUrl"] + para)
	if e != nil {
		Error(e.Error())
		return
	}
	b := make([]byte, 2048)
	response.Body.Read(b)
	n, _ := response.Body.Read(b)

	defer response.Body.Close()
	if e != nil {
		Error("verify %s error, rsp:%s\n", mobile, string(b[:n]))
	}
}

func VerifySmsCode(mobile, code string) bool {
	vCode, e := g_smscache.Get([]byte(mobile))

	if e == nil && string(vCode) == code {
		return true
	} else {
		return false
	}
}

////////////////////////////////////////////////////////////////////////////////
//
//新接口
//                                                                            //
////////////////////////////////////////////////////////////////////////////////

// const v5_url = "http://www.api.zthysms.com/sendSms.do"
// const v5_username = "shxyhy"
// const v5_password = "9BApAi"

// func verify_ex(code, product, mobile string) {
// 	tkey := time.Now().Format("20060102150405")

// 	md5Ctx := md5.New()
// 	md5Ctx.Write([]byte(v5_password))
// 	cipherStr := md5Ctx.Sum(nil)

// 	md5Ctx = md5.New()
// 	md5Ctx.Write([]byte(hex.EncodeToString(cipherStr) + tkey))

// 	pwd := hex.EncodeToString(md5Ctx.Sum(nil))
// 	para := "?username=" + v5_username
// 	para += "&tkey=" + tkey
// 	para += "&password=" + pwd
// 	para += "&mobile=" + mobile
// 	para += "&content=hello"

// 	response, e := http.Get(v5_url + para)
// 	b := make([]byte, 2048)
// 	response.Body.Read(b)

// 	defer response.Body.Close()
// 	if e != nil {
// 		b := make([]byte, 2048)
// 		n, _ := response.Body.Read(b)
// 		g_log.Error("verify %s error, rsp:%s\n", mobile, string(b[:n]))
// 	}
// }

// func RegisterAuth_ex(mobile, product string, expire int) {
// 	code := getCode()
// 	g_smscache.Set([]byte(mobile), []byte(code), expire)
// 	g_log.Info("---------------------------------------------------------Code=%s\n", code)
// 	verify_ex(code, product, mobile)
// }
