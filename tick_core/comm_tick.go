package tick_core

type TickServiceModel interface {
	GetUniqueId() string
	OverDoFunc()
	StopDoFunc()
	EverSecondDoFunc()
}

type BaseTickServiceModel struct {
}

func (t BaseTickServiceModel) GetUniqueId() string {
	return ""
}

func (t BaseTickServiceModel) OverDoFunc() {
}

func (t BaseTickServiceModel) StopDoFunc() {
}

func (t BaseTickServiceModel) EverSecondDoFunc() {
}
