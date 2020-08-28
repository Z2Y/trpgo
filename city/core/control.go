package core

import (
	"github.com/EngoEngine/ecs"
	"github.com/EngoEngine/engo"
	"github.com/EngoEngine/engo/common"
	"github.com/Z2Y/trpgo/city/core/input"
)

type ControlSystem struct {
	grid *WorldSystem

	camera *common.CameraSystem

	state     ControlState
	ZoomSpeed float32
}

type ControlState struct {
	x float32
	y float32

	mouseX float32
	mouseY float32

	dx float32
	dy float32

	mouseDown bool
}

var (
	Steps = []float32{0, 1, 2, 4, 8, 16, 32, 64}
)

func NextStep(i float32) float32 {
	s, p := i, float32(1)
	if i < 0 {
		s, p = -i, float32(-1)
	}
	for _, x := range Steps {
		if s <= x {
			return x * p
		}
	}
	return Steps[len(Steps)-1] * p
}

func (c *ControlSystem) setWorldCamera() {
	worldWidth, worldHeight := c.grid.Size()
	centerX := float32(worldWidth) / 2
	centerY := float32(0)

	common.CameraBounds.Min.X = centerX - worldWidth
	common.CameraBounds.Min.Y = centerY - worldHeight

	common.CameraBounds.Max.X = centerX + worldWidth
	common.CameraBounds.Max.Y = centerY + worldWidth

	engo.Mailbox.Dispatch(common.CameraMessage{Axis: common.XAxis,
		Value:       centerX,
		Incremental: false})
	engo.Mailbox.Dispatch(common.CameraMessage{Axis: common.YAxis,
		Value:       centerY,
		Incremental: false})
}

func (c *ControlSystem) Remove(basic ecs.BasicEntity) {
	// do nothing
}

func (c *ControlSystem) Update(float32) {
	var (
		mouseX = engo.Input.Mouse.X // *c.camera.Z() + (c.camera.X()-(engo.GameWidth()/2)*c.camera.Z()+(engo.ResizeXOffset/2))/engo.GetGlobalScale().X
		mouseY = engo.Input.Mouse.Y // *c.camera.Z() + (c.camera.Y()-(engo.GameHeight()/2)*c.camera.Z()+(engo.ResizeYOffset/2))/engo.GetGlobalScale().Y
	)

	switch engo.Input.Mouse.Action {
	case engo.Press:
		c.state.mouseX, c.state.mouseY = mouseX, mouseY
		c.state.mouseDown = true
	case engo.Move:
		if c.state.mouseDown {
			c.state.dx = c.state.mouseX - mouseX
			c.state.dy = c.state.mouseY - mouseY
			c.state.x = c.state.x + c.state.dx
			c.state.y = c.state.y + c.state.dy
			c.state.mouseX, c.state.mouseY = mouseX, mouseY
		}
	case engo.Release:
		c.state.dx, c.state.dy = 0, 0
		c.state.mouseDown = false
	}

	engo.Mailbox.Dispatch(common.CameraMessage{Axis: common.XAxis,
		Value:       c.state.dx,
		Incremental: true})
	engo.Mailbox.Dispatch(common.CameraMessage{Axis: common.YAxis,
		Value:       c.state.dy,
		Incremental: true})

	if engo.Input.Mouse.ScrollY != 0 {
		engo.Mailbox.Dispatch(common.CameraMessage{Axis: common.ZAxis, Value: (engo.Input.Mouse.ScrollY * c.ZoomSpeed), Incremental: true})
	}
}

func (c *ControlSystem) New(w *ecs.World) {
	for _, system := range w.Systems() {
		switch sys := system.(type) {
		case *common.CameraSystem:
			c.camera = sys
		case *WorldSystem:
			c.grid = sys
		}
	}

	if c.grid != nil && c.camera != nil {
		c.setWorldCamera()
	}
	c.listenTouchEvent()
}

func (s *ControlSystem) listenTouchEvent() {
	engo.Mailbox.Listen(input.TOUCH_MESSAGE, func(message engo.Message) {
		_, isTouch := message.(input.TouchMessage)
		if isTouch {
			s.Update(engo.Time.Delta())
		}
	})
}
