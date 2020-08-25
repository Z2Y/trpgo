package core

import (
	"math/rand"
	"time"

	"github.com/EngoEngine/ecs"
	"github.com/EngoEngine/engo"
	"github.com/EngoEngine/engo/common"
)

type WorldSystem struct {
	world    *ecs.World
	renderer *common.RenderSystem
	grids    []*Grid

	width  float32
	height float32
}

func (w *WorldSystem) New(world *ecs.World) {
	w.world = world

	for _, system := range world.Systems() {
		switch sys := system.(type) {
		case *common.RenderSystem:
			w.renderer = sys
		}
	}
}

func (w *WorldSystem) LoadWorldMap(worldMap *WorldMap) {
	w.width = float32(worldMap.xLen * gridSize)
	w.height = float32(worldMap.yLen*gridSize) * 0.5

	rand.Seed(time.Now().UnixNano())

	for x := 0; x < worldMap.xLen; x++ {
		for y := 0; y < worldMap.yLen; y++ {
			space := &common.SpaceComponent{
				Position: engo.Point{X: float32(x*(gridSize/2) + y*gridSize/2), Y: float32(y*gridSize/4 - x*gridSize/4)},
				Width:    gridSize,
				Height:   gridSize,
			}
			code := rand.Intn(len(Lands))
			w.grids = append(w.grids, &Grid{BasicEntity: ecs.NewBasic(), RenderComponent: Entitys[code], SpaceComponent: space})
		}
	}

	w.generateTrees(worldMap.xLen, worldMap.yLen)
}

func (w *WorldSystem) generateTrees(width, height int) {
	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			code := rand.Intn(40)
			if code < len(Trees) {
				space := &common.SpaceComponent{
					Position: engo.Point{X: float32(x*(gridSize/2) + y*gridSize/2), Y: float32(y*gridSize/4 - x*gridSize/4)},
					Width:    gridSize,
					Height:   gridSize,
				}
				render := &common.RenderComponent{
					Drawable:    Entitys[Trees[code]].Drawable,
					Scale:       Entitys[Trees[code]].Scale,
					StartZIndex: 2,
				}
				w.grids = append(w.grids, &Grid{BasicEntity: ecs.NewBasic(), RenderComponent: render, SpaceComponent: space})
			}
		}
	}
}

func (w *WorldSystem) Size() (float32, float32) {
	return w.width, w.height
}

func (w *WorldSystem) Center() engo.Point {
	return engo.Point{X: w.width / 2, Y: 0}
}

func (w *WorldSystem) Bounds() engo.AABB {
	return engo.AABB{
		Min: engo.Point{X: 0, Y: -w.height / 2},
		Max: engo.Point{X: w.width, Y: w.height / 2},
	}
}

func (w *WorldSystem) Contains(point engo.Point) bool {
	bounds := w.Bounds()
	return point.X > bounds.Min.X && point.X < bounds.Max.X && point.Y > bounds.Min.Y && point.Y < bounds.Max.Y
}

func (w *WorldSystem) Update(float32) {
	for _, grid := range w.grids {
		w.renderer.Add(&grid.BasicEntity, grid.RenderComponent, grid.SpaceComponent)
	}
}

func (w *WorldSystem) Remove(basic ecs.BasicEntity) {

}
