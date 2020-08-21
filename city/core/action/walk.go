package action

import (
	"github.com/EngoEngine/ecs"
	"github.com/EngoEngine/engo"
	"github.com/EngoEngine/engo/common"
	"github.com/Z2Y/trpgo/city/core"
)

const WALK_MESSAGE = "WALK_MESSAGE"

type WalkComponent struct {
	engo.Point
}

type walkEntity struct {
	*ecs.BasicEntity
	*WalkComponent
	*common.SpaceComponent
}

type WalkSystem struct {
	ids      map[uint64]struct{}
	entities []walkEntity
	land     *core.WorldSystem
}

func (s *WalkSystem) New(w *ecs.World) {
	s.ids = make(map[uint64]struct{})

	for _, system := range w.Systems() {
		switch sys := system.(type) {
		case *core.WorldSystem:
			s.land = sys
		}
	}

	s.listen()
}

func (s *WalkSystem) listen() {
	engo.Mailbox.Listen(ACTION_MESSAGE, func(message engo.Message) {
		msg, isAction := message.(ActionMessage)
		if isAction {
			for _, e := range s.entities {
				if e.ID() == msg.BasicEntity.ID() {
					e.WalkComponent.Point = msg.State.Speed
				}
			}
		}
	})
}

func (s *WalkSystem) inBound() {

}

func (s *WalkSystem) Add(basic *ecs.BasicEntity, speed *WalkComponent, space *common.SpaceComponent) {
	if _, ok := s.ids[basic.ID()]; ok {
		return
	}
	s.ids[basic.ID()] = struct{}{}
	s.entities = append(s.entities, walkEntity{basic, speed, space})
}

func (s *WalkSystem) Remove(basic ecs.BasicEntity) {
	delete := -1
	for index, e := range s.entities {
		if e.BasicEntity.ID() == basic.ID() {
			delete = index
			break
		}
	}
	if delete >= 0 {
		s.entities = append(s.entities[:delete], s.entities[delete+1:]...)
	}
}

func (s *WalkSystem) Update(dt float32) {
	for _, e := range s.entities {
		speed := engo.GameWidth() * dt
		nextPosition := engo.Point{X: e.SpaceComponent.Position.X + speed*e.WalkComponent.Point.X, Y: e.SpaceComponent.Position.Y + speed*e.WalkComponent.Point.Y}

		if nextPosition.Within(s.land) {
			e.SpaceComponent.Position = nextPosition
		}
	}
}
