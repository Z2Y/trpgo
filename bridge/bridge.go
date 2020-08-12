//+build mobilebind

package bridge

import (
	"github.com/EngoEngine/engo"
	"github.com/Z2Y/trpgo/city"
)

var running bool

func Start(width, height int) {
	running = true
	city.Start(width, height)
}

func Update() {
	engo.RunIteration()
}

func IsRunning() bool {
	return running
}

func Touch(x, y, id, action int) {
	engo.TouchEvent(x, y, id, action)
}

func Resume() {
	running = true
}

func Stop() {
	running = false
	engo.MobileStop()
}
