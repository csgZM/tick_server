package main

import (
	"fmt"
	"tick_server/tick_core"
	"time"
)

/*
目前项目存在的问题
2. 项目重启会导致计时器失效
*/

type TestTickServiceModel struct {
	tick_core.BaseTickServiceModel
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

func main() {
	//f, _ := os.OpenFile("cpu.pprof", os.O_CREATE|os.O_RDWR, 0644)
	//defer f.Close()
	//pprof.StartCPUProfile(f)
	//defer pprof.StopCPUProfile()
	//defer profile.Start(profile.MemProfile, profile.MemProfileRate(1)).Stop()
	tickPool := tick_core.NewTickWithConf(1000, 10000)
	for i := 0; i < 10001; i++ {
		test := &TestTickServiceModel{TestServiceData: fmt.Sprintf("test_%d", i)}
		baseTick := tick_core.BaseTick{
			Id:          test.GetUniqueId(),
			TrickTime:   time.Second * time.Duration(i%10+5),
			ServiceData: test,
			TickPool:    tickPool,
		}
		baseTick.StartTick()
	}
	for i := 0; i < 20; i++ {
		time.Sleep(time.Second * 1)
		fmt.Println(tickPool.GetFreeTimerListLen(), tickPool.GetRunningNum())
	}
}
