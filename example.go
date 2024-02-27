package tick

import (
	"fmt"
	"time"
)

type TreasureBoxRefresh struct {
	Info string
	BaseTick
}

func NewTreasureBoxRefresh(info string) *TreasureBoxRefresh {
	return &TreasureBoxRefresh{
		Info: info,
		BaseTick: BaseTick{
			TrickKey:  fmt.Sprintf(TickForExample, info),
			TrickTime: time.Second * 10,
		},
	}
}

func (p *TreasureBoxRefresh) StartTick() {
	go p.TickServe()
}

func (p *TreasureBoxRefresh) TickServe() {
	tick := p.BaseTick.NewUniqueFlagForTick()
	tick.StopTick()
	isTimeOver, err := tick.DoTick() //阻塞 开始计时10s
	if err == nil {
		if isTimeOver { //超时后处理逻辑
		} else { //计时终止后的处理逻辑

		}
	}
}

func (p *TreasureBoxRefresh) StopTick() {
	tick := p.BaseTick.NewUniqueFlagForTick()
	tick.StopTick()
}

/*
test git revert
*/
