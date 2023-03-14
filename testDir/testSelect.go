package main

import (
	"log"
	"sync"
	"time"
)

func write(ch chan int, num int) {
	for i := 0; i < 10001; i++ {
		ch <- num
	}
}
func mutiChannSelectT() {
	countMap := make(map[int]int, 0)
	mu := sync.Mutex{}
	ch1 := make(chan int, 100000)
	ch2 := make(chan int, 100000)
	ch3 := make(chan int, 100000)
	ch4 := make(chan int, 100000)
	wg := &sync.WaitGroup{}
	wg.Add(4)
	go func() {
		defer wg.Done()
		write(ch1, 1)
	}()
	go func() {
		defer wg.Done()
		write(ch2, 2)
	}()
	go func() {
		defer wg.Done()
		write(ch3, 3)
	}()
	go func() {
		defer wg.Done()
		write(ch4, 4)
	}()
	wg.Wait()

	time.Sleep(time.Second * 2)
	for i := 0; i < 10000; i++ {
		select {
		case data, ok := <-ch1:
			if ok {
				mu.Lock()
				countMap[data]++
				mu.Unlock()
			} else {
				log.Println("ch1-channel-closed!")
			}
		case data, ok := <-ch2:
			if ok {
				mu.Lock()
				countMap[data]++
				mu.Unlock()
			} else {
				log.Println("ch2-channel-closed!")
			}
		case data, ok := <-ch3:
			if ok {
				mu.Lock()
				countMap[data]++
				mu.Unlock()
			} else {
				log.Println("ch3-channel-closed!")
			}
		case data, ok := <-ch4:
			if ok {
				mu.Lock()
				countMap[data]++
				mu.Unlock()
			} else {
				log.Println("ch4-channel-closed!")
			}
		}
	}
	log.Println(countMap)
}
