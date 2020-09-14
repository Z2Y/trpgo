package ui

import (
	"github.com/EngoEngine/ecs"
	"github.com/EngoEngine/engo"
	"github.com/EngoEngine/engo/common"
)

type Image struct {
	UIBasic

	Texture  *common.Texture
	Position engo.Point
	Width    float32
	Height   float32
}

func NewImage(img Image) *Image {

	img.BasicEntity = ecs.NewBasic()

	img.SpaceComponent = common.SpaceComponent{
		Position: img.Position,
		Width:    img.Width,
		Height:   img.Height,
	}

	img.RenderComponent = common.RenderComponent{
		Drawable:    img.Texture,
		Scale:       engo.Point{X: img.Width / img.Texture.Width(), Y: img.Height / img.Texture.Height()},
		StartZIndex: UILayerIndex,
	}

	img.SetShader(common.HUDShader)

	return &img
}

func (img *Image) SetPosition(pos engo.Point) {
	img.Position = pos
	img.SpaceComponent.Position = pos
}
