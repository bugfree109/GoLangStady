package _Map

import (
	"fmt"
	"sync"
	"time"
)

type SafeMap struct {
	mu  sync.RWMutex
	m   map[int]interface{}
	rCh map[int]chan interface{}
}

func NewSafeMap() *SafeMap {
	safeMap := &SafeMap{
		m: make(map[int]interface{}, 0),
		//rCh: make(chan map[int]interface{}, 100),
		rCh: make(map[int]chan interface{}),
	}
	return safeMap
}

func (s *SafeMap) Set(key int, value interface{}) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.m[key] = value
	fmt.Printf("set Data key:%d,value:%v", key, value)
	if ch, exists := s.rCh[key]; exists {
		ch <- value
		close(ch)          // 关闭通道以避免泄漏
		delete(s.rCh, key) // 清理已通知的通道
	}
}

func (s *SafeMap) Get(key int, sleepTime uint64) (interface{}, bool) {
	s.mu.RLock()
	value, exists := s.m[key]
	s.mu.RUnlock()

	if exists {
		return value, true
	}

	s.mu.Lock()
	if _, exists := s.rCh[key]; !exists {
		s.rCh[key] = make(chan interface{}, 1)
	}
	ch := s.rCh[key]
	s.mu.Unlock()

	select {
	case update := <-ch:
		return update, true
	case <-time.After(time.Duration(sleepTime) * time.Second):
		fmt.Printf("time out key:%d\n", key)
		return nil, false
	}

	return nil, false
}
