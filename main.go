package main

import (
	"fmt"
	"tick_server/tick_core"
	"time"
)

/*
目前项目存在的问题
1. 项目重启会导致计时器失效
2. 日志优化
3. timeTickMap 性能需要优化
*/

type TestTickServiceModel struct {
	tick_core.BaseTickServiceModel
	TestServiceData string
}

// GetUniqueId 业务侧控制全局唯一计时器id
func (t TestTickServiceModel) GetUniqueId() string {
	return fmt.Sprintf("test_unique_%s", t.TestServiceData)
}

// OverDoFunc 计时器计时完成后的业务逻辑
func (t TestTickServiceModel) OverDoFunc() {
	fmt.Println("TestTickServiceModel over do function", t.TestServiceData)
}

// StopDoFunc 业务侧终止计时器后的业务逻辑
func (t TestTickServiceModel) StopDoFunc() {
	fmt.Println("TestTickServiceModel stop do function", t.TestServiceData)
}

func main() {
	tickPool := tick_core.NewTickWithConf(1000, 10000)
	test := &TestTickServiceModel{TestServiceData: fmt.Sprintf("test_%d", 1)}
	baseTick := tick_core.NewBaseTick(test.GetUniqueId(), time.Second*5, tick_core.InitTickTime, test, tickPool)
	baseTick.StartTick()        //开启计时
	time.Sleep(time.Second * 7) //等待计时完成 并执行计时完成逻辑
	baseTick.StartTick()        //开启计时
	baseTick.StopTick()         //业务关闭计时器并执行关闭后的业务逻辑
	time.Sleep(time.Second * 2)
}
