package tick


const (
	TickForExample = "TickForExample:%s"
)

type IOperateTick interface {
	StartTick()
	TickServe()
	StopTick()
	NewUniqueFlagForTick() *UniqueFlagForTick
	TimeOverDo()
}

/*
NewXSTickInfo内没有开启计时器条件,bool可以一直为true,计时器支持重新计时
*/
type BaseTick struct{}

func (p *BaseTick) NewUniqueFlagForTick() {}

func (p *BaseTick) StartTick() {}

func (p *BaseTick) TickServe() {}

func (p *BaseTick) StopTick() {}

func (p *BaseTick) TimeOverDo() {}

func NewXSTickInfo(uid uint32) (IOperateTick, bool) {

}
