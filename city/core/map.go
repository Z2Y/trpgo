package core

import (
	"github.com/EngoEngine/engo"
	"github.com/EngoEngine/engo/common"
)

type WorldMap struct {
	xLen      int
	yLen      int
	base      [][]int
	buildings [][]int
}

var SampleWorldMap = WorldMap{
	xLen: 100,
	yLen: 100,
	base: [][]int{},
	buildings: [][]int{
		{0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
	},
}

var Entitys = map[int]*common.RenderComponent{}

func RegistRenderComponent(code int, d *common.RenderComponent) {
	Entitys[code] = d
}

func InitRenderComponents() {
	// land1, _ := common.LoadedSprite("land/land_1.png")
	land := common.NewSpritesheetFromFile("land.png", 64, 64)
	landScale := engo.Point{X: 1, Y: 1}
	for i, cell := range land.Cells() {
		RegistRenderComponent(i, &common.RenderComponent{Scale: landScale, Drawable: cell})
	}
}
