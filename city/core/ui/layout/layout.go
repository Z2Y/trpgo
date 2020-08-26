package layout

import (
	"github.com/EngoEngine/engo"
)

func AlignCenter(parent, child engo.AABB) engo.Point {
	pWidth := (parent.Max.X - parent.Min.X)
	pHeight := (parent.Max.Y - parent.Min.Y)

	cWidth := (child.Max.X - child.Min.X)
	cHeight := (child.Max.Y - child.Min.Y)

	return engo.Point{X: parent.Min.X + (pWidth-cWidth)/2, Y: parent.Min.Y + (pHeight-cHeight)/2}
}

func AlignToWorldCenter(child engo.AABB) engo.Point {
	p := engo.AABB{Max: engo.Point{X: engo.WindowWidth() / engo.GetGlobalScale().X, Y: engo.WindowHeight() / engo.GetGlobalScale().Y}}
	return AlignCenter(p, child)
}
