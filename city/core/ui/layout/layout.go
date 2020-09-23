package layout

import (
	"github.com/EngoEngine/engo"
	"github.com/Z2Y/trpgo/city/config"
)

func AlignCenter(parent, child engo.AABB) engo.Point {
	pWidth := (parent.Max.X - parent.Min.X)
	pHeight := (parent.Max.Y - parent.Min.Y)

	cWidth := (child.Max.X - child.Min.X)
	cHeight := (child.Max.Y - child.Min.Y)

	return engo.Point{X: parent.Min.X + (pWidth-cWidth)/2, Y: parent.Min.Y + (pHeight-cHeight)/2}
}

func AlignRightBottom(parent engo.AABB, child engo.AABB, right, bottom float32) engo.Point {
	cWidth := (child.Max.X - child.Min.X)
	cHeight := (child.Max.Y - child.Min.Y)
	return engo.Point{X: parent.Max.X - right - cWidth, Y: parent.Max.Y - bottom - cHeight}
}

func AlignRightTop(parent engo.AABB, child engo.AABB, right, top float32) engo.Point {
	cWidth := (child.Max.X - child.Min.X)
	return engo.Point{X: parent.Max.X - right - cWidth, Y: parent.Min.Y + top}
}

func AlignToWorldCenter(child engo.AABB) engo.Point {
	p := engo.AABB{Max: engo.Point{X: float32(config.GameWidth), Y: float32(config.GameHeight)}}
	return AlignCenter(p, child)
}

func AlignToWorldRightBottom(child engo.AABB, right, bottom float32) engo.Point {
	p := engo.AABB{Max: engo.Point{X: float32(config.GameWidth), Y: float32(config.GameHeight)}}
	return AlignRightBottom(p, child, right, bottom)
}
