package tick_core

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

type TickWithConf struct {
	freeTimerList []*time.Timer //空闲的计时器存储列表
	effectiveNum  int32         //计时器数量超出改值将空闲的进行释放
	capacity      int32         //timer 最大容量
	runningNum    int32         //运行的计时器数量
	sync.Mutex
}

func (c *TickWithConf) GetRunningNum() int32 {
	c.Lock()
	defer c.Unlock()
	return c.runningNum
}

func (c *TickWithConf) GetFreeTimerListLen() int {
	c.Lock()
	defer c.Unlock()
	return len(c.freeTimerList)
}

func (c *TickWithConf) GetFreeTimer() *time.Timer {
	var tempTimer *time.Timer
	c.Lock()
	if len(c.freeTimerList) <= 0 {
		if c.runningNum > c.capacity {
			fmt.Println("超出规定最大可开启的计时器数量")
			return nil
		}
		tempTimer = c.newTimer()
	} else {
		tempTimer = c.freeTimerList[0]
		c.freeTimerList[0] = nil
		if len(c.freeTimerList) > 1 {
			c.freeTimerList = c.freeTimerList[1:]
		} else {
			c.freeTimerList = c.freeTimerList[0:0]
		}

	}
	c.Unlock()
	atomic.AddInt32(&c.runningNum, 1)
	if !tempTimer.Stop() {
		select {
		case <-tempTimer.C:
		default:
		}
	}
	return tempTimer
}

func (c *TickWithConf) PushToFreeTimerList(timer *time.Timer) {
	timer.Stop()
	c.Lock()
	atomic.AddInt32(&c.runningNum, -1)
	c.freeTimerList = append(c.freeTimerList, timer)
	c.Unlock()
}

func (c *TickWithConf) newTimer() *time.Timer {
	return time.NewTimer(InitTickTime)
}

func (c *TickWithConf) releaseTimer() {
	ticker := time.NewTicker(time.Second * 5)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			c.Lock()
			if len(c.freeTimerList) > int(c.effectiveNum) {
				for i := int(c.effectiveNum); i < len(c.freeTimerList); i++ {
					c.freeTimerList[i].Stop()
					c.freeTimerList[i] = nil
				}
				c.freeTimerList = c.freeTimerList[:int(c.effectiveNum)]
			}
			c.Unlock()
		}

	}
}

func NewTickWithConf(effectiveNum, capacity int32) *TickWithConf {
	tickWithConf := &TickWithConf{
		freeTimerList: NewFreeTimerList(effectiveNum),
		effectiveNum:  effectiveNum,
		capacity:      capacity,
	}
	go tickWithConf.releaseTimer()
	return tickWithConf
}

func NewFreeTimerList(effectiveNum int32) []*time.Timer {
	tempFreeTimerList := make([]*time.Timer, effectiveNum)
	for i := 0; i < int(effectiveNum); i++ {
		tempFreeTimerList[i] = time.NewTimer(InitTickTime)
	}
	return tempFreeTimerList
}
