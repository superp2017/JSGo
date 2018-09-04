package main

import (
	"JsGo/JsStore/JsRedis"
	"fmt"
	"math/rand"
	"time"
)

func shop() {
	var str string
	JsRedis.Redis_get("shop", &str)
	fmt.Printf("%s\n", str)
}

func admin() {
	type Admin struct {
		Account  string
		Password string
		Header   string
		Role     string
	}

	admin := &Admin{"admin", "1qaz2wsx", "http://stage.junsie.cn/static/admin.png", "admin"}

	JsRedis.Redis_hset("H_ADMIN", "admin", admin)
}

func main() {
	shop()
}

func main_e() {
	var v string = fmt.Sprintf("%d", time.Now().UnixNano())
	var k string = fmt.Sprintf("%d", time.Now().UnixNano())

	e := JsRedis.Redis_set(k, &v)
	if e != nil {
		fmt.Println(e.Error())
	} else {
		fmt.Printf("SET SUCCESS: key[%s] value[%s]\n", k, v)
	}

	b, e := JsRedis.Redis_exists(k)
	if e != nil {
		fmt.Println(e.Error())
	} else if b {
		fmt.Printf("EXIST SUCCESS: exist key = %s\n", k)
	} else {
		fmt.Printf("EXIST FAILED: exist key = %s\n", k)
	}

	e = JsRedis.Redis_get(k, &v)
	if e != nil {
		fmt.Println(e.Error())
	} else {
		fmt.Printf("GET SUCCESS: key[%s] value[%s]\n", k, v)
	}

	e = JsRedis.Redis_del(k)
	if e != nil {
		fmt.Println(e.Error())
	} else {

		b, e := JsRedis.Redis_exists(k)
		if e != nil {
			fmt.Println(e.Error())
		} else if b {
			fmt.Printf("DEL FAILED: key[%s] value[%s]\n", k, v)
		} else {
			fmt.Printf("DEL SUCCESS: key[%s] value[%s]\n", k, v)
		}

	}

	var t string = "TTable"
	v = fmt.Sprintf("%d", time.Now().UnixNano())
	k = fmt.Sprintf("%d", time.Now().UnixNano())
	e = JsRedis.Redis_hset(t, k, v)
	if e != nil {
		fmt.Println(e.Error())
	} else {
		fmt.Printf("HSET SUCCESS: table[%s] key[%s] value[%s]\n", t, k, v)
	}

	b, e = JsRedis.Redis_hexists(t, k)
	if e != nil {
		fmt.Println(e.Error())
	} else if b {
		fmt.Printf("HEXIST SUCCESS: table[%s] key[%s] value[%s]\n", t, k, v)
	} else {
		fmt.Printf("HEXIST FAILED: table[%s] key[%s] value[%s]\n", t, k, v)
	}

	e = JsRedis.Redis_hget(t, k, &v)
	if e != nil {
		fmt.Println(e.Error())
	} else {
		fmt.Printf("HGET SUCCESS: table[%s] key[%s] value[%s]\n", t, k, v)
	}

	e = JsRedis.Redis_hdel(t, k)
	if e != nil {
		fmt.Println(e.Error())
	} else {
		b, e = JsRedis.Redis_hexists(t, k)
		if e != nil {
			fmt.Println(e.Error())
		} else if b {
			fmt.Printf("HDEL FAILED: table[%s] key[%s] value[%s]\n", t, k, v)
		} else {
			fmt.Printf("HDEL SUCCESS: table[%s] key[%s] value[%s]\n", t, k, v)
		}
	}

	data := make(map[string]interface{})
	for i := 0; i < 2; i++ {
		k := fmt.Sprintf("%d", rand.Intn(100000000))
		data[k] = fmt.Sprintf("%d", rand.Intn(100000000))
		//fmt.Printf("data:%s,%v\n", k, data[k])
	}

	e = JsRedis.Redis_hmset(t, data)
	if e != nil {
		fmt.Println(e.Error())
	} else {
		fmt.Printf("HMSET SUCCESS: table[%s] insert %d ITEM\n", t, len(data))
	}

	n, e := JsRedis.Redis_hsize(t)
	if e != nil {
		fmt.Println(e.Error())
	} else if (int)(n) == len(data) {
		fmt.Printf("HSIZE SUCCESS: table[%s] size = %d\n", t, n)
	} else {
		fmt.Printf("HSIZE FAILED: table[%s] size = %d\n", t, n)
	}

	keys, e := JsRedis.Redis_hkeys(t)
	if e != nil {
		fmt.Println(e.Error())
	} else if len(keys) == len(data) {
		b := true
		for _, v := range keys {
			_, ok := data[v]
			if !ok {
				fmt.Printf("HKEYS FAILED: table[%s] keys[%s] not match\n", t, v)
				b = false
				break
			}
		}
		if b {
			fmt.Printf("HKEYS SUCCESS: table[%s]\n", t)
		}
	} else {
		fmt.Printf("HKEYS FAILED: table[%s] keys size[%d] not match data size[%d]\n", t, len(keys), len(data))
	}

	values := make(map[string]interface{})
	for k, _ := range data {
		var v string = ""
		values[k] = &v
	}
	e = JsRedis.Redis_hmget(t, &values)
	if e != nil {
		fmt.Println(e.Error())
	} else {
		b := true
		for k, v := range data {

			s, ok := values[k].(*string)

			if ok && v != *s {
				b = false
				fmt.Printf("HKEYS FAILED: table[%s] KV[%s:%s] not match data[%s:%s]\n", t, k, *s, k, v)
				break
			}
		}
		if b {
			fmt.Printf("HMGET SUCCESS: table[%s]\n", t)
		}
	}

	rand.Seed(time.Now().UnixNano())

	xdata := make(map[string]interface{})
	type MyData struct {
		A int
		B string
	}
	for i := 0; i < 20000; i++ {
		d := &MyData{A: 345, B: "36556"}
		d.A = rand.Intn(100000000)
		d.B = fmt.Sprintf("%d", rand.Intn(100000000))
		k := fmt.Sprintf("%d", rand.Intn(100000000))

		xdata[k] = d
		//fmt.Printf("data:%s,%v\n", k, data[k])
	}

	e = JsRedis.Redis_hmset(t, xdata)
	if e != nil {
		fmt.Println(e.Error())
	} else {
		fmt.Printf("HMSET SUCCESS: table[%s] insert %d ITEM\n", t, len(xdata))
	}

	xvalues := make(map[string]interface{})
	for k, _ := range xdata {
		v := &MyData{}
		xvalues[k] = v
	}
	e = JsRedis.Redis_hmget(t, &xvalues)
	if e != nil {
		fmt.Println(e.Error())
	} else {
		b := true
		for k, v := range xdata {

			s, ok := xvalues[k].(*MyData)
			u, uok := v.(*MyData)

			if ok && uok && u.A != s.A {
				b = false
				fmt.Printf("HKEYS FAILED: table[%s] KV[%s:%d] not match xdata[%s:%d]\n", t, k, s.A, k, u.A)
				break
			}
		}
		if b {
			fmt.Printf("HMGET SUCCESS: table[%s]", t)
		}
	}
}
