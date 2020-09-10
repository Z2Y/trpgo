package human

import (
	"fmt"

	"github.com/EngoEngine/engo"
	"github.com/EngoEngine/engo/common"

	"github.com/Z2Y/trpgo/city/core/control"
)

var (
	StopAction       *common.Animation
	DefaultAnimation *common.AnimationComponent
	Animations       map[string]*common.AnimationComponent
)

type Human struct {
	control.ActionEntity
	common.AnimationComponent
	common.RenderComponent

	CurrentAnimation *common.AnimationComponent
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
	entity := &Human{ActionEntity: control.NewActionEntity()}

	entity.RenderComponent = common.RenderComponent{
		Drawable:    DefaultAnimation.Drawables[0],
		Scale:       engo.Point{X: 0.25, Y: 0.25},
		StartZIndex: 1,
	}

	entity.SpaceComponent = common.SpaceComponent{
		Position: entity.getFixedPosition(point),
		Width:    30,
		Height:   30,
	}

	entity.AnimationComponent = *DefaultAnimation
	entity.SetAnimation(DefaultAnimation)

	entity.listenAction()

	return entity
}

func (h *Human) getFixedPosition(point engo.Point) engo.Point {
	center := getFootOffset(h.RenderComponent)
	h.Offset = center
	return engo.Point{X: point.X - center.X, Y: point.Y - center.Y}
}

func getFootOffset(render common.RenderComponent) engo.Point {
	width := render.Drawable.Width() * render.Scale.X
	height := render.Drawable.Height() * render.Scale.Y
	return engo.Point{X: width / 2, Y: height / 2}
}

func (h *Human) listenAction() {
	engo.Mailbox.Listen("ActionMessage", func(message engo.Message) {
		msg, isAction := message.(control.ActionMessage)
		if isAction && h.ID() == msg.BasicEntity.ID() {
			switch msg.State.Code {
			case control.ActionIdle:
				h.SetAnimation(DefaultAnimation)
			case control.ActionWalking:
				h.UpdateFace()
				h.SetAnimation(Animations["walking"])
			}
		}
	})
}

func (h *Human) SetAnimation(animation *common.AnimationComponent) {
	if h.CurrentAnimation != animation {
		h.AnimationComponent = *animation
		h.CurrentAnimation = animation
	}
}

func (h *Human) UpdateFace() {
	if h.ActionState.Speed.X == 0 {
		return
	}
	if h.ActionState.Speed.X < 0 {
		h.RenderComponent.Scale.X = -0.25
	} else {
		h.RenderComponent.Scale.X = 0.25
	}
	offset := h.Offset

	originPosition := engo.Point{X: h.SpaceComponent.Position.X + offset.X, Y: h.SpaceComponent.Position.Y + offset.Y}
	h.SpaceComponent.Position = h.getFixedPosition(originPosition)
}
