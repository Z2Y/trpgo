package human

import (
	"fmt"

	"github.com/EngoEngine/engo"
	"github.com/EngoEngine/engo/common"

	"github.com/Z2Y/trpgo/city/core/action"
)

var (
	StopAction       *common.Animation
	DefaultAnimation *common.AnimationComponent
	Animations       map[string]*common.AnimationComponent
)

type Human struct {
	action.ActionEntity
	common.AnimationComponent
	common.RenderComponent

	Offset engo.Point
}

func frames(size int) []int {
	f := make([]int, size)
	for i := 0; i < size; i++ {
		f[i] = i
	}
	return f
}

func initAnimation(name string) *common.AnimationComponent {
	sprite := common.NewSpritesheetFromFile(fmt.Sprintf("npc/%s.png", name), 260, 210)

	anim := common.NewAnimationComponent(sprite.Drawables(), 0.1)
	anim.AddDefaultAnimation(&common.Animation{Name: name, Frames: frames(sprite.CellCount())})
	return &anim
}

func Init() {
	Animations = make(map[string]*common.AnimationComponent)
	Animations["idle"] = initAnimation("idle")
	Animations["walking"] = initAnimation("walking")
	DefaultAnimation = Animations["idle"]
}

func NewHuman(point engo.Point) *Human {
	entity := &Human{ActionEntity: action.NewActionEntity()}

	entity.RenderComponent = common.RenderComponent{
		Drawable:    DefaultAnimation.Drawables[0],
		Scale:       engo.Point{X: 0.5, Y: 0.5},
		StartZIndex: 1,
	}

	entity.SpaceComponent = common.SpaceComponent{
		Position: entity.getFixedPosition(point),
		Width:    64,
		Height:   64,
	}

	entity.AnimationComponent = *DefaultAnimation

	entity.listenAction()

	return entity
}

func (h *Human) getFixedPosition(point engo.Point) engo.Point {
	center := getCenterOfRender(h.RenderComponent)
	h.Offset = center
	return engo.Point{X: point.X - center.X, Y: point.Y - center.Y}
}

func getCenterOfRender(render common.RenderComponent) engo.Point {
	width := render.Drawable.Width() * render.Scale.X
	height := render.Drawable.Height() * render.Scale.Y
	return engo.Point{X: width / 2, Y: height / 2}
}

func (h *Human) listenAction() {
	engo.Mailbox.Listen("ActionMessage", func(message engo.Message) {
		msg, isAction := message.(action.ActionMessage)
		if isAction && h.ID() == msg.BasicEntity.ID() {
			switch msg.State.Code {
			case action.ActionIdle:
				h.AnimationComponent = *DefaultAnimation
			case action.ActionWalking:
				h.UpdateFace()
				h.AnimationComponent = *Animations["walking"]
			}
		}
	})
}

func (h *Human) UpdateFace() {
	if h.ActionState.Speed.X < 0 {
		h.RenderComponent.Scale.X = -0.5
	} else {
		h.RenderComponent.Scale.X = 0.5
	}
	offset := h.Offset

	originPosition := engo.Point{X: h.SpaceComponent.Position.X + offset.X, Y: h.SpaceComponent.Position.Y + offset.Y}
	h.SpaceComponent.Position = h.getFixedPosition(originPosition)
}
