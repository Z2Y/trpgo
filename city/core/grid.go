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

	Offset        engo.Point
	X, Y          int
	Width, Height int
}

func (g *Grid) Blocked() bool {
	if g.Code == Waters[0] {
		return true
	}
	for _, e := range g.SubEntites {
		if e.Width > 0 && e.Height > 0 {
			return true
		}
	}
	return false
}

func (e *Entity) GetZIndex() float32 {
	return e.RenderComponent.StartZIndex
}

func AlignEntityToGrid(e *Entity) *Entity {
	if e == nil {
		return nil
	}

	drawable := e.RenderComponent.Drawable

	if e.Width == 0 || e.Height == 0 {
		gridCenter := GridCenter(e.X, e.Y, 1, 1)
		e.Offset = getEntityFooter(e.RenderComponent)
		e.SpaceComponent.Position = engo.Point{X: gridCenter.X - e.Offset.X, Y: gridCenter.Y - e.Offset.Y}
		e.RenderComponent.StartZIndex = gridCenter.Y
	} else {
		scale := float32((1+e.Width)*gridSize/2) / drawable.Width()
		e.RenderComponent.Scale = engo.Point{X: scale, Y: scale}
		entityFooter := getEntityFooter(e.RenderComponent)
		gridCenter := GridCenter(e.X, e.Y, e.Width, e.Height)
		e.Offset = engo.Point{X: entityFooter.X, Y: entityFooter.Y - float32((1+e.Height)*gridSize/8)}
		e.SpaceComponent.Position = engo.Point{X: gridCenter.X - e.Offset.X, Y: gridCenter.Y - e.Offset.Y}
		e.RenderComponent.StartZIndex = gridCenter.Y
	}
	return e
}

func getEntityFooter(render *common.RenderComponent) engo.Point {
	width := render.Drawable.Width() * render.Scale.X
	height := render.Drawable.Height() * render.Scale.Y
	return engo.Point{X: width / 2, Y: height}
}

func GridCenter(x, y, width, height int) engo.Point {
	return engo.Point{X: float32(x*(gridSize/2) + y*gridSize/2 + (1+width)*gridSize/4), Y: float32(y*gridSize/4 - x*gridSize/4 + (1+height)*gridSize/8)}
}

func GridIndex(x, y float32) (int, int) {
	py := int((x+2*y)/gridSize - 0.5)
	px := int((x-2*y)/gridSize + 0.5)
	return px, py
}
