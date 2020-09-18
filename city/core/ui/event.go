package ui

type UIMouseEvent struct {
}

func (UIMouseEvent) Type() string {
	return "UIMouseEvent"
}

type UICloseEvent struct {
	Target *UIBasic
}

func (UICloseEvent) Type() string {
	return "UICloseEvent"
}
