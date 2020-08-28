package input

import "github.com/EngoEngine/engo"

const TOUCH_MESSAGE = "TouchMessage"

type TouchMessage struct {
	X, Y   float32
	ID     int
	Action engo.Action
}

func (TouchMessage) Type() string {
	return TOUCH_MESSAGE
}
