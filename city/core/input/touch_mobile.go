//+build android OR ios

package input

import (
	"github.com/EngoEngine/engo"
)

func TouchEvent(x, y, id, action int) {
	go func() {
		engo.Mailbox.Dispatch(TouchMessage{X: x, Y: y, ID: id, Action: action})
	}()
}

func UpdateMouse(x, y, id, action int) {
	engo.TouchEvent(x, y, id, action)
}
