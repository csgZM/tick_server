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
	InitTickTime        = time.Duration(0)
)

// TimeTickMap 用于映射计时器id和计时信号管道
var timeTickMap sync.Map

var blockBaseTick []*BaseTick //超过tick池的最大容量后阻塞的计时器存储列表
var blockBaseTickMutex sync.Mutex

type BaseTick struct {
	Id           string //保证唯一 不唯一会关闭前一个的计时器
	TrickTime    time.Duration
	IntervalTime time.Duration
	ServiceData  TickServiceModel
	TickPool     *TickWithConf
}

func NewBaseTick(id string, trickTime, intervalTime time.Duration, serviceData TickServiceModel, tickPool *TickWithConf) *BaseTick {
	return &BaseTick{
		Id:           id,
		TrickTime:    trickTime,
		IntervalTime: intervalTime,
		ServiceData:  serviceData,
		TickPool:     tickPool,
	}
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
	timer := p.TickPool.GetFreeTimer()
	if timer == nil {
		blockBaseTickMutex.Lock()
		blockBaseTick = append(blockBaseTick, p)
		blockBaseTickMutex.Unlock()
		return
	}
	timer.Reset(p.TrickTime)
	timeTickMap.Store(p.Id, make(chan struct{}))
	defer p.TickPool.PushToFreeTimerList(timer)
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

func GetBlockBaseTickLen() int {
	blockBaseTickMutex.Lock()
	defer blockBaseTickMutex.Unlock()
	return len(blockBaseTick)
}
