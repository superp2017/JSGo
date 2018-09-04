package JsConfig

import (
	"log"
	"sync"
)

type St_WxJsApi struct {
	WeChatAccessToken  string
	WeChatJsapiTicket  string
	WeChatAppId        string
	WeChatNoncestr     string
	WeChatSecret       string
	WeChatOAuth2       string
	WeChatOAuth2Path   string
	WeChatRedirectHome string
	WeChatOriHome      string
	WeChatJsapiList    string
}

type St_Pay struct {
	WxAppId          string
	WxSecret         string
	WxMchId          string
	WxSecretKey      string
	WxPubPayCb       string
	WxPubPayUrl      string
	WxPubRefundUrl   string
	WxPubTransferUrl string
	WxPubSendredpack string
	WxSpbillCreateIp string
	CertFile         string
	KeyFile          string
	RootcaFile       string
}

type RunTimeConfig struct {
	WxJsApi St_WxJsApi
	WxMiniP St_Pay
	AppPay  St_Pay
	PubPay  St_Pay
}

var config *RunTimeConfig = nil
var mutext sync.Mutex

func InitConfig(c *RunTimeConfig) {
	mutext.Lock()
	defer mutext.Unlock()
	config = c
	if config == nil {
		log.Fatalln("InitConfig failed")
	}
}
func GetConfig() *RunTimeConfig {
	return config
}

//"Logger":{
//"Logger": "wxpub.log",
//"MaxSize": 1048576,
//"Console": "true"
//},
//
//"Redis":{
//"Ip": "139.196.108.155",
//"Port": "6379",
//"Password":"1qaz2wsx"
//},
//
//"MobileVerify":{
//"AppKey":"LTAIvgiaaRMJpz1h",
//"SecretKey":"YOBZTZSzMtmNSrmXrtQNeT7P9IbSzW",
//"VUrl":"http://test.junsie.com/alidayu/mobile/index.php"
//},
//
//"Net":{
//"Session":"true",
//"Get":"true"
//},
//
//
//
//"Http":{
//"Ip": "",
//"Listen": "9568"
//},
//
//"Https":{
//"Ip": "",
//"Listen": "9569",
//"Name": "MALL",
//"Pem":"cert/fullchain.pem",
//"Key":"cert/privkey.pem"
//},
//
//"Qiniu":{
//"Scope":"wxpub",
//"Bucket":"wxpub",
//"Domain":"http://resource.junsie.com/",
//"AK":"NLhOfHy6_Y3qZe7DJ2qW4fkZ4KPJTBSNkLonPjN6",
//"SK":"lrWcWN4IsnvXQt78vCPfgoRQv4bNqkMZ-pLoEKGR"
//},
//
//"WxJsApi":{
//"WeChatAccessToken":"https://api.weixin.qq.com/cgi-bin/token?grant_type=client_credential&appid=wx7c9376486e203c66&secret=556b824cd8ce64a1cb35057114f4e1e5",
//"WeChatJsapiTicket":"https://api.weixin.qq.com/cgi-bin/ticket/getticket",
//"WeChatAppId":"wx7c9376486e203c66",
//"WeChatNoncestr":"24394875hihfq0893i5hrmghf",
//"WeChatSecret":"556b824cd8ce64a1cb35057114f4e1e5",
//"WeChatOAuth2":"https://www.junsie.cn",
//"WeChatOAuth2Path":"/wxauthcb",
//"WeChatRedirectHome": "http://pub.junsie.cn",
//"WeChatOriHome":"https://pub.junsie.cn:9569",
//"WeChatJsapiList":"onMenuShareTimeline,onMenuShareAppMessage,onMenuShareQQ,onMenuShareWeibo,onMenuShareQZone,startRecord,stopRecord,onVoiceRecordEnd,playVoice,pauseVoice,stopVoice,onVoicePlayEnd,uploadVoice,downloadVoice,chooseImage,previewImage,uploadImage,downloadImage,translateVoice,getNetworkType,openLocation,getLocation,hideOptionMenu,showOptionMenu,hideMenuItems,showMenuItems,hideAllNonBaseMenuItem,showAllNonBaseMenuItem,closeWindow,scanQRCode,chooseWXPay,openProductSpecificView,addCard,chooseCard,openCard"
//
//},
//
//"WxMiniP":{
//"WxAppId":"wx11cdac22d7719783",
//"WxSecret":"ee059d475b23cfde69a77327b2018bd8",
//"WxMchId":"1317786501",
//"WxSecretKey":"3248v0n90vcm8u305uvn23m85ux40924",
//"WxPubPayCb":"http://pub.junsie.cn:9568/paidordersuccesscb",
//"WxPubPayUrl":"https://api.mch.weixin.qq.com/pay/unifiedorder",
//"WxPubRefundUrl":"https://api.mch.weixin.qq.com/secapi/pay/refund",
//"WxPubTransferUrl":"https://api.mch.weixin.qq.com/mmpaymkttransfers/promotion/transfers",
//"WxPubSendredpack":"https://api.mch.weixin.qq.com/mmpaymkttransfers/sendredpack",
//"WxSpbillCreateIp":"139.196.108.155",
//"CertFile":"cert/pub/apiclient_cert.pem",
//"KeyFile":"cert/pub/apiclient_key.pem",
//"RootcaFile":"cert/pub/rootca.pem"
//
//},
//"AppPay":{
//"WxAppId":"wx59b594e3ca3fd289",
//"WxSecret":"eb73bd344ef563ed12ff6920a84ec127",
//"WxMchId":"1317447701",
//"WxSecretKey":"09284m0cjhk9lzkq2c480t234v5f9r0o",
//"WxPubPayCb":"http://pub.junsie.cn:9568/paidordersuccesscb",
//"WxPubPayUrl":"https://api.mch.weixin.qq.com/pay/unifiedorder",
//"WxPubRefundUrl":"https://api.mch.weixin.qq.com/secapi/pay/refund",
//"WxPubTransferUrl":"https://api.mch.weixin.qq.com/mmpaymkttransfers/promotion/transfers",
//"WxPubSendredpack":"https://api.mch.weixin.qq.com/mmpaymkttransfers/sendredpack",
//"WxSpbillCreateIp":"139.196.108.155",
//"CertFile":"cert/app/apiclient_cert.pem",
//"KeyFile":"cert/app/apiclient_key.pem",
//"RootcaFile":"cert/app/rootca.pem"
//},
//
//
//"PubPay":{
//"WxAppId":"wx7c9376486e203c66",
//"WxSecret":"556b824cd8ce64a1cb35057114f4e1e5",
//"WxMchId":"1317786501",
//"WxSecretKey":"3248v0n90vcm8u305uvn23m85ux40924",
//"WxPubPayCb":"http://pub.junsie.cn:9568/paidordersuccesscb",
//"WxPubPayUrl":"https://api.mch.weixin.qq.com/pay/unifiedorder",
//"WxPubRefundUrl":"https://api.mch.weixin.qq.com/secapi/pay/refund",
//"WxPubTransferUrl":"https://api.mch.weixin.qq.com/mmpaymkttransfers/promotion/transfers",
//"WxPubSendredpack":"https://api.mch.weixin.qq.com/mmpaymkttransfers/sendredpack",
//"WxSpbillCreateIp":"139.196.108.155",
//"CertFile":"cert/pub/apiclient_cert.pem",
//"KeyFile":"cert/pub/apiclient_key.pem",
//"RootcaFile":"cert/pub/rootca.pem"
