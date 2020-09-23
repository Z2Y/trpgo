package ui

import (
	"math"
	"sort"

	"github.com/EngoEngine/ecs"
	"github.com/EngoEngine/engo"
	"github.com/EngoEngine/engo/common"
	"github.com/Z2Y/trpgo/city/config"
	"github.com/Z2Y/trpgo/city/core/input"
)

var (
	UIEntityFace     *UIFace
	UILayerIndex     = float32(1e4)
	UIPopLayerIndex  = float32(2e4)
	UISystemPriority = 1000
)

type UIBasic struct {
	ecs.BasicEntity
	common.RenderComponent
	common.SpaceComponent

	Width           float32
	Height          float32
	SubEntities     []*UIBasic
	MessageListener *engo.MessageManager

	EventResponder bool //是否处理事件
}

func (u *UIBasic) GetComponentUI() *UIBasic {
	return u
}

func (u *UIBasic) AddSubEntity(e UIFace) {
	u.SubEntities = append(u.SubEntities, e.GetComponentUI())
}

func (u *UIBasic) SetPosition(pos engo.Point) {
	if len(u.SubEntities) > 0 {
		offset := engo.Point{X: u.Position.X, Y: u.Position.Y}
		offset.Subtract(pos)
		for _, sub := range u.SubEntities {
			sub.SetPosition(engo.Point{X: sub.Position.X - offset.X, Y: sub.Position.Y - offset.Y})
		}
	}

	u.Position = pos
}

func (u *UIBasic) SetZIndex(idx float32) {
	if len(u.SubEntities) > 0 {
		offset := u.StartZIndex - idx
		for _, sub := range u.SubEntities {
			sub.SetZIndex(sub.StartZIndex - offset)
		}
	}

	u.StartZIndex = idx
	u.RenderComponent.SetZIndex(idx)
}

func (u *UIBasic) OnClick(f func()) {
	if !u.EventResponder {
		return
	}
	if u.MessageListener == nil {
		u.MessageListener = &engo.MessageManager{}
	}
	u.MessageListener.Listen("UIMouseEvent", func(engo.Message) {
		f()
	})
}

type UIFace interface {
	GetComponentUI() *UIBasic
	AddSubEntity(UIFace)
}

type respondEntityList []*UIBasic

func (r respondEntityList) Len() int {
	return len(r)
}

func (r respondEntityList) Less(i, j int) bool {
	if r[i].StartZIndex != r[j].StartZIndex {
		return r[i].StartZIndex > r[j].StartZIndex
	}
	if r[i].Position.Y != r[j].Position.Y {
		return r[i].Position.Y > r[j].Position.Y
	}

	return r[i].Position.X > r[j].Position.X
}

func (r respondEntityList) Swap(i, j int) {
	r[i], r[j] = r[j], r[i]
}

type UISystem struct {
	entities        map[uint64]*UIBasic
	respondEntities respondEntityList

	touchHandler *input.TouchHandler

	renderer *common.RenderSystem
	camera   *common.CameraSystem

	cameraY       float32
	sortingNeeded bool
}

func (ui *UISystem) New(world *ecs.World) {
	for _, system := range world.Systems() {
		switch sys := system.(type) {
		case *common.RenderSystem:
			ui.renderer = sys
		case *common.CameraSystem:
			ui.camera = sys
		}
	}
	ui.entities = make(map[uint64]*UIBasic)
	ui.touchHandler = input.NewTouchHandler()
	engo.Mailbox.Listen("UICloseEvent", func(message engo.Message) {
		msg, ok := message.(UICloseEvent)
		if ok {
			ui.Remove(msg.Target.BasicEntity)
		}
	})
	engo.Mailbox.Listen("WindowResizeMessage", config.UpdateWindowScale)
}

func (*UISystem) Priority() int {
	return UISystemPriority
}

func (ui *UISystem) Add(entity *UIBasic) {
	if _, ok := ui.entities[entity.ID()]; ok {
		return
	}
	ui.entities[entity.ID()] = entity

	ui.renderer.Add(&entity.BasicEntity, &entity.RenderComponent, &entity.SpaceComponent)

	for _, sub := range entity.SubEntities {
		ui.Add(sub)
	}

	if entity.EventResponder {
		ui.respondEntities = append(ui.respondEntities, entity)
		ui.sortingNeeded = true
	}
}

func (ui *UISystem) AddByInterface(o ecs.Identifier) {
	face := o.(UIFace)
	ui.Add(face.GetComponentUI())
}

func (ui *UISystem) removeRespondEntity(basic ecs.BasicEntity) {
	delete := -1
	for index, e := range ui.respondEntities {
		if e.ID() == basic.ID() {
			delete = index
			break
		}
	}
	if delete >= 0 {
		ui.respondEntities = append(ui.respondEntities[:delete], ui.respondEntities[delete+1:]...)
	}
}

func (ui *UISystem) Remove(basic ecs.BasicEntity) {

	entity, ok := ui.entities[basic.ID()]

	if ok {
		delete(ui.entities, basic.ID())
		for _, sub := range entity.SubEntities {
			ui.Remove(sub.BasicEntity)
		}
		ui.renderer.Remove(basic)

		if entity.EventResponder {
			ui.removeRespondEntity(basic)
			ui.sortingNeeded = true
		}
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

func (ui *UISystem) shouldUpdateIndex() bool {
	if ui.camera.Y() < 0 {
		return false
	}
	return math.Abs(float64(ui.camera.Y()-ui.cameraY)) > float64(engo.WindowHeight()/2)
}

func (ui *UISystem) getUIZIndexOffset() float32 {
	return (ui.camera.Y() + engo.WindowHeight()) / engo.GetGlobalScale().Y
}

func (ui *UISystem) update() {
	curPos := engo.Point{X: engo.Input.Mouse.X / config.WindowScale().X, Y: engo.Input.Mouse.Y / config.WindowScale().Y}
	updateZindex, zIndexOffset := ui.shouldUpdateIndex(), float32(0)

	resetInput := false

	if updateZindex {
		ui.cameraY = ui.camera.Y()
		zIndexOffset = ui.getUIZIndexOffset()
		for _, e := range ui.entities {
			e.SetZIndex(e.StartZIndex + zIndexOffset)
		}
	}

	if ui.sortingNeeded {
		sort.Sort(ui.respondEntities)
		ui.sortingNeeded = false
	}

	if engo.Input.Mouse.Action == engo.Press && engo.Input.Mouse.Button == engo.MouseButtonLeft {
		for _, e := range ui.respondEntities {
			if e.Contains(curPos) {
				if e.MessageListener != nil {
					e.MessageListener.Dispatch(UIMouseEvent{})
				}

				resetInput = true
				break
			}
		}
	}

	if resetInput {
		engo.Input.Mouse.Action = engo.Neutral
	}
}
