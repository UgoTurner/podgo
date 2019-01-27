package ui

type Panel struct {
	Title                                         string
	Name                                          string
	Highlight, Frame, Overwrite, Hidden, Editable bool
	Coordinate                                    Coordinate
	SelectionColor                                SelectionColor
}

func (p *Panel) EnableSelection() {
	p.SelectionColor.BgColorCurrent = p.SelectionColor.BgColorActive
	p.SelectionColor.FgColorCurrent = p.SelectionColor.FgColorActive
}

func (p *Panel) DisableSelection() {
	p.SelectionColor.BgColorCurrent = p.SelectionColor.BgColorUnactive
	p.SelectionColor.FgColorCurrent = p.SelectionColor.FgColorUnactive
}
