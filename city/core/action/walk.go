package action

import (
	"github.com/EngoEngine/ecs"
	"github.com/EngoEngine/engo"
	"github.com/EngoEngine/engo/common"
)

const WALK_MESSAGE = "WALK_MESSAGE"

type WalkComponent struct {
	engo.Point
}

type WalkMessage struct {
	*ecs.BasicEntity
	engo.Point
}

func (WalkMessage) Type() string {
	return WALK_MESSAGE
}

type walkEntity struct {
	*ecs.BasicEntity
	*WalkComponent
	*common.SpaceComponent
}

type WalkSystem struct {
	ids      map[uint64]struct{}
	entities []walkEntity
}

func (s *WalkSystem) New(*ecs.World) {
	s.ids = make(map[uint64]struct{})
	engo.Mailbox.Listen(WALK_MESSAGE, func(message engo.Message) {
		speed, isSpeed := message.(WalkMessage)
		if isSpeed {
			for _, e := range s.entities {
				if e.ID() == speed.BasicEntity.ID() {
					e.WalkComponent.Point = speed.Point
				}
			}
		}
	})
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
		e.SpaceComponent.Position.X = e.SpaceComponent.Position.X + speed*e.WalkComponent.Point.X
		e.SpaceComponent.Position.Y = e.SpaceComponent.Position.Y + speed*e.WalkComponent.Point.Y
	}
}
