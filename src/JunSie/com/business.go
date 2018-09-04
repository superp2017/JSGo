package com

import (
	"JsGo/JsBench/JsUser"
	"JsGo/JsConfig"
	"JsGo/JsHttp"
	"JsGo/JsLogger"
	"JsGo/JsMobile"
	"JsGo/JsStore/JsRedis"
	"JunSie/constant"
	"JunSie/util"
	"fmt"
)

func InitBusiness() {
	JsHttp.WhiteHttp("/applybusiness", ApplyBusiness)                              //申请商家入驻(申请)
	JsHttp.WhiteHttp("/queryapplyinfo", QueryApplyInfo)                            //查询商家入驻申请的信息
	JsHttp.WhiteHttp("/getglobalapplyinfo", GetGlobalApplyInfo)                    //获取全局的分页的申请列表(申请)
	JsHttp.WhiteHttp("/createbusiness", CreatBuiness)                              //创建创建商家基本信息(查看未启用的商家)
	JsHttp.WhiteHttp("/enablebusiness", EnableBusiness)                            //启用商户账号
	JsHttp.WhiteHttp("/getglobalbusiness", GetGlobalBusiness)                      //获取全局的分页的商家列表(分别看启用和未启用)
	JsHttp.WhiteHttp("/querybusinessinfo", QueryBusiness)                          //查询单商家信息(用于查看未启用的商家信息)
	JsHttp.WhiteHttp("/modifbusinessinfo", ModifyBusiness)                         //修改商家信息(只能修改未启用的商家信息)
	JsHttp.WhiteHttps("/modifubusinesspasswd", ModifyBusinessPasswd)               //修改商家的登录密码（启用才有密码修改）
	JsHttp.WhiteHttps("/modifybussinesswechatconfig", ModifyBussinessWeChatConfig) //修改商家微信配置（启用状态才可以用）
	JsHttp.WhiteHttp("/setexpressprice", SetExpressPrice)                          //设置运费
	JsHttp.WhiteHttp("/getexpressprice", GetExpressPrice)                          //获取运费
}

func InitBusinessMall() {
	JsHttp.WhiteHttps("/applybusiness", ApplyBusiness)   //申请商家入驻(申请)
	JsHttp.WhiteHttps("/queryapplyinfo", QueryApplyInfo) //查询商家入驻申请的信息
}

type Business struct { //商家基本信息
	BID             string                  //商家id
	Mobile          string                  //手机号
	BusinessName    string                  //店铺名称
	BusinessType    string                  //店铺类型
	BusinessAddress string                  //店铺地址  （必填）
	BosName         string                  //姓名 和地址中的名字重复
	Industry        string                  //行业
	BusinessLogo    string                  //商家标志（类似头像图片）
	BusinessDetail                          //补充完善商家信息
	RunTimeConfig   *JsConfig.RunTimeConfig //微信配置
	ConsignerAddrs  JsUser.RecAddr          //发货地址列表  consigner
	Status          string                  //状态	Status_ON    Status_OFF
	OpenTime        string                  //正式营业时间(启用账号时间)
	CreatTime       string                  //创建时间
}
type BusinessDetail struct { //补充完善商家信息（营业执照）business license
	CompanyName         string //公司名称   //company
	RegisterNumber      string //营业执照  注册账号注册号Register Number
	CompanyType         string //公司类型
	Address             string //住所 地址Address:
	LegalRepresentative string //法人代表Legal Representative
	TermBegin           string //经营期限（开始时间）Business Term
	TermEnd             string //经营期限（结束时间）Business Term
	ScopeOfBusiness     string //营业范围Scope of Business
	Date                string //签发日期Date:
	BLPicture           string //营业执照图片
}

type BussinessAccount struct {
	BID      string //商家id
	Account  string //账号
	Password string //密码
	Header   string //头像
	Role     string //权限
}
type ApplyBussiness struct {
	Name     string //姓名（检查字段）
	Mobile   string //手机号（检查字段）
	Industry string //行业
	LeaveMs  string //留言
}

//申请商家入驻
func ApplyBusiness(s *JsHttp.Session) {
	para := &ApplyBussiness{}
	err := s.GetPara(para)
	if err != nil {
		JsLogger.Error(err.Error())
		s.Forward("1", err.Error(), nil)
		return
	}
	if para.Name == "" || para.Mobile == "" {
		info := fmt.Sprintf("ApplyBuiness failed: Name=%s  Mobile=%s \n ", para.Name, para.Mobile)
		JsLogger.Error(info)
		s.Forward("2", info, nil)
		return
	}

	if exist, err := JsRedis.Redis_hexists(constant.H_ApplyBusiness, para.Mobile); err == nil && exist {
		info := fmt.Sprintf("ApplyBuiness failed:Name=%s  Mobile=%s was  already applied !!\n ", para.Name, para.Mobile)
		JsLogger.Error(info)
		s.Forward("3", info, nil)
		return
	}
	if err = JsRedis.Redis_hset(constant.H_ApplyBusiness, para.Mobile, para); err != nil {
		JsLogger.Error(err.Error())
		s.Forward("1", err.Error(), nil)
		return
	}
	s.Forward("0", "ApplyBusiness success\n", para)
}

func QueryApplyInfo(s *JsHttp.Session) {
	type Para struct {
		Mobile string //手机号
	}

	st := &Para{}
	if err := s.GetPara(st); err != nil {
		JsLogger.Error(err.Error())
		s.Forward("1", err.Error(), nil)
		return
	}
	if st.Mobile == "" {
		JsLogger.Error("GetApplyInfo failed,Mobile is empty\n")
		s.Forward("1", "GetApplyInfo failed,Mobile is empty\n", nil)
		return
	}
	d := &ApplyBussiness{}
	if err := JsRedis.Redis_hget(constant.H_ApplyBusiness, st.Mobile, d); err != nil {
		JsLogger.Error(err.Error())
		s.Forward("1", err.Error(), nil)
		return
	}
	s.Forward("1", "QueryApplyInfo success\n", d)
}

func GetGlobalApplyInfo(s *JsHttp.Session) {
	type Info struct {
		SIndex int //启始索引(数据库绝中的对位置，数组 从第几个开始)
		Size   int //个数（开始索引之后的 几个）
	}
	st := &Info{}
	if err := s.GetPara(st); err != nil {
		info := "GetGlobalApplyInfo,GetPara:" + err.Error()
		JsLogger.Error(info)
		s.Forward("1", info, nil)
		return
	}
	if st.SIndex < 0 || st.Size <= 0 {
		info := fmt.Sprintf("GetGlobalApplyInfo param error,SIndex=%d,Size=%d\n", st.SIndex, st.Size)
		JsLogger.Error(info)
		s.Forward("1", info, nil)
		return
	}

	ids, err := JsRedis.Redis_hkeys(constant.H_ApplyBusiness)
	if err != nil {
		s.Forward("1", err.Error(), nil)
		return
	}
	idslen := len(ids)
	data := getmoreApplyInfo(ids)
	if st.SIndex < idslen {
		if st.Size < idslen && (st.SIndex+st.Size) < idslen {
			s.Forward("0", "GetGlobalApplyInfo success", data[st.SIndex:st.SIndex+st.Size])
			return
		}
		s.Forward("0", "GetGlobalApplyInfo success", data[st.SIndex:])
		return
	}
	s.Forward("0", "GetGlobalApplyInfo success", nil)
}

//创建商家
func CreatBuiness(s *JsHttp.Session) {
	para := &Business{}
	if err := s.GetPara(para); err != nil {
		JsLogger.Error(err.Error())
		s.Forward("1", err.Error(), nil)
		return
	}
	if para.BosName == "" || para.Mobile == "" {
		info := fmt.Sprintf("CreatBuiness failed :Mobile=%s,Name=%s\n", para.Mobile, para.BosName)
		JsLogger.Error(info)
		s.Forward("2", info, nil)
		return
	}
	para.BID = util.IDer(constant.H_Business)
	para.CreatTime = util.CurTime()
	para.Status = constant.Status_OFF
	if err := JsRedis.Redis_hset(constant.H_Business, para.BID, para); err != nil {
		JsLogger.Error(err.Error())
		s.Forward("1", err.Error(), nil)
		return
	}
	s.Forward("0", "CreatBuiness success\n", para)
}

//启用商户账号
func EnableBusiness(s *JsHttp.Session) {
	type Para struct {
		BID string //商家id
	}
	st := &Para{}
	if err := s.GetPara(st); err != nil {
		JsLogger.Error(err.Error())
		s.Forward("1", err.Error(), nil)
		return
	}
	if st.BID == "" {
		JsLogger.Error("EnableBusiness FAILED,BID is empty\n")
		s.Forward("1", "EnableBusiness FAILED,BID is empty\n", nil)
		return
	}
	data := &Business{}
	if err := JsRedis.Redis_hget(constant.H_Business, st.BID, data); err != nil {
		JsLogger.Error(err.Error())
		s.Forward("1", err.Error(), nil)
		return
	}
	data.Status = constant.Status_ON
	data.OpenTime = util.CurTime()

	if err := JsRedis.Redis_hset(constant.H_Business, st.BID, data); err != nil {
		JsLogger.Error(err.Error())
		s.Forward("1", err.Error(), nil)
		return
	}
	go newBusinessAccount(data.BID, data.Mobile, data.BusinessLogo, "admin")
	s.Forward("0", "success", data)
}

//修改商家信息
func ModifyBusiness(s *JsHttp.Session) {
	st := &Business{}
	if err := s.GetPara(st); err != nil {
		JsLogger.Error(err.Error())
		s.Forward("1", err.Error(), nil)
		return
	}
	if st.BID == "" || st.BosName == "" {
		info := fmt.Sprintf("ModifyBusiness failed :Mobile=%s  BossName=%s\n", st.Mobile, st.BosName)
		JsLogger.Error(info)
		s.Forward("2", info, nil)
		return
	}
	data := &Business{}
	if err := JsRedis.Redis_hget(constant.H_Business, st.BID, data); err != nil {
		JsLogger.Error(err.Error())
		s.Forward("1", err.Error(), nil)
		return
	}
	JsLogger.Error("DATA=%v\n", data)
	///下面几个不准修改，使用旧的数据
	st.Mobile = data.Mobile
	st.Status = data.Status
	st.OpenTime = data.OpenTime
	st.CreatTime = data.CreatTime
	if err := JsRedis.Redis_hset(constant.H_Business, st.BID, st); err != nil {
		JsLogger.Error(err.Error())
		s.Forward("1", err.Error(), nil)
		return
	}
	JsLogger.Error("st=%v\n", st)
	s.Forward("0", "ModifyBusiness success\n", st)
}

//修改商家账号的密码
func ModifyBusinessPasswd(s *JsHttp.Session) {
	type Para struct {
		Mobile   string //手机号
		Password string //密码
		SmsCode  string //短信验证码
	}
	st := &Para{}
	if err := s.GetPara(st); err != nil {
		JsLogger.Error(err.Error())
		s.Forward("1", err.Error(), nil)
		return
	}
	if st.Password == "" || st.Mobile == "" || st.SmsCode == "" {
		info := fmt.Sprintf("ModifuBusinessPasswd failed :Cell=%s  Password=%s,SmsCode=%s\n", st.Mobile, st.Password, st.SmsCode)
		JsLogger.Error(info)
		s.Forward("2", info, nil)
		return
	}

	if ok := JsMobile.VerifySmsCode(st.Mobile, st.SmsCode); !ok {
		info := fmt.Sprintf("ModifuBusinessPasswd VerifySmsCode failed,Cell = %s,SmsCode =%s \n", st.Mobile, st.SmsCode)
		JsLogger.Error(info)
		s.Forward("1", info, nil)
		return
	}

	//获取商家账户信息
	data, err := GetAccountInfo(st.Mobile)
	if err != nil {
		JsLogger.Error(err.Error())
		s.Forward("1", err.Error(), nil)
		return
	}
	data.Password = st.Password
	if err := JsRedis.Redis_hset(constant.ADMIN, st.Mobile, data); err != nil {
		JsLogger.Error(err.Error())
		s.Forward("1", err.Error(), nil)
		return
	}
	s.Forward("0", "ModifuBusinessPasswd success\n", data)
}

func ModifyBussinessWeChatConfig(s *JsHttp.Session) {
	type Para struct {
		BID           string                 //商户id
		RunTimeConfig JsConfig.RunTimeConfig //微信配置
	}
	st := &Para{}
	if err := s.GetPara(st); err != nil {
		JsLogger.Error(err.Error())
		s.Forward("1", err.Error(), nil)
		return
	}
	if st.BID == "" {
		JsLogger.Error("ModifyBussinessWeChatConfig failed,BID is empty\n")
		s.Forward("1", "ModifyBussinessWeChatConfig failed,BID is empty\n", nil)
		return
	}
	data := Business{}
	if err := JsRedis.Redis_hget(constant.H_Business, st.BID, &data); err != nil {
		JsLogger.Error(err.Error())
		s.Forward("1", err.Error(), nil)
		return
	}
	data.RunTimeConfig = &st.RunTimeConfig
	if err := JsRedis.Redis_hset(constant.H_Business, st.BID, &data); err != nil {
		JsLogger.Error(err.Error())
		s.Forward("1", err.Error(), nil)
		return
	}
	s.Forward("0", "sucess\n", data)
}

//获取全局的分页的商家列表
func GetGlobalBusiness(s *JsHttp.Session) {
	type Info struct {
		Status string //状态
		SIndex int    //启始索引
		Size   int    //个数
	}
	st := &Info{}
	if err := s.GetPara(st); err != nil {
		info := "GetGlobalBusiness,GetPara:" + err.Error()
		JsLogger.Error(info)
		s.Forward("1", info, nil)
		return
	}
	if st.Status == "" || (st.Status != constant.Status_ON && st.Status != constant.Status_OFF) || st.SIndex < 0 || st.Size <= 0 {
		info := fmt.Sprintf("GetGlobalBusiness param error,Status=%s,SIndex=%d,Size=%d\n", st.Status, st.SIndex, st.Size)
		JsLogger.Error(info)
		s.Forward("1", info, nil)
		return
	}
	data := []*Business{}

	ids, err := JsRedis.Redis_hkeys(constant.H_Business)
	if err != nil {
		s.Forward("1", err.Error(), nil)
		return
	}
	list := getmoreBusiness(ids)
	for _, v := range list {
		if v.Status == st.Status {
			data = append(data, v)
		}
	}
	if st.SIndex < len(data) {
		if st.Size < len(data) && (st.SIndex+st.Size) < len(data) {
			s.Forward("0", "GetGlobalBusiness success", data[st.SIndex:st.SIndex+st.Size])
			return
		}
		s.Forward("0", "GetGlobalBusiness success", data[st.SIndex:])
		return
	}
	s.Forward("0", "GetGlobalBusiness success", nil)
}

//查询单个商家信息
func QueryBusiness(s *JsHttp.Session) {
	type Para struct {
		BID string
	}
	para := &Para{}
	if err := s.GetPara(para); err != nil {
		JsLogger.Error(err.Error())
		s.Forward("1", err.Error(), nil)
		return
	}
	if para.BID == "" {
		JsLogger.Error("QueryBusiness FAILED,BID is empty\n")
		s.Forward("1", "QueryBusiness FAILED,BID is empty\n", nil)
		return
	}
	data := Business{}
	if err := JsRedis.Redis_hget(constant.H_Business, para.BID, &data); err != nil {
		JsLogger.Error(err.Error())
		s.Forward("1", err.Error(), nil)
		return
	}
	s.Forward("0", "success", data)

}

//获取多个商家信息
func getmoreBusiness(ids []string) []*Business {
	data := make([]*Business, 0, len(ids))
	for _, v := range ids {
		d := &Business{}
		if err := JsRedis.Redis_hget(constant.H_Business, v, d); err == nil {
			data = append(data, d)
		}
	}
	return data
}

//获取多个申请信息
func getmoreApplyInfo(ids []string) []*ApplyBussiness {
	data := make([]*ApplyBussiness, 0, len(ids))
	for _, v := range ids {
		d := &ApplyBussiness{}
		if err := JsRedis.Redis_hget(constant.H_ApplyBusiness, v, d); err != nil {
			JsLogger.Error(err.Error())
			continue
		} else {
			data = append(data, d)
		}
	}
	return data
}

//创建商户账号
func newBusinessAccount(BID, Account, Head, Role string) error {
	st := BussinessAccount{
		BID:      BID,
		Account:  Account,
		Header:   Head,
		Role:     Role,
		Password: "123456",
	}
	return JsRedis.Redis_hset(constant.ADMIN, st.Account, &st)
}

//获取商家账户信息
func GetAccountInfo(Mobile string) (*BussinessAccount, error) {
	data := &BussinessAccount{}
	err := JsRedis.Redis_hget(constant.ADMIN, Mobile, data)
	return data, err
}

//设置邮费
func SetExpressPrice(s *JsHttp.Session) {
	type Para struct {
		Price map[string]int
	}
	st := &Para{}
	if err := s.GetPara(st); err != nil {
		JsLogger.Error(err.Error())
		s.Forward("1", err.Error(), nil)
		return
	}
	if st.Price == nil {
		str := fmt.Sprintf("SetExpressPrice failedPrice=%v", st.Price)
		JsLogger.Error(str)
		s.Forward("1", str, nil)
		return
	}
	if err := JsRedis.Redis_set("ExpressPrice", &st.Price); err != nil {
		JsLogger.Error(err.Error())
		s.Forward("1", err.Error(), nil)
		return
	}
	s.Forward("0", "suceess", nil)
}

//获取邮费
func GetExpressPrice(s *JsHttp.Session) {
	Price := make(map[string]int)
	if err := JsRedis.Redis_get("ExpressPrice", &Price); err != nil {
		JsLogger.Error(err.Error())
		s.Forward("1", err.Error(), nil)
		return
	}
	s.Forward("0", "success", Price)
}

//查询物流价格
func queryExpressPrice(Provence string) int {
	Price := make(map[string]int)
	//字段截取，取地址中的前两个文字
	nameRune := []rune(Provence)
	provence := string(nameRune[:2])
	if err := JsRedis.Redis_get("ExpressPrice", &Price); err != nil {
		JsLogger.Error(err.Error())
		return 10000
		//找不到默认运费100元
	}
	if v, ok := Price[provence]; ok {
		return v
	}
	return 10000
}
