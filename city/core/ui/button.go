package ui

import (
	"github.com/EngoEngine/ecs"
	"github.com/EngoEngine/engo"
	"github.com/EngoEngine/engo/common"
)

type Button struct {
	UIBasic

	Text         *Text
	Image        *common.Texture
	ImagePressed *common.Texture

	Position  engo.Point
	Width     float32
	Height    float32
	TextAlign int
}

func (b *Button) OnClick(f func()) {
	if b.MessageListener == nil {
		b.MessageListener = &engo.MessageManager{}
	}
	b.MessageListener.Listen("UIMouseEvent", func(engo.Message) {
		f()
	})
}

func NewButton(b Button) *Button {

	b.BasicEntity = ecs.NewBasic()

	b.SpaceComponent = common.SpaceComponent{
		Position: b.Position,
		Width:    b.Width,
		Height:   b.Height,
	}

	b.RenderComponent = common.RenderComponent{
		Drawable: b.Image,
		Scale:    engo.Point{X: b.Width / b.Image.Width(), Y: b.Height / b.Image.Height()},
	}

	b.SetShader(common.HUDShader)
	b.SetZIndex(UILayerIndex)

	if b.Text != nil {
		b.alignText()
	}

	return &b
}

func (b *Button) SetPosition(pos engo.Point) {
	b.Position = pos
	b.SpaceComponent.Position = pos
	if b.Text != nil {
		b.alignText()
	}
}

func (b *Button) alignText() {
	offsetX := (b.SpaceComponent.Width - b.Text.SpaceComponent.Width) / 2
	offsetY := (b.SpaceComponent.Height - b.Text.SpaceComponent.Height) / 2

	b.Text.SpaceComponent.Position = engo.Point{X: b.SpaceComponent.Position.X + offsetX, Y: b.SpaceComponent.Position.Y + offsetY}
}
