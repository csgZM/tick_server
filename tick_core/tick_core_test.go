package tick_core

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"runtime"
	"testing"
	"time"
)

var curMem uint64

const (
	_   = 1 << (10 * iota)
	KiB // 1024
	MiB // 1048576
)

type TestTickServiceModel struct {
	BaseTickServiceModel
	TestServiceData string
}

func (t TestTickServiceModel) GetUniqueId() string {
	return fmt.Sprintf("test_unique_%s", t.TestServiceData)
}

func (t TestTickServiceModel) OverDoFunc() {
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
		tickPool := NewTickWithConf(10000, 10000)
		test := &TestTickServiceModel{TestServiceData: "test2"}
		baseTick := NewBaseTick(test.GetUniqueId(), time.Second*5, time.Second*1, test, tickPool)

		baseTick.StartTickOfEvery()
		time.Sleep(1 * time.Second)
		stopTick, ok := timeTickMap.Load(baseTick.Id)
		assert.Equal(t, true, ok)
		assert.Equal(t, 0, len(stopTick.(chan struct{})))
		time.Sleep(8 * time.Second)
		_, ok = timeTickMap.Load(baseTick.Id)
		assert.Equal(t, false, ok)
	})
}

func TestTickPerformance(t *testing.T) { //测试timer池的复用带来的内存优化比较
	t.Run("测试timer池使用后的内存使用情况", func(t *testing.T) {
		tickPool := NewTickWithConf(10000, 10000)
		for i := 0; i < 9999; i++ {
			test := &TestTickServiceModel{TestServiceData: fmt.Sprintf("test_%d", i)}
			baseTick := NewBaseTick(test.GetUniqueId(), time.Second*time.Duration(i%5+1), InitTickTime, test, tickPool)
			baseTick.StartTick()
		}

		time.Sleep(time.Second * 7)

		for i := 0; i < 9999; i++ {
			test := &TestTickServiceModel{TestServiceData: fmt.Sprintf("test_%d", i)}
			baseTick := NewBaseTick(test.GetUniqueId(), time.Second*time.Duration(i%5+1), InitTickTime, test, tickPool)
			baseTick.StartTick()
		}

		time.Sleep(time.Second * 7)

		mem := runtime.MemStats{}
		runtime.ReadMemStats(&mem)
		curMem = mem.Alloc
		t.Logf("memory = %vKB, GC Times = %vn", curMem, mem.NumGC)
	})

	t.Run("测试未使用timer池的内存使用情况", func(t *testing.T) {
		for i := 0; i < 9999; i++ {
			test := &TestTickServiceModel{TestServiceData: fmt.Sprintf("test_%d", i)}
			baseTick := NewBaseTick(test.GetUniqueId(), time.Second*time.Duration(i%5+1), InitTickTime, test, nil)
			baseTick.StartTick()
		}
		time.Sleep(time.Second * 7)
		for i := 0; i < 9999; i++ {
			test := &TestTickServiceModel{TestServiceData: fmt.Sprintf("test_%d", i)}
			baseTick := NewBaseTick(test.GetUniqueId(), time.Second*time.Duration(i%5+1), InitTickTime, test, nil)
			baseTick.StartTick()
		}
		time.Sleep(time.Second * 7)
		mem := runtime.MemStats{}
		runtime.ReadMemStats(&mem)
		curMem = mem.Alloc
		t.Logf("memory = %vKB, GC Times = %vn", curMem, mem.NumGC)
	})
}
