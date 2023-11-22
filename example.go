package tick

import (
	"fmt"
	"time"
)

type TreasureBoxRefresh struct {
	BaseTick
}

func (p *TreasureBoxRefresh) NewUniqueFlagForTick() *UniqueFlagForTick { //内部使用
	return &UniqueFlagForTick{Id: fmt.Sprintf(TickForExample, "example")}
}

func (p *TreasureBoxRefresh) StartTick() {
	go p.TickServe()
}

func (p *TreasureBoxRefresh) TickServe() {
	tick := p.NewUniqueFlagForTick()
	tick.StopTick()
	isTimeOver, err := tick.DoTick(time.Second * 10) //阻塞 开始计时10s
	if err == nil {
		if isTimeOver { //超时后处理逻辑
		} else { //计时终止后的处理逻辑

		}
	}
}
