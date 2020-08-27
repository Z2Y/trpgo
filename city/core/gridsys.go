package core

import (
	"log"
	"math/rand"
	"time"

	"github.com/EngoEngine/ecs"
	"github.com/EngoEngine/engo"
	"github.com/EngoEngine/engo/common"
	"github.com/Z2Y/trpgo/city/core/mapgenerator"
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

	start := time.Now()
	height_map := mapgenerator.Generate(worldMap.xLen, worldMap.yLen)
	endAt := time.Now()
	log.Println("Map Gen End", endAt.Sub(start))

	for x := 0; x < worldMap.xLen; x++ {
		for y := 0; y < worldMap.yLen; y++ {
			space := &common.SpaceComponent{
				Position: engo.Point{X: float32(x*(gridSize/2) + y*gridSize/2), Y: float32(y*gridSize/4 - x*gridSize/4)},
				Width:    gridSize,
				Height:   gridSize,
			}
			code := rand.Intn(len(Lands))
			if height_map.Value(x, y) == 0 {
				code = Waters[0]
			} else if height_map.Value(x, y) >= 1 && height_map.Value(x, y) <= 2 {
				code = Sand[0]
			}
			w.grids = append(w.grids, &Grid{BasicEntity: ecs.NewBasic(), RenderComponent: Entitys[code], SpaceComponent: space})
		}
	}

	w.generateTrees(height_map, worldMap.xLen, worldMap.yLen)
}

func (w *WorldSystem) generateTrees(height_map *mapgenerator.HeightMap, width, height int) {
	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			if height_map.Value(x, y) >= 30 && height_map.Value(x, y) <= 64 {
				code := rand.Intn(10)
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
