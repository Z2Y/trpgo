package ui

import (
	"github.com/EngoEngine/ecs"
	"github.com/EngoEngine/engo"
	"github.com/EngoEngine/engo/common"
)

type List struct {
	UIBasic

	Background *common.Texture

	Position engo.Point
	Width    float32
	Height   float32

	Items []*UIBasic
}

func NewList(l List) *List {
	l.BasicEntity = ecs.NewBasic()

	l.SpaceComponent = common.SpaceComponent{
		Position: l.Position,
		Width:    l.Width,
		Height:   l.Height,
	}

	l.RenderComponent = common.RenderComponent{
		Drawable: l.Background,
		Scale:    engo.Point{X: l.Width / l.Background.Width(), Y: l.Height / l.Background.Height()},
	}

	return &l
}

