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
	camera   *common.CameraSystem

	ground map[int]map[int]*Grid // 地面

	width  float32
	height float32

	cameraX int
	cameraY int
	cameraZ int

	viewLenX int
	viewLenY int
}

func (w *WorldSystem) New(world *ecs.World) {
	w.world = world

	for _, system := range world.Systems() {
		switch sys := system.(type) {
		case *common.RenderSystem:
			w.renderer = sys
		case *common.CameraSystem:
			w.camera = sys
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

	w.ground = make(map[int]map[int]*Grid)
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
			if w.ground[x] == nil {
				w.ground[x] = make(map[int]*Grid)
			}
			w.ground[x][y] = &Grid{BasicEntity: ecs.NewBasic(), RenderComponent: Entitys[code], SpaceComponent: space, Code: code}
		}
	}

	w.generateTrees(height_map, worldMap.xLen, worldMap.yLen)
	w.generateBuildings(height_map, worldMap.xLen, worldMap.yLen)
}

func (w *WorldSystem) generateTrees(height_map *mapgenerator.HeightMap, width, height int) {
	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			if height_map.Value(x, y) >= 30 && height_map.Value(x, y) <= 64 {
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
					ground := w.ground[x][y]
					if ground != nil {
						ground.SubEntites = append(ground.SubEntites, &Grid{BasicEntity: ecs.NewBasic(), RenderComponent: render, SpaceComponent: space, Code: Trees[code]})
					}
				}
			}
		}
	}
}

func (w *WorldSystem) generateBuildings(height_map *mapgenerator.HeightMap, width, height int) {
	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			code := rand.Intn(1000)
			if code < len(Buildings) && (height_map.Value(x, y) >= 10 && height_map.Value(x, y) < 30) {
				space := &common.SpaceComponent{
					Position: engo.Point{X: float32(x*(gridSize/2) + y*gridSize/2), Y: float32(y*gridSize/4 - x*gridSize/4)},
					Width:    gridSize,
					Height:   gridSize,
				}
				render := &common.RenderComponent{
					Drawable:    Entitys[Buildings[code]].Drawable,
					Scale:       Entitys[Buildings[code]].Scale,
					StartZIndex: 2,
				}
				ground := w.ground[x][y]
				if ground != nil {
					ground.SubEntites = append(ground.SubEntites, &Grid{BasicEntity: ecs.NewBasic(), RenderComponent: render, SpaceComponent: space, Code: Buildings[code]})
				}
			}
		}
	}
}

func (w *WorldSystem) Size() (float32, float32) {
	return w.width, w.height
}

func (w *WorldSystem) Center() engo.Point {
	return engo.Point{X: w.width / 2, Y: gridSize / 4}
}

func (w *WorldSystem) Bounds() engo.AABB {
	return engo.AABB{
		Min: engo.Point{X: 0, Y: -w.height / 2},
		Max: engo.Point{X: w.width, Y: w.height / 2},
	}
}

func (w *WorldSystem) getGridPos(x, y, z float32) (int, int, int) {
	py := int((x + 2*y) / engo.GetGlobalScale().Y / gridSize)
	px := int((x - 2*y) / engo.GetGlobalScale().X / gridSize)
	pz := int(z + 0.5)
	return px, py, pz
}

func (w *WorldSystem) inView(x, y int) bool {
	return (x >= w.cameraX-w.viewLenX && x < w.cameraX+w.viewLenX) &&
		(y >= w.cameraY-w.viewLenY && y < w.cameraY+w.viewLenY)
}

func (w *WorldSystem) cameraMoved() bool {
	x, y, z := w.getGridPos(w.camera.X(), w.camera.Y(), w.camera.Z())
	if x != w.cameraX || y != w.cameraY || z != w.cameraZ {
		return true
	}
	return false
}

func (w *WorldSystem) Contains(point engo.Point) bool {
	bounds := w.Bounds()
	return point.X > bounds.Min.X && point.X < bounds.Max.X && point.Y > bounds.Min.Y && point.Y < bounds.Max.Y
}

func (w *WorldSystem) CanMove(x, y float32) bool {
	py := int((x+2*y)/gridSize - 0.5)
	px := int((x-2*y)/gridSize + 0.5)
	grid := w.ground[px][py]
	return (grid != nil && grid.Code != Waters[0])
}

func (w *WorldSystem) Update(float32) {
	lastCameraX, lastCameraY := w.cameraX, w.cameraY
	lastViewLenX, lastViewLenY := w.viewLenX, w.viewLenY
	if w.cameraMoved() {
		w.cameraX, w.cameraY, w.cameraZ = w.getGridPos(w.camera.X(), w.camera.Y(), w.camera.Z())
		w.viewLenX = int(engo.GameWidth()/engo.GetGlobalScale().X/gridSize) * (w.cameraZ + 1)
		w.viewLenY = int(engo.GameWidth()/engo.GetGlobalScale().Y/gridSize) * (w.cameraZ + 1)

		count := 0

		for i := lastCameraX - lastViewLenX; i < lastCameraX+lastViewLenX; i++ {
			for j := lastCameraY - lastViewLenY; j < lastCameraY+lastViewLenY; j++ {
				grid := w.ground[i][j]
				if grid != nil && !w.inView(i, j) {
					w.renderer.Remove(grid.BasicEntity)
					if grid.SubEntites != nil {
						for _, e := range grid.SubEntites {
							w.renderer.Remove(e.BasicEntity)
						}
					}
				}
			}
		}

		for i := w.cameraX - w.viewLenX; i < w.cameraX+w.viewLenX; i++ {
			for j := w.cameraY - w.viewLenY; j < w.cameraY+w.viewLenY; j++ {
				grid := w.ground[i][j]
				if grid != nil {
					count += 1
					w.renderer.Add(&grid.BasicEntity, grid.RenderComponent, grid.SpaceComponent)
					if grid.SubEntites != nil {
						for _, e := range grid.SubEntites {
							count += 1
							w.renderer.Add(&e.BasicEntity, e.RenderComponent, e.SpaceComponent)
						}
					}
				} else {
					space := &common.SpaceComponent{
						Position: engo.Point{X: float32(i*(gridSize/2) + j*gridSize/2), Y: float32(j*gridSize/4 - i*gridSize/4)},
						Width:    gridSize,
						Height:   gridSize,
					}
					fillGrid := &Grid{BasicEntity: ecs.NewBasic(), RenderComponent: Entitys[Lands[0]], SpaceComponent: space, Code: Lands[0]}
					if w.ground[i] == nil {
						w.ground[i] = make(map[int]*Grid)
					}
					w.ground[i][j] = fillGrid
					w.renderer.Add(&fillGrid.BasicEntity, fillGrid.RenderComponent, fillGrid.SpaceComponent)
				}
			}
		}
	}
}

func (w *WorldSystem) Remove(basic ecs.BasicEntity) {

}
