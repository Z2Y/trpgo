package control

import (
	"github.com/EngoEngine/ecs"
	"github.com/EngoEngine/engo"
	"github.com/Z2Y/trpgo/city/core"
)

const WALK_MESSAGE = "WALK_MESSAGE"

type WalkComponent struct {
	engo.Point
}

type WalkSystem struct {
	ids      map[uint64]struct{}
	entities []*ActionEntity
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

func (s *WalkSystem) Add(entity *ActionEntity) {
	if _, ok := s.ids[entity.ID()]; ok {
		return
	}
	s.ids[entity.ID()] = struct{}{}
	s.entities = append(s.entities, entity)
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
		if e.WalkComponent.Point.X == 0 && e.WalkComponent.Point.Y == 0 {
			continue
		}
		nextPosition := engo.Point{X: e.SpaceComponent.Position.X + e.WalkComponent.Point.X, Y: e.SpaceComponent.Position.Y + e.WalkComponent.Point.Y}
		footX, footY := nextPosition.X+e.Offset.X, nextPosition.Y+e.Offset.Y+e.SpaceComponent.Height/2
		gx, gy := s.land.GetGridPos(footX, footY)

		grid := s.land.GetGrid(gx, gy)

		if grid == nil || grid.Blocked() {
			continue
		}

		ox, oy := s.land.GetGridPos(e.SpaceComponent.Position.X+e.Offset.X, e.SpaceComponent.Position.Y+e.Offset.Y)
		e.SpaceComponent.Position = nextPosition

		if ox != gx || oy != gy {
			e.RenderComponent.SetZIndex(footY)
			continue
		}

		for x := gx - 1; x <= gx+1; x++ {
			for y := gy - 1; y < gy+1; y++ {
				grid := s.land.GetGrid(x, y)
				if grid != nil && len(grid.SubEntites) > 0 {
					e.RenderComponent.SetZIndex(footY)
				}
			}
		}
	}
}
