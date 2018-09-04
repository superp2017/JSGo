package com

import (
	"JsGo/JsStore/JsRedis"
	"JunSie/constant"
	"sync"
	"time"
	"strings"
)

var WaitPay map[string]*Order = make(map[string]*Order) //待支付的订单
var IsTimerRun bool = false                             //是否运行定时
var mutex_pay sync.Mutex
var mutex_timer sync.Mutex

func init() {
	list, err := JsRedis.Redis_hkeys(constant.H_Order)
	if err == nil {
		for _, v := range list {
			if strings.Count(v, "")-1 < 30 {
				continue
			}
			d := &Order{}
			if err := JsRedis.Redis_hget(constant.H_Order, v, d); err == nil {
				WaitPay[d.OrderID] = d
			}
		}
	}

	if len(WaitPay) > 0 {
		IsTimerRun = true
		go StartOrderTimer()
	}
}

///添加一个订单到待支付的队列
func appendWaitPayOrder(order *Order) {
	mutex_pay.Lock()
	defer mutex_pay.Unlock()
	WaitPay[order.OrderID] = order
	if len(WaitPay) > 0 {
		if !IsTimerRun {
			mutex_timer.Lock()
			defer mutex_timer.Unlock()
			IsTimerRun = true
			go StartOrderTimer()
		}
	}
}

//从待支付的队列中删除
func removeWaitPayOrder(orderID string) {
	mutex_pay.Lock()
	defer mutex_pay.Unlock()
	if _, ok := WaitPay[orderID]; ok {
		delete(WaitPay, orderID)
	}
}

////启动定时任务
func StartOrderTimer() {
	go timeTask(checkOrderTimeout, 1*60)
}

////定时任务
func timeTask(task func(), minus int64) {
	for {
		if !IsTimerRun {
			break
		}
		//////////定时任务//////////////////
		t := time.NewTimer(time.Duration(minus))
		<-t.C
		task()
	}
}

//检查无效的订单
func checkOrderTimeout() {
	delList := make(map[string][]*Order)
	for _, v := range WaitPay {
		if v.Current.Status == constant.OrderStatus_Creat && time.Now().Unix()-v.Current.TimeStamp > 3600 {
			delList[v.UID] = append(delList[v.UID], v)
			go updateInventory(v, false)
		}
	}

	if len(delList) > 0 {
		for k, v := range delList {
			go removeInvalidOrder(k, v)
		}
		mutex_pay.Lock()
		defer mutex_pay.Unlock()
		for _, v := range delList {
			for _, v1 := range v {
				delete(WaitPay, v1.OrderID)
			}
		}
	}

	mutex_pay.Lock()
	defer mutex_pay.Unlock()
	if len(WaitPay) > 0 {
		if !IsTimerRun {
			IsTimerRun = true
			go StartOrderTimer()
		}
	} else {
		IsTimerRun = false
	}
}

//删除无效的订单
func removeInvalidOrder(UID string, Order []*Order) {
	list := []string{}
	JsRedis.Redis_hget(constant.H_UserOrder, UID, &list)
	for _, v := range Order {
		index := -1
		for i, v1 := range list {
			if v1 == v.OrderID {
				index = i
				break
			}
		}
		if index != -1 {
			list = append(list[:index], list[index+1:]...)
		}
		go JsRedis.Redis_hdel(constant.H_Order, v.OrderID)
	}
	go JsRedis.Redis_hset(constant.H_UserOrder, UID, &list)
}
