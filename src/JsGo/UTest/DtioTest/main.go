package main

import (
	"JsGo/JsBench/JsUser"
	"JsGo/JsStore"
	"fmt"
	"math/rand"
	"time"
)

var n []int
var c chan int

const (
	TOTAL int64 = 5000
	NUM   int   = 3000
)

var uids []string

func Insert(x *int) {
	user := &JsUser.User{}
	user.Nickname = fmt.Sprintf("Nickname = %d", time.Now().UnixNano())
	user.Header = fmt.Sprintf("Header = %d", time.Now().UnixNano())
	user.Birthday = "1958-06-03"
	user.Mobile = "15951862395"
	user.Sex = 1
	user.Role = "Admin"
	user.Address = fmt.Sprintf("Addrsss = %d", time.Now().UnixNano())
	user.City = "Shanghai"
	user.Country = "China"
	user.Profession = "Software"
	user.Tag = "normal"
	user.Email = "mengzhaofeng@163.com"
	user.Account = fmt.Sprintf("%d", time.Now().UnixNano())
	user.Password = "sdflkjowiutrojsw"
	user.Openid = "sdfsdfsdfsdfjlksdjfsadfkas"
	user.Unionid = "iowujqpoiuw nvoiqwerpovwer"
	for i := 0; i < 10000; i++ {

		user.NewUser()
		(*x)++
	}

	c <- 1
}

func Count(len int) {

	for {
		x := 0
		for i := 0; i < len; i++ {
			fmt.Printf("%d ", n[i])
			x += n[i]
			n[i] = 0
		}
		if x == 0 {
			fmt.Printf("x == 0\n")
			data := JsStore.CurLock()
			for _, v := range data {
				fmt.Printf("Table[%s] ID[%s] locked, Waiting:%d, lock time %s\n", v.Table, v.Id, v.C, v.T)
			}
		}
		fmt.Println()
		time.Sleep(time.Second)
	}
}

func Query(x *int) {
	for i := 0; i < NUM; i++ {
		v := rand.Intn((int)(TOTAL))
		//v := i
		user := &JsUser.User{}
		e := JsStore.ShareLock("User", uids[v], user)
		if e != nil {
			fmt.Println(e.Error())
		}
		(*x)++
	}
	fmt.Printf("Query:%d, finished\n", *x)

	c <- 1
}

func Update(x *int) {
	for i := 0; i < NUM; i++ {
		v := rand.Intn((int)(TOTAL))
		//v := i
		user := &JsUser.User{}
		e := JsStore.WriteLock("User", uids[v], user)
		if e != nil {
			fmt.Println(e.Error())
		}

		user.Nickname = fmt.Sprintf("NickName = %d", time.Now().UnixNano())
		JsStore.WriteBack("User", user.UID, user)
		(*x)++
	}

	fmt.Printf("Update:%d, finished\n", *x)

	c <- 1
}

func Init() {
	uids = make([]string, 1)

	ret, e := JsStore.HKeys("User", "", "", TOTAL)
	if e == nil {
		uids = ret
	} else {
		fmt.Println(e.Error())
	}
}

func main() {
	// c = make(chan int, 1)

	// len := 10
	// n = make([]int, len)
	// for i := 0; i < len; i++ {
	// 	n[i] = 0
	// 	go Insert(&n[i])
	// }

	// go Count(len)

	// <-c

	// fmt.Println("finished")

	fmt.Printf("begin: %s\n", time.Now().String())
	Init()
	fmt.Printf("Init finised: %s\n", time.Now().String())
	len := 20
	c = make(chan int, len)

	n = make([]int, len)
	go Count(len)
	for i := 0; i < len/2; i++ {
		n[i] = 0
		go Query(&n[i])
	}

	for i := len / 2; i < len; i++ {
		n[i] = 0
		go Update(&n[i])
	}

	for i := 0; i < len; i++ {
		<-c
	}

	fmt.Printf("Query finised: %s\n", time.Now().String())

	data := JsStore.CurLock()
	for _, v := range data {
		fmt.Printf("Table[%s] ID[%s] locked, Waiting:%d, lock time %s\n", v.Table, v.Id, v.C, v.T)
	}

}

func hkeys() {

	ret, e := JsStore.HKeys("User", "", "", 10000)
	if e == nil {
		for _, v := range ret {
			fmt.Printf("%s\n", v)
		}
	} else {
		fmt.Println(e.Error())
	}
}

func newIm() {
	b, e := JsStore.HExist("User", "ae2724f6-8c10-4c6b-a81f-dd8fd21b4c02")
	if e != nil {
		fmt.Println(e.Error())
	} else if b {
		fmt.Println("exist")
	} else {
		fmt.Println("not exist")
	}

	s, e := JsStore.HSize("User")
	if e != nil {
		fmt.Println(e.Error())
	} else if s > 0 {
		fmt.Printf("s = %d\n", s)
	} else {
		fmt.Println("s = -1")
	}

	b, e = JsStore.Exist("Users")
	if e != nil {
		fmt.Println(e.Error())
	} else if b {
		fmt.Println("exist")
	} else {
		fmt.Println("not exist")
	}
}
