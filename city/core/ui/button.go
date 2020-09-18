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

	TextAlign int
}

func NewButton(b Button) *Button {

	b.BasicEntity = ecs.NewBasic()

	b.SpaceComponent = common.SpaceComponent{
		Width:  b.Width,
		Height: b.Height,
	}

	b.RenderComponent = common.RenderComponent{
		Drawable:    b.Image,
		Scale:       engo.Point{X: b.Width / b.Image.Width(), Y: b.Height / b.Image.Height()},
		StartZIndex: UILayerIndex,
	}

	b.SetShader(common.HUDShader)
	b.EventResponder = true

	if b.Text != nil {
		b.alignText()
		b.AddSubEntity(b.Text)
	}

	return &b
}

func (b *Button) alignText() {
	offsetX := (b.SpaceComponent.Width - b.Text.SpaceComponent.Width) / 2
	offsetY := (b.SpaceComponent.Height - b.Text.SpaceComponent.Height) / 2

	b.Text.SpaceComponent.Position = engo.Point{X: b.SpaceComponent.Position.X + offsetX, Y: b.SpaceComponent.Position.Y + offsetY}
}
