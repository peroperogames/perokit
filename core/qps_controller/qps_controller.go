package qps_controller

import (
	"bytes"
	"fmt"
	"math/rand"
	"runtime/debug"
	"strconv"
	"sync"
	"time"
)

//定义panic时抛出错误
type panicError struct {
	value interface{}
	stack []byte
}

//实现error接口
func (p *panicError) Error() string {
	return fmt.Sprintf("%v\n\n%s", p.value, p.stack)
}

func newPanicError(v interface{}) error {
	stack := debug.Stack()

	// The first line of the stack trace is of the form "goroutine N [status]:"
	// but by the time the panic reaches Do the goroutine may no longer exist
	// and its status will have changed. Trim out the misleading line.
	if line := bytes.IndexByte(stack[:], '\n'); line >= 0 {
		stack = stack[line+1:]
	}
	return &panicError{value: v, stack: stack}
}

type Result struct {
	Val interface{}
	Err error
}

type keyFn struct {
	key string
	fn  func() (interface{}, error)
}

type Group struct {
	mu  sync.Mutex
	c   chan struct{}
	m   map[string]chan Result
	fns []keyFn
}

func (g *Group) Do(fn func() (interface{}, error)) (v interface{}, err error) {
	<-g.c
	v, err = fn()
	return
}

//异步调用时，返回一个通道等待结果
func (g *Group) DoChan(fn func() (interface{}, error)) <-chan Result {
	ch := make(chan Result, 1)
	t := time.Now().String() + strconv.Itoa(rand.Int()) //使用时间+随机数字当key
	g.mu.Lock()
	g.m[t] = ch
	g.mu.Unlock()
	g.fns = append(g.fns, keyFn{ //把方法加入队列
		key: t,
		fn:  fn,
	})
	return ch
}

func (g *Group) doCall(key string, fn func() (interface{}, error)) {
	g.mu.Lock()
	ch := g.m[key] //取得返回channel
	g.mu.Unlock()
	defer func() {
		g.mu.Lock() //锁住，删除这个channel
		delete(g.m, key)
		g.mu.Unlock()
	}()

	defer func() { //捕获panic
		if r := recover(); r != nil {
			err := newPanicError(r)
			ch <- Result{
				Val: nil,
				Err: err,
			}
			return
		}
	}()

	val, err := fn() //执行返回
	ch <- Result{
		Val: val,
		Err: err,
	}
	return
}

//获取QPS剩余量
func (g *Group) GetNum() (n int) {
	g.mu.Lock()
	n = len(g.c)
	g.mu.Unlock()
	return
}

//设置qps更新时间和QPS
func NewGroup(t time.Duration, QPS int) *Group {
	g := &Group{
		mu:  sync.Mutex{},
		c:   make(chan struct{}, QPS),
		m:   map[string]chan Result{},
		fns: []keyFn{},
	}
	go func() { //起一个携程，根据设置的QPS时间刷新QPS量
		tick := time.NewTicker(t)
		for {
			<-tick.C
			g.mu.Lock()
			l := QPS - len(g.c)
			for i := 0; i < l; i++ {
				g.c <- struct{}{}
			}
			g.mu.Unlock()
		}
	}()

	//异步调用时需要的消费线程
	go func() { //起一个携程，消费队列
		for {
			<-g.c
			g.mu.Lock()
			if len(g.fns) > 0 {
				go g.doCall(g.fns[len(g.fns)-1].key, g.fns[len(g.fns)-1].fn)
				g.fns = g.fns[:len(g.fns)-1]
			}
			g.mu.Unlock()
		}
	}()
	return g
}
