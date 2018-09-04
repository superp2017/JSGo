package main

import (
	"JsGo/JsAliPay/alipay"
	"JsGo/JsHttp"
	"fmt"
	"io/ioutil"
	"log"
	"time"
)

var (
	appID     = "2017082408349309"
	partnerID = "2088221180975634"

	// // RSA2(SHA256)
	// aliPublicKey = []byte(`-----BEGIN PUBLIC KEY-----
	// 	MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAuqBFkVdN/quRFMGq0yQKpFTcCUPdXpcHfqESow3+ZPse8srWCzDhWONe8yMB3FBYy9xfhdcpL0bit1twpWHN/iE/HCEg/0H4n5Q5XJIqBB6g55IIoGkowgGm+6mwTMBdzMIhpo+96d0owusNWAuxfRZ2psJ5ocGcGiGRRx6cY5hS8IkaiOe5JE743dvEGzwZ03ZylC7x2BSLDvpUQa4BHOCgPQs4TcDd1CuSC4oUdbvc7pfYFaWJb3sQ875DA02fNHX/53Lfg1+mstZfa+i1NMlgZRnESMjkmgTmYELQ6PWQ6VXakP30rxaIxqdZtHpHs1jIGJC8l1m6yLVDKvqd4QIDAQAB
	// 	-----BEGIN PUBLIC KEY-----`)

	// publicKey = []byte(`-----BEGIN PUBLIC KEY-----
	// 	MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA2Keyi6ZFDxdBp/T+mkMfwa1V6qxXu4c3jmXqAps7nSUpjQniZhCoRz4fLyLgPJ6WQ1BA+ltmf4jdnWfIKIBESTeuKuhY5fhEWovaG9noder4fVOKO8vjsBwCBBchKpa2DqucvTuNb9Mn9+c1Ga3EBTd2+UfRHIAnkgENM3Vui7VbNHm3JAo0rWIYNQKn3i/wmvPiHfSxp25mPmaQi0p6tw5HnfNfXkwcCGm7u3gcwwSNtzYtIb9IK84ImTl638HSK7E11y9TEyF6PRhvdxCUNHmBaJfEoZxhlf4shjBYyTULU3QLA9JTIxSIB8Hdpq72Z67K2dx9BbVEUJtZ1dFsBwIDAQAB
	// 	-----END PUBLIC KEY-----
	// 	`)

	// privateKey = []byte(`-----BEGIN RSA PRIVATE KEY-----
	// 	MIIEpAIBAAKCAQEA2Keyi6ZFDxdBp/T+mkMfwa1V6qxXu4c3jmXqAps7nSUpjQniZhCoRz4fLyLgPJ6WQ1BA+ltmf4jdnWfIKIBESTeuKuhY5fhEWovaG9noder4fVOKO8vjsBwCBBchKpa2DqucvTuNb9Mn9+c1Ga3EBTd2+UfRHIAnkgENM3Vui7VbNHm3JAo0rWIYNQKn3i/wmvPiHfSxp25mPmaQi0p6tw5HnfNfXkwcCGm7u3gcwwSNtzYtIb9IK84ImTl638HSK7E11y9TEyF6PRhvdxCUNHmBaJfEoZxhlf4shjBYyTULU3QLA9JTIxSIB8Hdpq72Z67K2dx9BbVEUJtZ1dFsBwIDAQABAoIBAGztc5FO2W3K7ZG6Vb/Ne8vukEHawIliZIZNqygAUCqkIo3sqE1UlqarDxat3DveKAJT+TdmeNQqRfH72eCzDKIbQpAEHZ4SApvbsJ9MIdoXdzPbqqaBzkoe6syPrHczKvmZQsbJiETuzyuOrV3pxaIxzrlqaDKCJGL98Fss8ZV/fdubbXRKkOnaCiyaImRiH6cbh2FriW5isWBNW542u9WUpiWGVlDQ5nxj5xNSuKPAzea7UClnCk+RCM5xrKVTu57pZh2sBZpJbOwcUY8o+geU9xrJx2bthBrQ7NMBpGh58EkxCjorvJaioqVbmElxL63Z9j8lRDl/pDH6mBMHvgECgYEA+O+6oXcqHGJa8ggAbOZUmWDQkx/wAkYG0915McW6mruEC7YzOejCXCxmj8wh/jASFGRWaUxiuAqxfIt+PVbWK8McFcRB7cIxDmkNVt0rueSOF2NdladhhfVOpO05J3zUwo9XYGoxuVszM0nmQw2PrBQyu1jlSO2//mQ3lW5Z14ECgYEA3s16GaPM9Rg9dOfIC7YD7XIVZcuyZqzv7PoeFapd1wekD7o363mtboUV/j3a4eBNmHOK4rv+OtUt0r5DRe75yGrvvS1neTARGu/84XlSGJo7wmgL6GUs9VN7PqjPupTwak4puTJzYJB9o+1HsZ5Q9XYTl9jjKdj/EJzzJG+sR4cCgYEAqEX9BZq0550A1yzbhMGqHEgalel350GI6fyDKUb83g21s+kE9bdGcuI8riWSMO4zun8c/m75KGlqEsOEoVgqzEhGmtwgqOSlHpWaw8YcAbvi5SJxJ3GO9eudrtUA1pWGiMI2kWEXnbFtidUBhwAKx4qbxJLR4xt7ti3ueN+wcYECgYA/pe91l6ebdNtJpFUvk0W39VlLhU9nqYu45RLnGY5JOXOS0p3a9R2obviDcuQulsdT/93zO8U6xV+bzqKlPcm5iWMHZgsjQaoBSgGx39imEplzxglw0EZxpvGUSuFc6eNsWvvsXg87zMs3ozdR9GooVRzvyhPLBqSG+G81P3m1zQKBgQCGloZ/ZA5tMa3eiyewe236vOLR+K1lYJPEUlMNUNhj26Jw41cytvkm42loY758uRZlTySAWNJCISNxc+0DY1U0TiPNUQOCCP8Lu6mItj2ijdM1EHXlNQEWYsL1Tcl62kzTrsB2OOcYW5uvRYSpZqMOEYnjJPG3/lam5g1RDPl5iw==
	// 	-----END RSA PRIVATE KEY-----`)

	// RSA2(SHA256)
	aliPublicKeyFile = string("cert/alipay/ali_public_key.pem")
	publicKeyFile    = string("cert/alipay/app_public_key.pem")
	privateKeyFile   = string("cert/alipay/app_private_key.pem")

	// RSA(SHA1)
//	aliPublicKey = []byte(`-----BEGIN PUBLIC KEY-----
//MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQDIgHnOn7LLILlKETd6BFRJ0Gqg
//S2Y3mn1wMQmyh9zEyWlz5p1zrahRahbXAfCfSqshSNfqOmAQzSHRVjCqjsAw1jyq
//rXaPdKBmr90DIpIxmIyKXv4GGAkPyJ/6FTFY99uhpiq0qadD/uSzQsefWo0aTvP/
//65zi3eof7TcZ32oWpwIDAQAB
//-----END PUBLIC KEY-----
//`)
//
//	publicKey = []byte(`-----BEGIN PUBLIC KEY-----
//MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQC/27vaZOkiSI+7/I0PQXHLWV+l
//uPhXd2sJIT3YnjUSDbW1Lo6HES0yDP/LOAdVHfzxU09+BnKJbSHAsmBuf/ZQej5y
//lYi7KUNekTf9zRiaT5mrt2T6GNUptbF/o5Ew4dIAdqvbe1+KQZhzkgoJ1o6uNqFH
//jVkE05TcQ7NQYr42JwIDAQAB
//-----END PUBLIC KEY-----`)
//
//
//	privateKey = []byte(`-----BEGIN RSA PRIVATE KEY-----
//MIICXAIBAAKBgQC/27vaZOkiSI+7/I0PQXHLWV+luPhXd2sJIT3YnjUSDbW1Lo6H
//ES0yDP/LOAdVHfzxU09+BnKJbSHAsmBuf/ZQej5ylYi7KUNekTf9zRiaT5mrt2T6
//GNUptbF/o5Ew4dIAdqvbe1+KQZhzkgoJ1o6uNqFHjVkE05TcQ7NQYr42JwIDAQAB
//AoGAAgpbOBpkpFmzNaOc+HGQvUHpE4EhGwUJHWK+HqSoGdYNfazOFT+ohGTA/69L
//3Krh+ssRCF0XMMD5X+xFvEceHn47yr3TNJeArsT61UGORm4z0bWPwurjQx884t56
//dXY2X4NnEHPJA1AlphWASZu4h8TkBzsMhfmfJQDURBuWn7ECQQD+x13z+baTCShv
//BMKrB+fVZa/yfVx3Mk2m5COn3EosF/+SUxPUONav8b7MqNaR20pSJBxmqpybKP5I
//BbtO7FOpAkEAwMcomveKwRlsP7qse30NY7TvJDoUZPezGutwDNlI5YjjOVh3RaYd
//SgtCHzqYRQRhiL3ESDHjNXBpj/ayJYxdTwJAIr859w41cjQriYiSrBS174qgxmeG
//dtMrd/lhS4FltEHJn0EpUSY3UWOc6/iS2u2XY0B9hxr5pMegdl4hv4/HkQJANCxy
//j+ZZFkPUKTdTgSRqIEcSxeI2LNFhFvMLY17XPNAcdyO7PA1mNejwH1WTanJyFzkM
//y2E9FfRzjXP96O2hPwJBAJKUyGfGQXVPqbCoYWaX/Bqj7ok8dal74OCRbKp9WrBe
//FOEq/sfp2vYGaCw9uyczDwRKcliKibgAEPmbZ1ToKt0=
//-----END RSA PRIVATE KEY-----
//`)
)

var client *alipay.AliPay

var g_aliPay map[string]*alipay.AliPayTradeAppPay

func Alinit() {
	privateKey, err := ioutil.ReadFile(privateKeyFile)
	if err != nil {
		log.Fatalln(err.Error())
	}

	publicKey, err := ioutil.ReadFile(publicKeyFile)
	if err != nil {
		log.Fatalln(err.Error())
	}

	aliPublicKey, err := ioutil.ReadFile(aliPublicKeyFile)
	if err != nil {
		log.Fatalln(err.Error())
	}

	client = alipay.New(appID, partnerID, publicKey, privateKey, true)
	client.AliPayPublicKey = aliPublicKey

	g_aliPay = make(map[string]*alipay.AliPayTradeAppPay)

	JsHttp.Http("/newaliorder", newAliOrder)
	JsHttp.Http("/newaliorderex", newAliOrderEx)
	JsHttp.Http("/newaliordercb", newAliOrderCb)
	// exitCode := m.Run()
	// os.Exit(exitCode)
}

func newAliOrder(session *JsHttp.Session) {
	fmt.Println("========== TradeAppPay ==========")
	var p = alipay.AliPayTradeAppPay{}
	p.NotifyURL = "http://test.junsie.cn:9568/newaliordercb"
	p.Body = "body"
	p.Subject = "商品标题"
	p.OutTradeNo = fmt.Sprintf("%d", time.Now().Nanosecond())
	p.TotalAmount = "0.15"
	p.ProductCode = "p_1010101"

	order, err := client.TradeAppPay(p)
	// fmt.Println(order, err)

	if err != nil {
		fmt.Println(err.Error())
		session.Forward("1", err.Error(), nil)
		return
	}

	g_aliPay[p.OutTradeNo] = &p
	session.Forward("0", "success", order)
}

func newAliOrderEx(session *JsHttp.Session) {
	fmt.Println("========== TradeAppPay ==========")
	var p = alipay.AliPayTradeAppPay{}
	p.NotifyURL = "http://test.junsie.cn:9568/newaliordercb"
	p.Body = "body"
	p.Subject = "商品标题"
	p.OutTradeNo = fmt.Sprintf("%d", time.Now().Nanosecond())
	p.TotalAmount = "0.15"
	p.ProductCode = "p_1010101"

	orders, err := client.TradeAppPayEx(p)
	// fmt.Println(order, err)

	if err != nil {
		fmt.Println(err.Error())
		session.Forward("1", err.Error(), nil)
		return
	}

	session.Forward("0", "success", orders)
}

func newAliOrderCb(session *JsHttp.Session) {

	session.Forward("1", "success", nil)
}
