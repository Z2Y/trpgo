//+build android OR ios

package input

import "github.com/EngoEngine/engo"

func TouchEvent(x, y, id, action int) {
	engo.TouchEvent(x, y, id, action)
	engo.Mailbox.Dispatch(TouchMessage{X: float32(x), Y: float32(y), ID: id, Action: engo.Input.Mouse.Action})
	engo.Input.Mouse.Action = engo.Neutral
}
