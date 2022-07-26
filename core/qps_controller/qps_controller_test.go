package qps_controller

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

func TestGroup_Do(t *testing.T) {
	mug := sync.WaitGroup{}
	g := NewGroup(time.Second, 10)
	fn := func() (interface{}, error) {
		defer mug.Done()
		fmt.Println(1) //如果使用数字和或者队列测试需要加锁，会导致协程等待锁使得下面的每秒输出的值被阻塞，使输出的值不准确,直接输出也会有点顺序问题，但是因为不用等待锁，大体上不会差多少
		return nil, nil
	}
	tk := time.NewTicker(time.Second)
	go func() {
		for {
			<-tk.C
			fmt.Println(0)
		}
	}()
	for i := 0; i < 10; i++ {
		mug.Add(1)
		go g.Do(fn)
	}
	mug.Wait()
	time.Sleep(time.Second)
}

func TestGroup_DoChan(t *testing.T) {
	mug := sync.WaitGroup{}
	g := NewGroup(time.Second, 10)
	fn := func() (interface{}, error) {
		defer mug.Done()
		fmt.Println(1)
		return 1, nil
	}
	tk := time.NewTicker(time.Second)
	go func() {
		for {
			<-tk.C
			fmt.Println(0)
		}
	}()
	results := []<-chan Result{}
	for i := 0; i < 50; i++ {
		mug.Add(1)
		results = append(results, g.DoChan(fn))
	}
	mug.Wait()
	time.Sleep(time.Second)
	for i := 0; i < len(results); i++ {
		fmt.Println(i, ":	", <-results[i])
	}
}
