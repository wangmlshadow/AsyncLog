package LogModule

import (
	"sync"
)

type CountDownLatch struct {
	mutex *sync.Mutex
	cond *sync.Cond
	count int
}

func NewCountDownLatch(count int) *CountDownLatch {
	mutex := &sync.Mutex{}
	cond := sync.NewCond(mutex)
	return &CountDownLatch{mutex: mutex, cond: cond, count: count}
}

func (c *CountDownLatch) Wait() {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	for c.count > 0 {
		c.cond.Wait()
	}
}

func (c *CountDownLatch) CountDown() {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.count -= 1
	if c.count == 0 {
		c.cond.Broadcast()
	}
}