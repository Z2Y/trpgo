package scene

import (
	"image/color"

	"github.com/EngoEngine/ecs"
	"github.com/EngoEngine/engo"
	"github.com/EngoEngine/engo/common"

	"github.com/Z2Y/trpgo/city/asset"
	"github.com/Z2Y/trpgo/city/core"
	"github.com/Z2Y/trpgo/city/core/action"
	"github.com/Z2Y/trpgo/city/core/ui"
	"github.com/Z2Y/trpgo/city/human"
)

var (
	background = color.RGBA{0xfa, 0xf8, 0xef, 0xff}
)

type Game struct {
	world *core.WorldSystem
}

func NewGame() *Game {
	g := &Game{}
	return g
}

func (g *Game) Preload() {
	asset.LoadAsset("npc/idle.png")
	asset.LoadAsset("npc/walking.png")
	asset.LoadAsset("land.png")
	asset.LoadAsset("ui/blueSheet.xml")
	asset.LoadAsset("font/CN.ttf")
	human.Init()
	core.InitRenderComponents()
}

func (g *Game) Setup(u engo.Updater) {
	common.SetBackground(background)
	w := u.(*ecs.World)

	g.world = &core.WorldSystem{}
	g.world.LoadWorldMap(&core.SampleWorldMap)
	ui.SetDefaultFont("font/CN.ttf")

	hero := human.NewHuman(g.world.Center())

	w.AddSystem(&common.RenderSystem{})
	w.AddSystem(&common.FPSSystem{Display: false})
	w.AddSystem(&common.AnimationSystem{})
	w.AddSystem(g.world)
	w.AddSystem(&core.ControlSystem{ZoomSpeed: -0.125})
	w.AddSystem(&action.AutoActionSystem{})

	for _, system := range w.Systems() {
		switch sys := system.(type) {
		case *common.RenderSystem:
			sys.Add(&hero.BasicEntity, &hero.RenderComponent, &hero.SpaceComponent)
		case *common.AnimationSystem:
			sys.Add(&hero.BasicEntity, &hero.AnimationComponent, &hero.RenderComponent)
		case *action.AutoActionSystem:
			sys.Add(&hero.ActionEntity)
		}
	}
}

func (g *Game) Type() string {
	return "city game"
}
