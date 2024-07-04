# tick_model

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


## 🛠 使用说明

``` go

func main() {
	tickPool := tick_core.NewTickWithConf(1000, 10000)
	test := &TestTickServiceModel{TestServiceData: fmt.Sprintf("test_%d", 1)}
	baseTick := tick_core.NewBaseTick(test.GetUniqueId(), time.Second*5, tick_core.InitTickTime, test, tickPool)
	baseTick.StartTick()
	time.Sleep(time.Second * 7)
	baseTick.StartTick() //开启计时
	baseTick.StopTick()
	time.Sleep(time.Second * 2)
}

```