package control

import (
	"github.com/EngoEngine/ecs"
	"github.com/EngoEngine/engo"
	"github.com/EngoEngine/engo/common"
	"github.com/Z2Y/trpgo/city/core"
	"github.com/Z2Y/trpgo/city/core/input"
)

const (
	UpButton    = "up"
	LeftButton  = "left"
	RightButton = "right"
	DownButton  = "down"
)

type ControlSystem struct {
	grid *core.WorldSystem
	hero *ActionEntity

	camera       *common.CameraSystem
	touchHandler *input.TouchHandler

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
	center := c.grid.Center()

	common.CameraBounds.Min.X = center.X - worldWidth/2
	common.CameraBounds.Min.Y = center.Y - worldHeight/2

	common.CameraBounds.Max.X = center.X + worldWidth/2
	common.CameraBounds.Max.Y = center.Y + worldHeight/2
}

func (c *ControlSystem) setupButton() {
	engo.Input.RegisterButton(UpButton, engo.KeyW, engo.KeyArrowUp)
	engo.Input.RegisterButton(LeftButton, engo.KeyA, engo.KeyArrowLeft)
	engo.Input.RegisterButton(RightButton, engo.KeyD, engo.KeyArrowRight)
	engo.Input.RegisterButton(DownButton, engo.KeyS, engo.KeyArrowDown)
}

func (c *ControlSystem) Remove(basic ecs.BasicEntity) {
	c.hero = nil
}

func (c *ControlSystem) Add(entity *ActionEntity) {
	c.hero = entity
}

func (c *ControlSystem) Update(dt float32) {
	for {
		ok := c.touchHandler.Update()
		if c.hero != nil {
			c.updateHero(dt)
			c.followHero()
		}
		if !ok {
			break
		}
	}
}

func (c *ControlSystem) updateHero(dt float32) {
	speed := engo.Point{X: 0, Y: 0}
	if engo.Input.Button(UpButton).Down() {
		speed.Y -= 1
	}
	if engo.Input.Button(LeftButton).Down() {
		speed.X -= 1
	}
	if engo.Input.Button(RightButton).Down() {
		speed.X += 1
	}
	if engo.Input.Button(DownButton).Down() {
		speed.Y += 1
	}

	preState := c.hero.ActionState.Code
	c.hero.ActionState.elapsed += dt

	if speed.X != c.hero.ActionState.Speed.X || speed.Y != c.hero.ActionState.Speed.Y {
		if speed.X == 0 && speed.Y == 0 {
			c.hero.ActionState.Code = ActionIdle
		} else {
			c.hero.ActionState.Code = ActionWalking
		}
		c.hero.ActionState.elapsed = 0
		c.hero.ActionState.duration = 0.2
		c.hero.ActionState.Speed = speed
	}

	switch c.hero.ActionState.Code {
	case ActionIdle:
		if preState != ActionIdle {
			engo.Mailbox.Dispatch(ActionMessage{BasicEntity: &c.hero.BasicEntity, State: c.hero.ActionState})
		}
	case ActionWalking:
		if c.hero.ActionState.elapsed == 0 {
			engo.Mailbox.Dispatch(ActionMessage{BasicEntity: &c.hero.BasicEntity, State: c.hero.ActionState})
		}
	}
}

func (c *ControlSystem) followHero() {
	engo.Mailbox.Dispatch(common.CameraMessage{Axis: common.XAxis,
		Value:       c.hero.SpaceComponent.Position.X + c.hero.Offset.X,
		Incremental: false})
	engo.Mailbox.Dispatch(common.CameraMessage{Axis: common.YAxis,
		Value:       c.hero.SpaceComponent.Position.Y + c.hero.Offset.Y,
		Incremental: false})

	if engo.Input.Mouse.ScrollY != 0 {
		engo.Mailbox.Dispatch(common.CameraMessage{Axis: common.ZAxis, Value: (engo.Input.Mouse.ScrollY * c.ZoomSpeed), Incremental: true})
	}

}

func (c *ControlSystem) TrackCamera() {
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
		case *core.WorldSystem:
			c.grid = sys
		}
	}

	if c.grid != nil && c.camera != nil {
		c.setWorldCamera()
	}
	c.setupButton()
	c.touchHandler = input.NewTouchHandler()
}
