package com

import (
	"JsGo/JsHttp"
	"JsGo/JsLogger"
	"JsGo/JsStore/JsRedis"
	"JunSie/constant"
)

func GetOderList(s *JsHttp.Session) {
	type para struct {
		Status string //状态 //订单发货 "Send"//已发货列表 "Receive"//订单评价完成 "Success"
		SIndex int    //从第几个开始
		Size   int    //个数
	}
	st := &para{}
	if e := s.GetPara(st); e != nil {
		JsLogger.Error(e.Error())
		s.Forward("1", e.Error(), nil)
		return
	}
	if st.Size < 1 || st.SIndex < 0 {
		info := "Size<1||SIndex<0:"
		JsLogger.Error(info, st.SIndex, st.Size)
		s.Forward("2", info, nil)
		return
	}
	type parat struct {
		Num int
		IDs []string
	}
	data := parat{}

	switch st.Status {
	case "OrderStatus_Paid":
		JsRedis.Redis_hget(constant.H_Order, constant.OrderStatus_Paid, &data)
	case "OrderStatus_Send":
		JsRedis.Redis_hget(constant.H_Order, constant.OrderStatus_Send, &data)
	case "OrderStatus_Receive":
		JsRedis.Redis_hget(constant.H_Order, constant.OrderStatus_Receive, &data)
	default:
		info := "Status error"
		JsLogger.Error(info)
		s.Forward("2", info, nil)
		return
	}
	idlist := []string{}
	lenList := len(data.IDs)
	if st.SIndex > lenList {
		s.Forward("0", "success", nil)
		return
	} else {
		if (st.SIndex + st.Size) < lenList {
			idlist = data.IDs[st.SIndex:st.Size]
		} else {
			idlist = data.IDs[st.SIndex:lenList]
		}
	}
	type paras struct {
		Num    int
		Orders []*OrderABS
	}
	dataf := &paras{}
	for _, v := range idlist {
		datas := &OrderABS{}
		err := JsRedis.Redis_hget(constant.H_Order, v, datas)
		if err == nil {
			dataf.Orders = append(dataf.Orders, datas)
		}
	}
	dataf.Num = lenList
	s.Forward("0", "succes", dataf)
}

//支付完成调用
//添加到到————待发货列表（正序，先下单的放在前面）//订单发货列表 "Send"
func AddOrderSend(id string) {
	type para struct {
		Num int
		IDs []string
	}
	data := para{}
	if err := JsRedis.Redis_hget(constant.H_Order, constant.OrderStatus_Paid, &data); err != nil {
		JsRedis.Redis_hget(constant.H_Order, constant.OrderStatus_Paid, &data)
	}
	data.IDs = append(data.IDs, id)
	data.Num = len(data.IDs)
	JsRedis.Redis_hset(constant.H_Order, constant.OrderStatus_Paid, &data)
}

//发货完成调用
//从待发货列表中转到————已发货列表（倒序，刚发货的放在最前面）//已发货列表 "Receive"
func AddOrderReceive(id string) {
	type para struct {
		Num int
		IDs []string
	}
	data := para{}
	//先填加到下一个表中，然后在删除
	if err := JsRedis.Redis_hget(constant.H_Order, constant.OrderStatus_Send, &data); err != nil {
		JsRedis.Redis_hget(constant.H_Order, constant.OrderStatus_Send, &data)
	}
	temp := data.IDs
	data.IDs = make([]string, 0, len(temp)+1)
	data.IDs = append(data.IDs, id)
	data.IDs = append(data.IDs, temp...)
	data.Num = len(data.IDs)
	if err := JsRedis.Redis_hset(constant.H_Order, constant.OrderStatus_Send, &data); err != nil {
		JsLogger.Error(err.Error())
		return
	}
	if err := JsRedis.Redis_hget(constant.H_Order, constant.OrderStatus_Paid, &data); err != nil {
		JsRedis.Redis_hget(constant.H_Order, constant.OrderStatus_Paid, &data)
	}
	for i, v := range data.IDs {
		if v == id {
			data.IDs = append(data.IDs[:i], data.IDs[i+1:]...)
			data.Num = len(data.IDs)
			JsRedis.Redis_hset(constant.H_Order, constant.OrderStatus_Paid, &data)
			break
		}
	}

}

//用户确认收货
//已发货列表转移到————已完成列表（倒序，最新的放在最前面）//订单评价完成 "Success"
func AddOrderSuccess(id string) {
	type para struct {
		Num int
		IDs []string
	}
	data := para{}
	if err := JsRedis.Redis_hget(constant.H_Order, constant.OrderStatus_Receive, &data); err != nil {
		JsRedis.Redis_hget(constant.H_Order, constant.OrderStatus_Receive, &data)
	}
	temp := data.IDs
	data.IDs = make([]string, 0, len(temp)+1)
	data.IDs = append(data.IDs, id)
	data.IDs = append(data.IDs, temp...)
	data.Num = len(data.IDs)
	if err := JsRedis.Redis_hset(constant.H_Order, constant.OrderStatus_Receive, &data); err != nil {
		JsLogger.Error(err.Error())
		return
	}
	if err := JsRedis.Redis_hget(constant.H_Order, constant.OrderStatus_Send, &data); err != nil {
		JsRedis.Redis_hget(constant.H_Order, constant.OrderStatus_Send, &data)
	}
	for i, v := range data.IDs {
		if v == id {
			data.IDs = append(data.IDs[:i], data.IDs[i+1:]...)
			data.Num = len(data.IDs)
			JsRedis.Redis_hset(constant.H_Order, constant.OrderStatus_Send, &data)
			break
		}
	}
}
