package main

import (
	"JsGo/JsHttp"
	"JsGo/JsLogger"
	"JsGo/JsMobile"
	"JsGo/JsStore/JsRedis"
	"JunSie/util"
)

const (
	DB_Insurance = "Insurnce"
	DB_Surrender = "Surrender"
)

//受益人
type AcciInsured struct {
	DocType       string `json:"docType"`       // 证件类型个人证件类型 "docType":"01",
	Sex           string `json:"sex"`           // 性别代码参见性别代码 "sex":"01",
	CustomerFlag  string `json:"customerFlag"`  // 关系人标志参见关系人标志默认为1-投保人 customerFlag":"2",
	DocNo         string `json:"docNo"`         // 证件号码 "docNo":"330101198001010070",
	AppliRelation string `json:"appliRelation"` // 被保险人与投保人关系 参见关系人关系代码"appliRelation":"05",
	CustomerName  string `json:"customerName"`  // 受益人名称  "customerName":"李双",
	PhoneNo       string `json:"phoneNo"`       // 联系电话 "phoneNo":"13000000000"
}

type ItemAcc struct { //标的主信息
	NominativeInd   string        `json:"nominativeInd"`   // 受益比例（单位：%）一个受益人：默认100多个受益人：受益比例总和不能大于100"nominativeInd":"1",
	Quantity        string        `json:"quantity"`        // 人数 "quantity":"1",
	OccupationCode  string        `json:"occupationCode"`  // 职业代码参见职业类别代码如无职业信息则传默认值：0000000"occupationCode":"0000",
	AcciInsuredList []AcciInsured `json:"acciInsuredList"` // 被保人清单
}

type RiskDynamic struct {
	FieldAA string `json:"fieldAA"` // 医疗美容机构 //动态扩展信息Az ~ Az 根据不同产品具体定义各个字段的含义
	FieldAB string `json:"fieldAB"` // 医疗美容项目"fieldAB":"01"
}

type Customer struct { //投保人列表
	CustomerSameInd string `json:"customerSameInd"` //是否同投保人：1:是，0：否
	CustomerFlag    string `json:"customerFlag"`    // 关系人标志参见关系人标志默认为1-投保人"customerFlag":"1",
	CustomerName    string `json:"customerName"`    // 被保人姓名 "customerName":"毕沐阳",
	CustomerType    string `json:"customerType"`    // 投保人类型参见关系人代码 "customerType":"2",
	DocNo           string `json:"docNo"`           // 证件号码 "docNo":"110101198001017933",
	DocType         string `json:"docType"`         // 证件类型 个人证件类型"docType":"1",
	Email           string `json:"email"`           // 电子邮箱 "email":"2583970208@qq.com",
	PhoneNo         string `json:"phoneNo"`         // 联系电话 "phoneNo":"13544270502",
	Sex             string `json:"sex"`             // 性别代码参见性别代码 注：企业用户可空"sex":"01",
	CustomerAddress string `json:"customerAddress"` // 地址（企业用户必填：企业地址）"customerAddress":"beijingxxxx",
	ContactName     string `json:"contactName"`     // 联系人（企业用户必填） "contactName":"xxxx",
	OfficePhone     string `json:"officePhone"`     // 办公电话（企业用户必填，格式 区号-号码）"officePhone":"0730-8100016",
	ContactMobile   string `json:"contactMobile"`   // 联系电话（企业用户必填） "contactMobile":"13100000000"
}

type Order struct {
	StartDate       string        `json:"startDate"`       // 起保日期 "startDate":"2018-04-18 00:00:00"
	EndDate         string        `json:"endDate"`         // 截止日期 "endDate":"2018-04-18 00:00:00"
	Premium         string        `json:"premium"`         // 保费( 总保费 = 份数*人数*每人保费) 产品方案为定额时可为空"premium":"30"
	UWCount         int           `json:"uwCount"`         // 投保份数Integer（默认1份） "uwCount":1
	CustomerList    []Customer    `json:"customerList"`    // 投保人列表
	ItemAcciList    []ItemAcc     `json:"itemAcciList"`    // 标的主信息
	RiskDynamicList []RiskDynamic `json:"riskDynamicList"` // 意外险标的信息
}

type OrderReq struct {
	OrderList []Order `json:"orderList"`
}

type InsuranceData struct {
	CreateOrderReq OrderReq `json:"createOrderReq"`
}

type Insurance struct {
	Data            InsuranceData `json:"data"`            // 返回数据，异常时该数据可能为空 保单data
	AgrtCode        string        `json:"agrtCode"`        // 协议号(易安分配)   "agrtCode":"03430002003000"
	RequestTime     string        `json:"requestTime"`     // 请求时间yyyy-MM-dd HH:mm:ss"requestTime":"2016-07-19 09:17:13",
	DataSource      string        `json:"dataSource"`      // 数据来源(易安分配)  "dataSource":"O-SHMK",
	OutBusinessCode string        `json:"outBusinessCode"` // 保险订单号(第三方平台订单号，用于幂等性校验)  "outBusinessCode":"123456789017",
	InterfaceCode   string        `json:"interfaceCode"`   // 接口标识(CancelInsurance) "interfaceCode":"CreateOrder"
}

type XMInsurance struct {
	ID          string    //保单id
	UID         string    //用户id
	OrderID     string    //订单id
	OrderStatus string    //订单状态
	TimeStamp   int64     // 时间戳
	Policy      Insurance //保单
	Status      string    //状态
	CreatTime   string    //创建时间
}

//创建保单
func CreatInsurance(s *JsHttp.Session) {
	st := &XMInsurance{}
	if err := s.GetPara(st); err != nil {
		s.Forward("1", err.Error(), nil)
		return
	}
	///////////检查字段////////////////
	if err := checkInsurance(st); err != nil {
		s.Forward("1", err.Error(), nil)
		return
	}
	st.ID = util.IDer(DB_Insurance)
	st.Status = "0"
	st.CreatTime = util.CurTime()
	st.Policy.RequestTime = util.CurTime()
	st.Policy.OutBusinessCode = st.ID
	if err := JsRedis.Redis_hset(DB_Insurance, st.ID, st); err != nil {
		s.Forward("1", err.Error(), nil)
		return
	}
	go appendUserInsurance(st.UID, st.ID)
	s.Forward("0", "success", st)
}

//获取单个保单
func QueryInsurance(s *JsHttp.Session) {
	type Para struct {
		ID string
	}
	st := &Para{}
	if err := s.GetPara(st); err != nil {
		s.Forward("1", err.Error(), nil)
		return
	}
	if st.ID == "" {
		info := "QueryInsurance ID is empty\n"
		JsLogger.Error(info)
		s.Forward("1", info, nil)
		return
	}
	d, e := getInsurance(st.ID)
	if e != nil {
		s.Forward("1", e.Error(), nil)
		return
	}

	s.Forward("0", "success", d)
}

//获取用户的所有保单列表
func GetUserInsurance(s *JsHttp.Session) {
	type Para struct {
		UID string //用户ID
	}
	st := &Para{}
	if err := s.GetPara(st); err != nil {
		s.Forward("1", err.Error(), nil)
		return
	}
	if st.UID == "" {
		info := "GetUserInsurance UID is empty\n"
		JsLogger.Error(info)
		s.Forward("1", info, nil)
		return
	}
	list, err := getUserInsurance(st.UID)
	if err != nil {
		s.Forward("1", err.Error(), nil)
		return
	}
	s.Forward("0", "success", getMoreInsurance(list))
}

//更改保单的订单状态
func UpdateInsuranceOrder(ID, OrderID, Status string, timeStamp int64) error {
	d, e := getInsurance(ID)
	if e != nil {
		return e
	}
	d.OrderID = OrderID
	d.Status = Status
	d.OrderStatus = Status
	d.TimeStamp = timeStamp
	return JsRedis.Redis_hset(DB_Insurance, ID, d)
}

//获取单个保单
func getInsurance(ID string) (*XMInsurance, error) {
	st := &XMInsurance{}
	err := JsRedis.Redis_hget(DB_Insurance, ID, st)
	return st, err
}

//获取多个保单
func getMoreInsurance(IDs []string) []*XMInsurance {
	data := []*XMInsurance{}
	for _, v := range IDs {
		d, err := getInsurance(v)
		if err == nil {
			data = append(data, d)
		}
	}
	return data
}

//检查保单字段
func checkInsurance(data *XMInsurance) error {
	if data.UID == "" {
		return JsLogger.ErrorLog("checkInsurance failed,UID is empty\n")
	}
	//if data.Policy.InterfaceCode==""{
	//	return JsLogger.ErrorLog("checkInsurance failed,InterfaceCode is empty\n")
	//}
	//if len(data.Policy.Data.CreateOrderReq.OrderList)==0{
	//	return JsLogger.ErrorLog("checkInsurance failed,OrderList is empty\n")
	//}
	//for _,v:=range data.Policy.Data.CreateOrderReq.OrderList{
	//	if v.StartDate==""{
	//		return JsLogger.ErrorLog("checkInsurance failed,OrderList StartDate is empty\n")
	//	}
	//	if v.UWCount<=0{
	//		return JsLogger.ErrorLog("checkInsurance failed,OrderList UWCount <=0\n")
	//	}
	//	if len(v.CustomerList)==0{
	//		return JsLogger.ErrorLog("checkInsurance failed,CustomerList  len =0\n")
	//	}
	//	for _,u:=range v.CustomerList{
	//		if u.CustomerType==""||(u.CustomerType!="1"&&u.CustomerType!="2"){
	//			return JsLogger.ErrorLog("checkInsurance failed,CustomerList  CustomerType error\n")
	//		}
	//		if u.CustomerName==""{
	//			return JsLogger.ErrorLog("checkInsurance failed,CustomerList  CustomerName is empty\n")
	//		}
	//		if u.DocType==""{
	//			return JsLogger.ErrorLog("checkInsurance failed,CustomerList  DocType is empty\n")
	//		}
	//		if u.DocNo==""{
	//			return JsLogger.ErrorLog("checkInsurance failed,CustomerList  DocNo is empty\n")
	//		}
	//		if u.CustomerFlag==""{
	//			return JsLogger.ErrorLog("checkInsurance failed,CustomerList  CustomerFlag is empty\n")
	//		}
	//	}
	//
	//	if len(v.ItemAcciList)==0{
	//		return JsLogger.ErrorLog("checkInsurance failed,ItemAcciList  len =0\n")
	//	}
	//	for _,u:=range v.ItemAcciList{
	//		if u.Quantity==""{
	//			return JsLogger.ErrorLog("checkInsurance failed,ItemAcciList  Quantity is empty\n")
	//		}
	//		if u.NominativeInd==""{
	//			return JsLogger.ErrorLog("checkInsurance failed,ItemAcciList  NominativeInd is empty\n")
	//		}
	//		if u.OccupationCode==""{
	//			return JsLogger.ErrorLog("checkInsurance failed,ItemAcciList  OccupationCode is empty\n")
	//		}
	//		if len(u.AcciInsuredList)==0{
	//			return JsLogger.ErrorLog("checkInsurance failed,ItemAcciList  AcciInsuredList len is 0\n")
	//		}
	//	}
	//}

	return nil
}

//添加保单到用户
func appendUserInsurance(UID, ID string) error {
	return JsRedis.Redis_Sset(UID, ID)
}

//获取用户的保单id列表
func getUserInsurance(uid string) ([]string, error) {
	data := []string{}
	d, err := JsRedis.Redis_Sget(uid)
	for _, v := range d {
		data = append(data, string(v.([]byte)))
	}
	return data, err
}

func insuranceSMS(name, prod, phone, customer, number, hos, money string) {
	par := make(map[string]string)
	par["name1"] = name
	par["prod"] = prod
	par["phone"] = phone
	par["name2"] = customer
	par["number"] = number
	par["hosp"] = hos
	par["money"] = money
	JsMobile.ComJsMobileVerify("喜妹儿", "13961150133", "SMS_133150974", "c", 300, par)
	JsMobile.ComJsMobileVerify("喜妹儿", "18616379727", "SMS_133150974", "c", 300, par)
}
