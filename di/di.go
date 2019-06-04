package di

import "github.com/dssysolyatin/ssecret/ui"

var _ DI = (*di)(nil)

type DI interface {
	GetUI() ui.UI
}

type di struct {
	ui ui.UI
}

func NewDI() *di {
	return &di{}
}

func (d *di) GetUI() ui.UI {
	if d.ui == nil {
		d.ui = ui.NewConsoleUI()
	}
	return d.ui
}
