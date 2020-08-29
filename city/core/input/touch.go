package input

import "github.com/EngoEngine/engo"

const TOUCH_MESSAGE = "TouchMessage"

type TouchMessage struct {
	X, Y   int
	ID     int
	Action int
}

func (TouchMessage) Type() string {
	return TOUCH_MESSAGE
}

type TouchHandler struct {
	touches chan *TouchMessage
}

func NewTouchHandler() *TouchHandler {
	handler := &TouchHandler{touches: make(chan *TouchMessage)}
	handler.listen()
	return handler
}

func (t *TouchHandler) Update() bool {
	select {
	case msg, ok := <-t.touches:
		if ok {
			UpdateMouse(msg.X, msg.Y, msg.ID, msg.Action)
		}
		return ok
	default:
		return false
	}
}

func (t *TouchHandler) listen() {
	engo.Mailbox.Listen(TOUCH_MESSAGE, func(message engo.Message) {
		msg, isTouch := message.(TouchMessage)
		if isTouch {
			t.touches <- &msg
		}
	})
}
