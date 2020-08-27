package core

import (
	"github.com/EngoEngine/engo"
	"github.com/EngoEngine/engo/common"
	"github.com/Z2Y/trpgo/city/asset"
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

var (
	Entitys = map[int]*common.RenderComponent{}
	Lands   = []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11}
	Waters  = []int{13}
	Sand    = []int{12}
	Trees   = []int{101, 102, 103, 104, 105}
)

func RegistRenderComponent(code int, d *common.RenderComponent) {
	Entitys[code] = d
}

func InitRenderComponents() {
	land := common.NewSpritesheetFromFile("land.png", 64, 64)
	landScale := engo.Point{X: 1, Y: 1}
	for i, cell := range land.Cells() {
		RegistRenderComponent(i, &common.RenderComponent{Scale: landScale, Drawable: cell})
	}

	treeScale := engo.Point{X: 0.5, Y: 0.5}
	RegistRenderComponent(101, &common.RenderComponent{Scale: treeScale, Drawable: asset.LoadedSubSprite("foliagePack_004.png")})
	RegistRenderComponent(102, &common.RenderComponent{Scale: treeScale, Drawable: asset.LoadedSubSprite("foliagePack_005.png")})
	RegistRenderComponent(103, &common.RenderComponent{Scale: treeScale, Drawable: asset.LoadedSubSprite("foliagePack_006.png")})
	RegistRenderComponent(104, &common.RenderComponent{Scale: treeScale, Drawable: asset.LoadedSubSprite("foliagePack_007.png")})
	RegistRenderComponent(105, &common.RenderComponent{Scale: treeScale, Drawable: asset.LoadedSubSprite("foliagePack_008.png")})
	RegistRenderComponent(106, &common.RenderComponent{Scale: treeScale, Drawable: asset.LoadedSubSprite("foliagePack_012.png")})
}
