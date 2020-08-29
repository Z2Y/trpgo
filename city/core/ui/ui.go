package ui

import (
	"github.com/EngoEngine/ecs"
	"github.com/EngoEngine/engo"
	"github.com/EngoEngine/engo/common"
	"github.com/Z2Y/trpgo/city/core/input"
)

var (
	UIEntityFace *UIFace
	UILayerIndex = float32(99)
)

type UIBasic struct {
	ecs.BasicEntity
	common.RenderComponent
	common.SpaceComponent

	MessageListener *engo.MessageManager
}

func (u *UIBasic) GetComponentUI() *UIBasic {
	return u
}

type UIFace interface {
	GetComponentUI() *UIBasic
}

type UISystem struct {
	entities     []*UIBasic
	touchHandler *input.TouchHandler

	renderer *common.RenderSystem
}

func (ui *UISystem) New(world *ecs.World) {
	for _, system := range world.Systems() {
		switch sys := system.(type) {
		case *common.RenderSystem:
			ui.renderer = sys
		}
	}
	ui.touchHandler = input.NewTouchHandler()
}

func (ui *UISystem) Add(entity *UIBasic) {
	ui.entities = append(ui.entities, entity)

	ui.renderer.Add(&entity.BasicEntity, &entity.RenderComponent, &entity.SpaceComponent)
}

func (ui *UISystem) AddByInterface(o ecs.Identifier) {
	face := o.(UIFace)
	ui.Add(face.GetComponentUI())
}

func (ui *UISystem) Remove(basic ecs.BasicEntity) {
	delete := -1
	for index, e := range ui.entities {
		if e.ID() == basic.ID() {
			delete = index
			break
		}
	}
	if delete >= 0 {
		ui.entities = append(ui.entities[:delete], ui.entities[delete+1:]...)
		ui.renderer.Remove(basic)
	}
}

func (ui *UISystem) Update(dt float32) {
	for {
		ok := ui.touchHandler.Update()
		ui.update()
		if !ok {
			break
		}
	}
}

func (ui *UISystem) update() {
	curPos := engo.Point{X: engo.Input.Mouse.X, Y: engo.Input.Mouse.Y}
	for _, e := range ui.entities {
		if e.Contains(curPos) && e.MessageListener != nil {
			if engo.Input.Mouse.Action == engo.Press && engo.Input.Mouse.Button == engo.MouseButtonLeft {
				e.MessageListener.Dispatch(UIMouseEvent{})
			}
		}
	}
}
