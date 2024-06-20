package tick_core

import (
	"fmt"
	"sync"
	"time"
)

var (
	DefaultTickTime     = time.Second * 30
	DefaultIntervalTime = time.Second * 1
	CloseDeadLineTime   = time.Millisecond * 100
)

// TimeTickMap 用于映射计时器id和计时信号管道
var timeTickMap sync.Map

type BaseTick struct {
	Id           string //保证唯一 不唯一会关闭前一个的计时器
	TrickTime    time.Duration
	IntervalTime time.Duration
	ServiceData  TickServiceModel
}

type TickWithConf struct {
	timeTickMap   map[string]chan struct{} //计时器信号控制map
	timerFreeList []*time.Timer            //空闲的计时器存储列表
	effectiveNum  int32                    //计时器数量超出改值将空闲的进行释放
	capacity      int32                    //timer 最大容量
	runningNum    int32                    //运行的计时器数量
	sync.Mutex
}

func (c *TickWithConf) GetFreeTimer() *time.Timer {
	if len(c.timerFreeList) <= 0 {
		if c.runningNum > c.capacity {
			fmt.Println("超出规定最大可开启的计时器数量")
			return nil
		}
		return c.NewTimer()
	}
	tempTimer := c.timerFreeList[0]
	c.Lock()
	if len(c.timerFreeList) > 1 {
		c.timerFreeList = c.timerFreeList[1:]
	} else {
		c.timerFreeList = c.timerFreeList[0:0]
	}
	c.Unlock()
	return tempTimer
}

func (c *TickWithConf) NewTimer() *time.Timer {
	return time.NewTimer(0)
}

func (p *BaseTick) StartTick() {
	if p.TrickTime == 0 {
		fmt.Println("warn:计时器时间为0，使用默认时间30s serviceId is ", p.Id)
		p.TrickTime = DefaultTickTime
	}
	p.stopTick() // 如果已经存在，则先关闭旧的计时器
	go p.doTick()
}

func (p *BaseTick) StopTick() {
	p.stopTick()
}

func (p *BaseTick) StartTickOfEvery() {
	if p.TrickTime == 0 {
		fmt.Println("warn:计时器时间为0，使用默认时间30s serviceId is ", p.Id)
		p.TrickTime = DefaultTickTime
	}
	if p.IntervalTime == 0 {
		fmt.Println("warn:间隙计时器间隙时间为0，使用默认间隙时间1s serviceId is ", p.Id)
		p.IntervalTime = DefaultIntervalTime
	}
	p.stopTick() // 如果已经存在，则先关闭旧的计时器
	go p.doTickForEvery()
}

// DoTick 延时操作
func (p *BaseTick) doTick() {
	timer := time.NewTimer(p.TrickTime)
	timeTickMap.Store(p.Id, make(chan struct{}))
	defer timer.Stop()
	for {
		stopTick, ok := timeTickMap.Load(p.Id)
		if !ok {
			fmt.Println("计时器异常退出")
		}
		select {
		case <-timer.C:
			p.ServiceData.OverDoFunc()
			timeTickMap.Delete(p.Id)
			return
		case <-stopTick.(chan struct{}):
			p.ServiceData.StopDoFunc()
			timeTickMap.Delete(p.Id)
			return
		}
	}
}

// DoTickForEvery 支持每秒操作函数
func (p *BaseTick) doTickForEvery() {
	timer := time.NewTimer(p.TrickTime)
	timeTickMap.Store(p.Id, make(chan struct{}))
	ticker := time.NewTicker(time.Second * 1) //定义一个1秒间隔的定时器
	defer ticker.Stop()
	defer timer.Stop()
	for {
		stopTick, ok := timeTickMap.Load(p.Id)
		if !ok {
			fmt.Println("计时器异常退出")
		}
		select {
		case <-timer.C:
			p.ServiceData.OverDoFunc()
			timeTickMap.Delete(p.Id)
			return
		case <-ticker.C:
			p.ServiceData.EverSecondDoFunc()
		case <-stopTick.(chan struct{}):
			p.ServiceData.StopDoFunc()
			timeTickMap.Delete(p.Id)
			return
		}
	}
}

func (p *BaseTick) stopTick() {
	stopTick, ok := timeTickMap.Load(p.Id)
	timer := time.NewTimer(CloseDeadLineTime)
	defer timer.Stop()
	if ok {
		select {
		case stopTick.(chan struct{}) <- struct{}{}:
			return
		case <-timer.C:
			timeTickMap.Delete(p.Id) //超时处理，判断为异常情况导致map的key未删除
			return
		}
	}
}
