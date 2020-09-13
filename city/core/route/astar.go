package route

import (
	"container/heap"
	"math"

	"github.com/EngoEngine/engo"
	"github.com/Z2Y/trpgo/city/core"
)

type AstarRoute struct {
	World *core.WorldSystem

	closed map[engo.Point]*RoutePath
	opened map[engo.Point]*RoutePath
}

type RoutePath struct {
	Point  engo.Point
	Parent *RoutePath

	gWeight int
	hWeight int
	fWeight int
}

type openList []*RoutePath

func (o openList) Len() int {
	return len(o)
}

func (o openList) Less(i, j int) bool {
	return o[i].fWeight < o[j].fWeight
}

func (o openList) Swap(i, j int) {
	o[i], o[j] = o[j], o[i]
}

func (o *openList) Push(p interface{}) {
	*o = append(*o, p.(*RoutePath))
}

func (o *openList) Pop() interface{} {
	size := len(*o)
	x := (*o)[size-1]
	*o = (*o)[0 : size-1]
	return x
}

func (r *RoutePath) updateWeight(target engo.Point) {
	r.hWeight = int(math.Abs(float64(r.Point.X-target.X))+math.Abs(float64(r.Point.Y-target.Y))) * 10

	if r.Parent != nil {
		deltaX := math.Abs(float64(r.Parent.Point.X-r.Point.X)) * 10
		deltaY := math.Abs(float64(r.Parent.Point.Y-r.Point.Y)) * 10
		r.gWeight = r.Parent.gWeight + int(math.Sqrt(deltaX*deltaX+deltaY*deltaY))
	}

	r.fWeight = r.gWeight + r.hWeight
}

func (r *RoutePath) Reserve() *RoutePath {
	if r == nil || r.Parent == nil {
		return r
	}
	p := r.Parent.Reserve()
	r.Parent.Parent = r
	r.Parent = nil
	return p
}

func newOpenList(source, target engo.Point) *openList {
	l := openList{}
	sourcePath := &RoutePath{
		Point:   source,
		Parent:  nil,
		gWeight: 0,
	}
	sourcePath.updateWeight(target)
	heap.Init(&l)
	heap.Push(&l, sourcePath)
	return &l
}

func (a *AstarRoute) FindPath(source, target engo.Point) *RoutePath {
	openL := newOpenList(source, target)
	a.opened = make(map[engo.Point]*RoutePath)
	a.closed = make(map[engo.Point]*RoutePath)

	var current *RoutePath

	for len(*openL) > 0 {
		current = heap.Pop(openL).(*RoutePath)

		if current.Point == target {
			return current
		}

		delete(a.opened, current.Point)
		a.closed[current.Point] = current

		for _, p := range a.adjacentPoints(current.Point) {
			if a.isClosed(p) {
				continue
			}

			newPath := &RoutePath{
				Point:  p,
				Parent: current,
			}

			newPath.updateWeight(target)

			oldPath, ok := a.opened[p]
			if !ok {
				a.opened[p] = newPath
				heap.Push(openL, newPath)
			} else {
				if newPath.fWeight < oldPath.fWeight {
					oldPath.Parent = current
					oldPath.updateWeight(target)
				}
			}
		}
	}
	return nil
}

func (a *AstarRoute) isClosed(p engo.Point) bool {
	px, py := int(p.X), int(p.Y)
	_, closed := a.closed[p]
	if closed || !a.World.InView(px, py) {
		return true
	}
	grid := a.World.GetGrid(px, py)
	return grid == nil || grid.Blocked()
}

func (a *AstarRoute) adjacentPoints(p engo.Point) []engo.Point {
	var adjacents []engo.Point

	sx, sy := int(p.X), int(p.Y)

	for x := sx - 1; x <= sx+1; x++ {
		for y := sy - 1; y <= sy+1; y++ {
			adjacents = append(adjacents, engo.Point{X: float32(x), Y: float32(y)})
		}
	}
	return adjacents
}
