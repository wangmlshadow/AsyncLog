package main

import (
	"fmt"
	"sync"
	"src/LogModule"
	"time"
	"runtime"
)

var LOG = LogModule.LOG

var wg sync.WaitGroup

func threadFunc() {
	defer wg.Done()
	for i := 0; i < 100000; i++ {
		LOG(i)
	}
}

func type_test() {
	fmt.Println("type test...")
	LOG(0)
	LOG(1.0)
	LOG('c')
	LOG("abc")
}

func single_thread() {
	fmt.Println("single thread...")
	for i := 0; i < 100000; i++ {
		LOG(i)
	}
}

func multi_thread() {
	fmt.Println("multi thread...")
	for i := 0; i < 4; i++ {
		wg.Add(1)
		go threadFunc()
	}
	wg.Wait()
}

func main() {
	fmt.Println("main...")
	type_test()
	// single_thread()
	// multi_thread()
	for i := 0; i < 20; i++ {
		time.Sleep(time.Second)
		runtime.GC()
	}
}