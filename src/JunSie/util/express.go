package util

//物流
import (
	"encoding/json"
	"log"
	"net/http"

	"JsGo/JsConfig"
	"JsGo/JsLogger"
	"crypto/md5"
	"fmt"
	"io/ioutil"
	"strings"
)

type ExpressFlow struct {
	Time    string `json:"time"`    //时间，原始格式
	Ftime   string `json:"ftime"`   //格式化后时间
	Context string `json:"context"` //内容
}

type ExpressInfo struct {
	Message       string        `json:"message"`   //消息体，请忽略
	ExpressNumber string        `json:"nu"`        //快递单号
	Ischeck       string        `json:"ischeck"`   //是否签收标记，请忽略，明细状态请参考state字段
	Condition     string        `json:"condition"` //快递单明细状态标记，暂未实现，请忽略
	ExpressName   string        `json:"com"`       //快递公司编码,一律用小写字母
	Status        string        `json:"status"`    //通讯状态，请忽略
	State         string        `json:"state"`     //快递单当前签收状态，包括0在途中、1已揽收、2疑难、3已签收、4退签、5同城派送中、6退回、7转单等7个状态，其中4-7需要另外开通才有效
	Data          []ExpressFlow `json:"data"`      //物流信息
}

var Express_Url string
var Express_Customer string
var Express_Key string

func init() {
	var e error
	Express_Url, e = JsConfig.GetConfigString([]string{"Express", "Url"})
	if e != nil {
		log.Fatalln(e.Error())
	}
	Express_Customer, e = JsConfig.GetConfigString([]string{"Express", "Customer"})
	if e != nil {
		log.Fatalln(e.Error())
	}
	Express_Key, e = JsConfig.GetConfigString([]string{"Express", "Sign"})
	if e != nil {
		log.Fatalln(e.Error())
	}
}

//查询订单物流
func QueryExpress(Name, Number string) (*ExpressInfo, error) {
	type Para struct {
		Com string `json:"com"`
		Num string `json:"num"`
	}
	st := Para{
		Com: Name,
		Num: Number,
	}
	b, e := json.Marshal(&st)
	if e != nil {
		return nil, e
	}
	para := string(b)
	str := para + Express_Key + Express_Customer
	has := md5.Sum([]byte(str))
	Sign := fmt.Sprintf("%x", has)
	Sign = strings.ToUpper(Sign)
	url := Express_Url
	url += "?customer="
	url += Express_Customer
	url += "&sign="
	url += Sign
	url += "&param="
	url += para

	JsLogger.Error("usr=====%s", url)

	resp, e := http.Get(url)
	defer resp.Body.Close()
	if e != nil {
		JsLogger.Error(e.Error())
		return nil, e
	}
	by, err := ioutil.ReadAll(resp.Body)
	JsLogger.Error(string(by))
	if err != nil {
		JsLogger.Error(err.Error())
		return nil, err
	}
	data := &ExpressInfo{}
	if err := json.Unmarshal(by, data); err != nil {
		JsLogger.Error(err.Error())
		return nil, err
	}
	return data, nil
}
