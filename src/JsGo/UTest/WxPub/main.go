package main

import (
	"JsGo/JsExit"
	"JsGo/JsHttp"
	"JsGo/JsLogger"
	"JsGo/JsMobile"
	"JsGo/JsQiniu"
	"JsGo/JsWeChat/JsAppAuth"
	"JsGo/JsWeChat/JsWechatAuth"
	"JsGo/JsWeChat/JsWechatJsApi/JsJdk"
	"log"

	"github.com/chanxuehong/wechat/mp/user/oauth2"
)

func exitCb() int {
	JsHttp.Close()
	return 0
}

func authCb(user *oauth2.UserInfo, session *JsHttp.Session) {
	log.Println(user)
}

func appUserInfoCb(user *oauth2.UserInfo, session *JsHttp.Session) {
	log.Println(user)
}

func ExchangeToken(s *JsHttp.Session) {
	type Para struct {
		OpenId  string
		UnionId string
	}

	para := &Para{}
	e := s.GetPara(para)
	if e != nil {
		s.Forward("1", e.Error(), "")
		return
	}

	s.MarkSession()
	s.Forward("0", "success", "")
}

func nothing(s *JsHttp.Session) {
	s.Forward("0", "success", "")
}

func main() {
	JsHttp.WhiteHttp("/exchagetoken", ExchangeToken)
	JsHttp.Http("/nothing", nothing)

	JsHttp.EnableGet()
	JsHttp.EnableSession()
	JsWechatAuth.WxauthInit(authCb)
	JsAppAuth.AppInit(appUserInfoCb)
	JsExit.RegisterExitCb(exitCb)
	JsJdk.JsJdkInit()
	JsQiniu.QiniuInit()
	JsMobile.AlidayuInit()

	Alinit()
	InitOrder()

	JsLogger.Console(true)

	JsHttp.Run()

}

// package main

// import (
// 	"encoding/pem"
// 	"fmt"
// 	"io/ioutil"
// )

// var pk = []byte(`-----BEGIN RSA PRIVATE KEY-----
// 	MIIEogIBAAKCAQEAmeQtyv2QiL1Pzj3h1XWXmtTovxFJhoaenXZo4rD59kg/21df
// 	1u5JfEhPz1MHzqMF5ycyn6g6cN9idGjYkbM5tMqcnkWc/UeaLte6tQE1PjSJbYZQ
// 	vuXhBVb/uOyFwGe2BbPnHk9kQalhja6dOk5iMzAii4LpQlLVH4+pUNV/Ina5LggP
// 	86U0wMbjYMXIVkJ4ws5Lz/gmBngjdWlZu2O9eoAvlazwMlEce2Z/HqXx89gbKvdi
// 	ZyskDyFeBIGfEu4+HSDXoj1Oqdz503jA+5H4fOZkOiHP0PbJVuYGFkT/IbYyKyuG
// 	m2RR/xKpI7VHfuo9Axp4Qe1jVSCpZZPoJrMIOwIDAQABAoIBACBa9DHQnBnTy7qu
// 	EhRCNAzOpNy/MrlBBopOwgCfev6H9D0WosTatsKVpYoOh/6vEeemuyMMSLVAkj+t
// 	Z5NCnmhfjQxN0JMEAevFWbECvwIEI7zOV527UVNBFmT3/asVYxTR3U7nHEod5a/X
// 	PaFrp4Pho/G3JXnXVo3bM5sWODOzdPwmgzdqUG7of3kplzbixUSrJ2llUmlfsZqK
// 	ElbPB2chLn0FS6AJAxvPo1GBiMIMsw+98x/xmLVlH8erHghFEObodivyXiS/6Gzb
// 	ayoxEn4IQ27i9ZTBFXmGYwVO2ITuSHAWs9e75PqwdGrMKN01XPNGpHmKnuR3AlD8
// 	l1p4XxECgYEAzPOWyphUvaHye7R6chB1pLDwsxLVF/RLI1/cG9+DK7+FL9NxkVWv
// 	IccZBv6nZN9M0hFh1x6pqvTXVQd4A16GtuEa/ARlkcaTfIXj/XCEbONlT7GLwd3i
// 	KFQ6csG5VCg1YbY1d5vNlPycEZQXTauX44j0UktP53/7VhyegI31JBkCgYEAwDjR
// 	HPegzcAKezl00SBwIooQWlvhfkZ0guUCxk4mWRlCDV+U2oR7bNthAw5Y1xHN8eaH
// 	Y6y+n3yvyf2rzvT7sUw8T7pLtM3FicYf4gneZc7ZW+SGvJjxBv9wpsBAIP4Yisdb
// 	PbdqD6gfxTk2y18DT7+a2Py1EiwY6ysal/AkeXMCgYBcjfSe8UPzj1sN+mcBc+Vs
// 	xmsss2iANNZp1zRzcfCupQLkojw7QdKhEmR/AClgKGdsxmTE3RgKGB/WSlUsUFfN
// 	5sJk5SdpOaAJL/3RyipDcj2iS6+tkSI8zCzI/itPkgjpY3up1DZ2/c0NMy+C5+bj
// 	3klXkKM5DFbYgHwj2ffGoQKBgBAZciI790LkT7xsXoVZcyrhZ2c6BNPfsMh5x9a2
// 	Gu4heG/ITp5StEe0xBZOcFBrFFWrWjGV+U1AUzTWwzoNOLtryC1hTA/zoBTe/DKh
// 	Yvgh8ACLTmGjaaSNZnEA7x4UShfthI3Ru9dd3HNXTGiSJ6PZR23fFIdWHCwuKwcI
// 	vPVTAoGAVqAhOWsmf9vrh5BIiVEwbjTfjKrG/VJdOM/lL5IergwaClHtq4TYIdGA
// 	jHmA1LuTZW6b+yTSTD0cobRMWehAX0ZZrlNwib0O4e5k/FVrBfNlJ6emHWJV+Vcj
// 	OyunKbzdw9AsBubrzsABpPrgMtEL8v5YvNJzmECRZ+AwgzzLgBI=
// 	-----END RSA PRIVATE KEY-----`)

// func main() {
// 	// var block *pem.Block
// 	// block, _ = pem.Decode(pk)
// 	// if block == nil {
// 	// 	fmt.Println("private key error")
// 	// }

// 	//var cert tls.Certificate
// 	//加载PEM格式证书到字节数组
// 	certPEMBlock, err := ioutil.ReadFile("ali_public_key.pem")
// 	if err != nil {
// 		fmt.Println(err.Error())
// 		return
// 	}
// 	//获取下一个pem格式证书数据 -----BEGIN CERTIFICATE-----   -----END CERTIFICATE-----
// 	certDERBlock, _ := pem.Decode(certPEMBlock)
// 	if certDERBlock == nil {
// 		fmt.Println("private key error")
// 	}
// }
