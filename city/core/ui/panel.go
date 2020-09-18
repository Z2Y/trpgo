package ui

import (
	"github.com/EngoEngine/ecs"
	"github.com/EngoEngine/engo"
	"github.com/EngoEngine/engo/common"
)

type Panel struct {
	UIBasic

	Image *common.Texture
}

func NewPanel(panel Panel) *Panel {
	panel.BasicEntity = ecs.NewBasic()

	panel.SpaceComponent = common.SpaceComponent{
		Width:  panel.Width,
		Height: panel.Height,
	}

	panel.RenderComponent = common.RenderComponent{
		Drawable:    panel.Image,
		Scale:       engo.Point{X: panel.Width / panel.Image.Width(), Y: panel.Height / panel.Image.Height()},
		StartZIndex: UILayerIndex,
	}

	panel.SetShader(common.HUDShader)

	return &panel
}
