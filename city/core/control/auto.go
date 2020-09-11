package control

import (
	"math/rand"

	"github.com/EngoEngine/ecs"
	"github.com/EngoEngine/engo"
	"github.com/EngoEngine/engo/common"
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
	WalkComponent

	Offset engo.Point
	ActionState
}

type ActionState struct {
	Code  int
	Speed engo.Point

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
		e.ActionState.Update(dt)
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

func (s *ActionState) Update(dt float32) {
	s.elapsed = s.elapsed + dt

	if s.elapsed >= s.duration {
		s.Next()
	}
}

func (s *ActionState) Finished() bool {
	return s.elapsed >= s.duration
}

func (s *ActionState) Next() {
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
