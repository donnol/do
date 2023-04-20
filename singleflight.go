package do

import (
	"fmt"
	"sync"
)

var (
	singleFlightMap = NewMap[string, *singleFlight]()

	// 应该是一个key对应一个wg，不能全局共用一个wg
	wgMap = NewMap[string, *sync.WaitGroup]()
)

type singleFlight struct {
	val any
	err error
}

type SingleFlightCall[R any] func() (R, error)

// SingleFlight make sure only one request is doing with one key
func SingleFlight[R any](key string, fn SingleFlightCall[R]) (r R, err error) {
	wg := initWg(key)

	// 已经有请求在执行时等待其结果返回
	if c, ok := singleFlightMap.Lookup(key); ok {
		wg.Wait()
		removeWg(key)

		return c.val.(R), c.err
	}

	// 执行
	c := &singleFlight{}
	wg.Add(1)
	singleFlightMap.Insert(key, c)

	func() {
		defer func() {
			if v := recover(); v != nil {
				c.err = fmt.Errorf("single flight run err: %v", v)
			}
			wg.Done()
		}()
		r, err := fn()
		c.val, c.err = r, err
	}()

	return c.val.(R), c.err
}

func ForgotKey(key string) {
	singleFlightMap.Remove(key)
	removeWg(key)
}

func initWg(key string) *sync.WaitGroup {

	var wg *sync.WaitGroup
	if v, ok := wgMap.Lookup(key); ok {
		wg = v
	} else {
		wg = new(sync.WaitGroup)
		wgMap.Insert(key, wg)
	}

	return wg
}

func removeWg(key string) {
	wgMap.Remove(key)
}
