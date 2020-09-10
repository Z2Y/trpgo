package core

import (
	"github.com/EngoEngine/ecs"
	"github.com/EngoEngine/engo"
	"github.com/EngoEngine/engo/common"
)

const gridSize = 64

type Grid struct {
	ecs.BasicEntity
	*common.RenderComponent
	*common.SpaceComponent

	SubEntites []*Entity

	Code int
}

type Entity struct {
	ecs.BasicEntity
	*common.RenderComponent
	*common.SpaceComponent

	X, Y          int
	Width, Height int
}

func AlignEntityToGrid(e *Entity) *Entity {
	if e == nil {
		return nil
	}

	drawable := e.RenderComponent.Drawable

	if e.Width == 0 || e.Height == 0 {
		gridCenter := getGridCenter(e.X, e.Y, 1, 1)
		entityFooter := getEntityFooter(e.RenderComponent)
		e.SpaceComponent.Position = engo.Point{X: gridCenter.X - entityFooter.X, Y: gridCenter.Y - entityFooter.Y}
	} else {
		scale := float32((1+e.Width)*gridSize/2) / drawable.Width()
		e.RenderComponent.Scale = engo.Point{X: scale, Y: scale}
		entityFooter := getEntityFooter(e.RenderComponent)
		gridCenter := getGridCenter(e.X, e.Y, e.Width, e.Height)
		e.SpaceComponent.Position = engo.Point{X: gridCenter.X - entityFooter.X, Y: gridCenter.Y - entityFooter.Y + float32((1+e.Height)*gridSize/8)}
	}
	return e
}

func getEntityFooter(render *common.RenderComponent) engo.Point {
	width := render.Drawable.Width() * render.Scale.X
	height := render.Drawable.Height() * render.Scale.Y
	return engo.Point{X: width / 2, Y: height}
}

func getGridCenter(x, y, width, height int) engo.Point {
	return engo.Point{X: float32(x*(gridSize/2) + y*gridSize/2 + (1+width)*gridSize/4), Y: float32(y*gridSize/4 - x*gridSize/4 + (1+height)*gridSize/8)}
}
