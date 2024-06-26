package tick_core

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

type TestTickServiceModel struct {
	BaseTickServiceModel
	TestServiceData string
}

func (t TestTickServiceModel) GetUniqueId() string {
	return fmt.Sprintf("test_unique_%s", t.TestServiceData)
}

func (t TestTickServiceModel) OverDoFunc() {
	fmt.Println("TestTickServiceModel over do function", t.TestServiceData)
}

func (t TestTickServiceModel) StopDoFunc() {
	fmt.Println("TestTickServiceModel stop do function", t.TestServiceData)
}

func TestTick(t *testing.T) {
	tickPool := NewTickWithConf(10, 100)
	test := &TestTickServiceModel{TestServiceData: "test"}
	baseTick := BaseTick{
		Id:          test.GetUniqueId(),
		TrickTime:   time.Second * 5,
		ServiceData: test,
		TickPool:    tickPool,
	}
	t.Run("测试计时器 timeout", func(t *testing.T) {
		baseTick.StartTick()
		time.Sleep(2 * time.Second)
		stopTick, ok := timeTickMap.Load(baseTick.Id)
		assert.Equal(t, true, ok)
		assert.Equal(t, 0, len(stopTick.(chan struct{})))
		time.Sleep(8 * time.Second)
		_, ok = timeTickMap.Load(baseTick.Id)
		assert.Equal(t, false, ok)
		assert.Equal(t, true, true)
	})

	t.Run("测试计时器 stop操作", func(t *testing.T) {
		baseTick.StartTick()
		time.Sleep(1 * time.Second)
		stopTick, ok := timeTickMap.Load(baseTick.Id)
		assert.Equal(t, true, ok)
		assert.Equal(t, 0, len(stopTick.(chan struct{})))
		baseTick.StopTick()
		time.Sleep(2 * time.Second)
		_, ok = timeTickMap.Load(baseTick.Id)
		assert.Equal(t, false, ok)
	})
}

func TestIntervalTick(t *testing.T) {
	t.Run("测试间隔计时器", func(t *testing.T) {
		test2 := &TestTickServiceModel{TestServiceData: "test2"}
		baseTick2 := BaseTick{
			Id:           test2.GetUniqueId(),
			TrickTime:    time.Second * 5,
			ServiceData:  test2,
			IntervalTime: time.Second * 1,
		}
		baseTick2.StartTickOfEvery()
		time.Sleep(1 * time.Second)
		stopTick, ok := timeTickMap.Load(baseTick2.Id)
		assert.Equal(t, true, ok)
		assert.Equal(t, 0, len(stopTick.(chan struct{})))
		time.Sleep(8 * time.Second)
		_, ok = timeTickMap.Load(baseTick2.Id)
		assert.Equal(t, false, ok)
	})
}

func TestTickPerformance(t *testing.T) {
	tickPool := NewTickWithConf(10, 100)
	t.Run("测试timer 是否复用", func(t *testing.T) {
		for i := 0; i < 50; i++ {
			test := &TestTickServiceModel{TestServiceData: fmt.Sprintf("test_%d", i)}
			baseTick := BaseTick{
				Id:          test.GetUniqueId(),
				TrickTime:   time.Second * time.Duration(i%10+1),
				ServiceData: test,
				TickPool:    tickPool,
			}
			baseTick.StartTick()
		}
		for i := 0; i < 20; i++ {
			time.Sleep(time.Second * 1)
			fmt.Println(tickPool.GetFreeTimerListLen(), tickPool.runningNum)
		}
	})
}
