package control

import (
	"github.com/EngoEngine/ecs"
	"github.com/EngoEngine/engo"
	"github.com/EngoEngine/engo/common"
	"github.com/Z2Y/trpgo/city/core"
	"github.com/Z2Y/trpgo/city/core/input"
	"github.com/Z2Y/trpgo/city/core/route"
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
			c.followHero()
			c.routeToMousePoint()
		}
		if !ok {
			break
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

func (c *ControlSystem) routeToMousePoint() {
	var (
		mouseX = ((engo.Input.Mouse.X * c.camera.Z() * engo.GameWidth() / engo.WindowWidth()) + (c.camera.X()-(engo.GameWidth()/2)*c.camera.Z())/engo.GetGlobalScale().X)
		mouseY = ((engo.Input.Mouse.Y * c.camera.Z() * engo.GameHeight() / engo.WindowHeight()) + (c.camera.Y()-(engo.GameHeight()/2)*c.camera.Z())/engo.GetGlobalScale().Y)
	)

	if engo.Input.Mouse.Action == engo.Press {
		gx, gy := c.grid.GetGridPos(mouseX, mouseY)

		footX, footY := c.grid.GetGridPos(c.hero.SpaceComponent.Position.X+c.hero.Offset.X, c.hero.SpaceComponent.Position.Y+c.hero.Offset.Y+c.hero.SpaceComponent.Height/2)
		router := &route.AstarRoute{World: c.grid}
		if c.hero.ActionState.Route != nil {
			c.hero.ActionState.Route.Parent = nil
			footX, footY = c.grid.GetGridPos(c.hero.ActionState.Route.Point.X, c.hero.ActionState.Route.Point.Y)
		}
		path := router.FindPath(engo.Point{X: float32(footX), Y: float32(footY)}, engo.Point{X: float32(gx), Y: float32(gy)})
		c.hero.ActionState.NextRoute = path.Reserve()
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
