package ui

import (
	"github.com/EngoEngine/ecs"
	"github.com/EngoEngine/engo/common"
)

type UIBasic struct {
	ecs.BasicEntity
	common.RenderComponent
	common.SpaceComponent
}

type UISystem struct {
	entities []*UIBasic

	renderer *common.RenderSystem
}

func (ui *UISystem) New(world *ecs.World) {
	for _, system := range world.Systems() {
		switch sys := system.(type) {
		case *common.RenderSystem:
			ui.renderer = sys
		}
	}
}

func (ui *UISystem) Add(entity *UIBasic) {
	ui.entities = append(ui.entities, entity)
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
	}
}

func (ui *UISystem) Update(dt float32) {
	for _, e := range ui.entities {
		ui.renderer.Add(&e.BasicEntity, &e.RenderComponent, &e.SpaceComponent)
	}
}
