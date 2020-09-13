package control

import (
	"github.com/EngoEngine/ecs"
	"github.com/EngoEngine/engo"
	"github.com/Z2Y/trpgo/city/core"
)

var WALK_SPEED_SCALE = float32(96)

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
		if e.ActionState.Speed.X == 0 && e.ActionState.Speed.Y == 0 {
			continue
		}

		speedScale := dt * WALK_SPEED_SCALE
		speed := e.ActionState.Speed
		speed.MultiplyScalar(speedScale)
		nextPosition := engo.Point{X: e.SpaceComponent.Position.X + speed.X, Y: e.SpaceComponent.Position.Y + speed.Y}

		if e.ActionState.Route != nil {
			dest := engo.Point{X: e.ActionState.Route.Point.X - e.Offset.X, Y: e.ActionState.Route.Point.Y - e.Offset.Y - e.SpaceComponent.Height/2}
			if nextPosition.PointDistance(dest) < speedScale {
				e.ActionState.elapsed = e.ActionState.duration
				nextPosition = dest
			}
		}

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
