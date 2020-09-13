package control

import (
	"math/rand"

	"github.com/EngoEngine/ecs"
	"github.com/EngoEngine/engo"
	"github.com/EngoEngine/engo/common"
	"github.com/Z2Y/trpgo/city/core"
	"github.com/Z2Y/trpgo/city/core/route"
)

const (
	ActionIdle = iota
	ActionWalking
)

const ACTION_MESSAGE = "ActionMessage"

type ActionEntity struct {
	ecs.BasicEntity
	common.SpaceComponent
	common.RenderComponent

	Offset engo.Point
	ActionState
}

func (e *ActionEntity) FootPosition() engo.Point {
	footX, footY := e.SpaceComponent.Position.X+e.Offset.X, e.SpaceComponent.Position.Y+e.Offset.Y+e.SpaceComponent.Height/2
	return engo.Point{X: footX, Y: footY}
}

type ActionState struct {
	Code  int
	Speed engo.Point

	Route     *route.RoutePath
	NextRoute *route.RoutePath

	duration float32
	elapsed  float32
}

type ActionMessage struct {
	*ecs.BasicEntity
	State ActionState
}

func (ActionMessage) Type() string {
	return ACTION_MESSAGE
}

type AutoActionSystem struct {
	entities []*ActionEntity

	walkSys *WalkSystem
}

type WalkAble interface {
	Walk(engo.Point)
}

func NewActionEntity() ActionEntity {
	return ActionEntity{BasicEntity: ecs.NewBasic()}
}

func (s *AutoActionSystem) New(w *ecs.World) {
	for _, system := range w.Systems() {
		switch sys := system.(type) {
		case *WalkSystem:
			s.walkSys = sys
		}
	}

	if s.walkSys == nil {
		s.walkSys = &WalkSystem{}
		w.AddSystem(s.walkSys)
	}
}

func (s *AutoActionSystem) Add(entity *ActionEntity) {
	s.entities = append(s.entities, entity)
	s.walkSys.Add(entity)
}

func (s *AutoActionSystem) Remove(basic ecs.BasicEntity) {
	delete := -1
	for index, e := range s.entities {
		if e.ID() == basic.ID() {
			delete = index
			break
		}
	}
	if delete >= 0 {
		s.entities = append(s.entities[:delete], s.entities[delete+1:]...)
		s.walkSys.Remove(basic)
	}
}

func (s *AutoActionSystem) Update(dt float32) {
	for _, e := range s.entities {
		preState := e.ActionState.Code
		e.ActionState.Update(e, dt)
		switch e.ActionState.Code {
		case ActionIdle:
			if preState != ActionIdle {
				engo.Mailbox.Dispatch(ActionMessage{BasicEntity: &e.BasicEntity, State: e.ActionState})
			}
		case ActionWalking:
			if e.ActionState.elapsed == 0 {
				engo.Mailbox.Dispatch(ActionMessage{BasicEntity: &e.BasicEntity, State: e.ActionState})
			}
		}
	}
}

func (s *ActionState) Update(e *ActionEntity, dt float32) {
	s.elapsed = s.elapsed + dt

	if s.elapsed >= s.duration {
		s.Next(e)
	}
}

func (s *ActionState) Finished() bool {
	return s.elapsed >= s.duration
}

func (s *ActionState) Next(e *ActionEntity) {
	if s.NextRoute != nil {
		s.Route = s.NextRoute
		s.NextRoute = s.Route.Parent
		s.routeState(e.FootPosition(), s.Route)
	} else {
		s.Route = nil
		s.controlState()
	}
}

func (s *ActionState) controlState() {
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

	speed, _ = speed.Normalize()

	if speed.X != s.Speed.X || speed.Y != s.Speed.Y {
		if speed.X == 0 && speed.Y == 0 {
			s.Code = ActionIdle
		} else {
			s.Code = ActionWalking
		}
		s.elapsed = 0
		s.duration = 0.01
		s.Speed = speed
	}
}

func (s *ActionState) randomState() {
	s.Code = rand.Intn(2)

	s.Speed.X, s.Speed.Y = 0, 0
	if s.Code == ActionWalking {
		s.Speed.X = 0.1 * float32((rand.Intn(3) - 1))
		s.Speed.Y = 0.1 * float32((rand.Intn(3) - 1))
	}

	if s.Speed.X == 0 && s.Speed.Y == 0 {
		s.Code = ActionIdle
	}

	s.duration = rand.Float32() * 4
	s.elapsed = 0
}

func (s *ActionState) routeState(source engo.Point, path *route.RoutePath) {
	// sx, sy := core.GridIndex(source.X, source.Y)
	tx, ty := int(path.Point.X), int(path.Point.Y)
	target := core.GridCenter(tx, ty, 1, 1)
	path.Point = target
	if source.Equal(target) {
		s.Speed.X, s.Speed.Y = 0, 0
		s.duration = 0
	} else {
		speed := engo.Point{X: target.X - source.X, Y: target.Y - source.Y}
		s.Speed, _ = speed.Normalize()

		s.Code = ActionWalking
		s.duration = source.PointDistance(target) / WALK_SPEED_SCALE
	}
	s.elapsed = 0
}
