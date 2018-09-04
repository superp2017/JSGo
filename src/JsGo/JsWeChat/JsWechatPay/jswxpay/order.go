package jswxpay

type ST_PayCb struct {
	AppId                string `xml:"appid"`
	Mch_id               string `xml:"mch_id"`
	Device_info          string `xml:"device_info"`
	Nonce_str            string `xml:"nonce_str"`
	Sign                 string `xml:"sign"`
	Sign_type            string `xml:"sign_type"`
	Result_code          string `xml:"result_code"`
	Err_code             string `xml:"err_code"`
	Err_code_des         string `xml:"err_code_des"`
	Openid               string `xml:"openid"`
	Is_subscribe         string `xml:"is_subscribe"`
	Trade_type           string `xml:"trade_type"`
	Bank_type            string `xml:"bank_type"`
	Total_fee            string `xml:"total_fee"`
	Settlement_total_fee string `xml:"settlement_total_fee"`
	Fee_type             string `xml:"fee_type"`
	Cash_fee             string `xml:"cash_fee"`
	Cash_fee_type        string `xml:"cash_fee_type"`
	Transaction_id       string `xml:"transaction_id"`
	Out_trade_no         string `xml:"out_trade_no"`
	Attach               string `xml:"attach"`
	Time_end             string `xml:"time_end"`
}

type StRefundCb struct {
	Result_code           string `xml:result_code`
	Err_code              string `xml:err_code`
	err_code_des          string `xml:err_code_des`
	Appid                 string `xml:appid`
	Mch_id                string `xml:mch_id`
	Device_info           string `xml:device_info`
	Nonce_str             string `xml:nonce_str`
	Sign                  string `xml:sign`
	Transaction_id        string `xml:transaction_id`
	Out_trade_no          string `xml:out_trade_no`
	Out_refund_no         string `xml:out_refund_no`
	Refund_id             string `xml:refund_id`
	Refund_channel        string `xml:refund_channel`
	Refund_fee            string `xml:refund_fee`
	Settlement_refund_fee string `xml:settlement_refund_fee`
	Total_fee             string `xml:total_fee`
	Settlement_total_fee  string `xml:settlement_total_fee`
	Fee_type              string `xml:fee_type`
	Cash_fee              string `xml:cash_fee`
	Cash_refund_fee       string `xml:cash_refund_fee`
}

type StRefund struct {
	Result_code           string `xml:result_code`
	Err_code              string `xml:err_code`
	Err_code_des          string `xml:err_code_des`
	Appid                 string `xml:appid`
	Mch_id                string `xml:mch_id`
	Device_info           string `xml:device_info`
	Nonce_str             string `xml:nonce_str`
	Sign                  string `xml:sign`
	Transaction_id        string `xml:transaction_id`
	Out_trade_no          string `xml:out_trade_no`
	Out_refund_no         string `xml:out_refund_no`
	Refund_id             string `xml:refund_id`
	Refund_channel        string `xml:refund_channel`
	Refund_fee            string `xml:refund_fee`
	Settlement_refund_fee string `xml:settlement_refund_fee`
	Total_fee             string `xml:total_fee`
	Settlement_total_fee  string `xml:settlement_total_fee`
	Fee_type              string `xml:fee_type`
	Cash_fee              string `xml:cash_fee`
	Cash_refund_fee       string `xml:cash_refund_fee`
}

type StOrder struct {
	OrderId       string
	RefundId      string //退款ID
	TransferId    string //转移转让  ？??
	AppId         string
	Mch_id        string //？？
	Nonce_str     string //？？？？
	ProjectId     string
	GoodsId       []string
	SellerId      string //卖方ID
	Status        string
	Type          string
	Amount        int //总计
	UserName      string
	UserHeader    string
	Desc          string //降序排列？？？
	Subject       string
	Channel       string //通道？？？
	OpenId        string
	Uid           string
	LDate         string //日期???
	TimeStamp     int64  //创建订单时间戳
	PaidTimeStamp int64  //支付成功的时间戳
	TerminalIp    string //终端IP
	TransactionId string //交易ID
	Time_end      string
	Trade_type    string //交易
	Bank_type     string
	PayCb         *ST_PayCb
	Charge        map[string]string //票据
	Refund        *StRefund
	RefundCb      map[string]string
	RefundFee     int
	CreateDate    string            //创建日期
	ExData        map[string]string //扩展的数据
}

type StOrderTable struct {
	ProjectId        string
	ComOrder         string
	ComTransfer      string
	UnTransfer       string
	IncomeOrder      string //订单收入
	ExpenditureOrder string
}

type StBill struct {
	Id       string
	TotalFee int
	Year     string
	Mon      string
	Ids      []string
}

type ST_Transfer struct {
	Tid              string
	OpenId           string
	OrderId          string
	UserName         string
	UserHeader       string
	LDate            string
	TimeStamp        int64
	Theme            string
	ProjectId        string
	Desc             string
	MchId            string
	AppId            string
	CheckName        string
	ReUserName       string
	Amount           int
	Spbill_create_ip string
	TransferCb       map[string]string
	LastError        string
	LastCb           map[string]string
}
