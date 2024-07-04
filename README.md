# tick_server

## 📖介绍
项目设计的目的是将计时功能与具体业务分离，业务侧无需关注计时相关逻辑代码，且支持业务侧对计时器进行关闭操作，同时对整个项目的计时器进行统一管理

## 🚀功能：

- 业务侧只需要实现TickServiceModel即可使用计时器功能，方便使用，可支持业务侧计时器计时结束、终止、计时内规定间隔时间触相关业务逻辑函数
- 实现timer池，用来进行timer管理和复用
- timer池可进行动态扩容和空闲timer释放

## 软件架构
### 流程图
![输入图片说明](%E9%A1%B9%E7%9B%AE%E6%B5%81%E7%A8%8B%E5%9B%BE-2024-07-03-1614.png)

## 🧰安装教程
go get github.com/csgZM/tick_server

## 🛠 使用说明

``` go

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
	tickPool := tick_core.NewTickWithConf(1000, 10000) //timer池配置
	test := &TestTickServiceModel{TestServiceData: fmt.Sprintf("test_%d", 1)}
	baseTick := tick_core.NewBaseTick(test.GetUniqueId(), time.Second*5, tick_core.InitTickTime, test, tickPool)
	baseTick.StartTick()        //开启计时
	time.Sleep(time.Second * 7) //等待计时完成 并执行计时完成逻辑
	baseTick.StartTick()        //开启计时
	baseTick.StopTick()         //业务关闭计时器并执行关闭后的业务逻辑
	time.Sleep(time.Second * 2)
}

```