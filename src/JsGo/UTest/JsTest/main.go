package main

import (
	"fmt"
	"time"
)

type I interface {
}

type A struct {
	a int
}

type B struct {
	A
	b string
}

func print(i I) {
	b, ok := i.(*A)
	if ok {
		fmt.Printf("v = %d\n", b.a)
		b.a = 10
	} else {
		fmt.Printf("is not A\n")
	}

}

func main() {
	e := make(chan int)
	c := make(chan int, 10)
	go get(c)
	go put(c, e)

	<-e
	fmt.Printf("finished \n")
}

func get(c chan int) {
	for i := 0; i < 10; i++ {
		v := <-c
		fmt.Printf("v = %d\n", v)
	}
}

func put(c chan int, e chan int) {
	for i := 0; i < 10; i++ {

		c <- i + 53
		fmt.Printf("put v = %d\n", i+53)
		c <- i + 85
		fmt.Printf("put v = %d\n", i+53)
		c <- i + 96

		time.Sleep(time.Second * 2)
	}
	e <- 1
}
