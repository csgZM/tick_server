package tick

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

type UniqueFlagForTick struct {
	Id string
}

// TimeTickMap 用于映射计时器id和计时信号管道
var TimeTickMap sync.Map

// DoTick 延时操作
func (t *UniqueFlagForTick) DoTick(duration time.Duration) (IsTimeOver bool, err error) {
	timer := time.NewTimer(duration)
	TimeTickMap.Store(t.Id, make(chan struct{}))
	defer timer.Stop()
	for {
		stopTick, ok := TimeTickMap.Load(t.Id)
		if !ok {
			return false, errors.New("计时器异常退出")
		}
		select {
		case <-timer.C:
			fmt.Println("Time out!")
			TimeTickMap.Delete(t.Id)
			return true, nil
		case <-stopTick.(chan struct{}):
			TimeTickMap.Delete(t.Id)
			return false, nil
		}
	}
}

// DoTickForEvery 支持每秒操作函数
func (t *UniqueFlagForTick) DoTickForEvery(duration time.Duration, f func()) (IsTimeOver bool, err error) {
	timer := time.NewTimer(duration)
	TimeTickMap.Store(t.Id, make(chan struct{}))
	ticker := time.NewTicker(time.Second * 1) //定义一个1秒间隔的定时器
	defer ticker.Stop()
	defer timer.Stop()
	for {
		stopTick, ok := TimeTickMap.Load(t.Id)
		if !ok {
			return false, errors.New("计时器异常退出")
		}
		select {
		case <-timer.C:
			TimeTickMap.Delete(t.Id)
			return true, nil
		case <-ticker.C:
			f()
		case <-stopTick.(chan struct{}):
			TimeTickMap.Delete(t.Id)
			return false, nil
		}
	}
}

func (t *UniqueFlagForTick) StopTick() {
	stopTick, ok := TimeTickMap.Load(t.Id)
	if ok {
		stopTick.(chan struct{}) <- struct{}{}
	}
}
