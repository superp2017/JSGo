package main

type CancleReq struct {
	OrderCode string `json:"orderCode"` //订单编号（创建订单接口返回值） "orderCode": "971271114707697664",
	OrderExt  string `json:"orderExt"`  //创建订单返回orderExt"orderExt": "735430230738468864",UUID:
	PolicyNo  string `json:"policyNo"`  //保单号"policyNo": "8G2012216201800000000087",
	ValidDate string `json:"validDate"` //批改生效时间：须晚于当前时间" validDate": "2018-03-07 16:00:00"
	//yyyy-MM-dd hh:mm:ss
} //请求报文体data

type SurrenderData struct {
	CancleInsuranceReq CancleReq `json:"cancelInsuranceReq"` //请求报文体data
}

type Surrender struct {
	InterfaceCode   string        `json:"interfaceCode"`   //接口标识(CreateOrder) "interfaceCode": "CreateOrder",
	RequestTime     string        `json:"requestTime"`     //请求时间yyyy-MM-dd HH:mm:ss "requestTime": "2018-02-01 14:32:53",
	DataSource      string        `json:"dataSource"`      //数据来源(易安分配)  "dataSource": "O- SHMK",
	AgrtCode        string        `json:"agrtCode"`        //协议号(易安分配) "agrtCode": " 03430002003000 ",
	OutBusinessCode string        `json:"outBusinessCode"` //保险订单号(第三方平台订单号，用于幂等性校验) "outBusinessCode": "1231212312312",
	Data            SurrenderData `json:"data"`
}

type XMSurrender struct {
	UID       string    //用户id
	ID        string    //退保id
	Cancle    Surrender //退保结构
	CreatDate string    //创建时间
}

func newSurrender(sur *XMSurrender) {

}
