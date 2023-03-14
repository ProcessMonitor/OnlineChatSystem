package main

import (
	"fmt"
	"time"
)

func ChannT() {
	ch1 := make(chan int)
	ch2 := make(chan int)
	go func() {
		time.Sleep(time.Second)
		ch1 <- 1
	}()

	go func() {
		ch2 <- 3
	}()

	for {
		select {
		case i := <-ch1:
			fmt.Printf("从ch1读取了数据%d", i)
		case j := <-ch2:
			fmt.Printf("从ch2读取了数据%d", j)
		}
		time.Sleep(time.Second)
	}
}
func buffChanT() {
	bufChan := make(chan int, 5)

	go func() {
		time.Sleep(time.Second)
		for {
			<-bufChan
			time.Sleep(5 * time.Second)
		}
	}()

	for {
		select {
		case bufChan <- 1:
			fmt.Println("add success")
			time.Sleep(time.Second)
		default:
			fmt.Println("资源已满，请稍后再试")
			time.Sleep(time.Second)
		}
	}
}
