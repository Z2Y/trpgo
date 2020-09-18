package scene

import (
	"image/color"
	"log"

	"github.com/EngoEngine/ecs"
	"github.com/EngoEngine/engo"
	"github.com/EngoEngine/engo/common"

	"github.com/Z2Y/trpgo/city/asset"
	"github.com/Z2Y/trpgo/city/core"
	"github.com/Z2Y/trpgo/city/core/control"
	"github.com/Z2Y/trpgo/city/core/ui"
	"github.com/Z2Y/trpgo/city/core/ui/layout"
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
	asset.LoadAsset("land/foliagePack_default.xml")
	asset.LoadAsset("building/building.xml")
	asset.LoadAsset("land.png")
	human.Init()
	core.InitRenderComponents()
}

func (g *Game) Setup(u engo.Updater) {
	common.SetBackground(background)
	w := u.(*ecs.World)

	g.world = &core.WorldSystem{}
	g.world.LoadWorldMap(&core.SampleWorldMap)
	// ui.SetDefaultFont("font/CN.ttf")

	hero := human.NewHuman(g.world.Center())

	w.AddSystem(&common.RenderSystem{})
	w.AddSystem(&common.FPSSystem{Display: true, Font: ui.DefaultFont})
	w.AddSystem(&common.AnimationSystem{})
	w.AddSystem(g.world)
	w.AddSystem(&control.ControlSystem{ZoomSpeed: -0.125})
	w.AddSystem(&control.AutoActionSystem{})

	g.SetupUI(w)

	for _, system := range w.Systems() {
		switch sys := system.(type) {
		case *common.RenderSystem:
			sys.Add(&hero.BasicEntity, &hero.RenderComponent, &hero.SpaceComponent)
		case *common.AnimationSystem:
			sys.Add(&hero.BasicEntity, &hero.AnimationComponent, &hero.RenderComponent)
		case *control.ControlSystem:
			sys.Add(&hero.ActionEntity)
		case *control.AutoActionSystem:
			sys.Add(&hero.ActionEntity)
		case *control.WalkSystem:
			sys.Add(&hero.ActionEntity)
		}
	}
}

func (g *Game) SetupUI(w *ecs.World) {
	w.AddSystemInterface(&ui.UISystem{}, ui.UIEntityFace, nil)

	buttonFnt := ui.NewTextFont("font/CN.ttf", 12)

	buttonBg := asset.LoadedSubSprite("blue_button00.png")
	text := ui.NewText(ui.Text{Value: "系统", Font: buttonFnt})
	button := ui.NewButton(ui.Button{Text: text, Image: buttonBg, UIBasic: ui.UIBasic{Width: 40, Height: 20}})

	buildBtn := ui.NewButton(ui.Button{Text: ui.NewText(ui.Text{Value: "建造", Font: buttonFnt}), Image: buttonBg, UIBasic: ui.UIBasic{Width: 40, Height: 20}})

	w.AddEntity(button)
	w.AddEntity(buildBtn)

	button.SetPosition(layout.AlignToWorldRightBottom(button.SpaceComponent.AABB(), 8, 8))
	buildBtn.SetPosition(layout.AlignToWorldRightBottom(buildBtn.SpaceComponent.AABB(), 56, 8))

	button.OnClick(func() {
		log.Println("系统 clicked")

		panel := ui.NewPanel(ui.Panel{Image: asset.LoadedSubSprite("blue_panel.png"), UIBasic: ui.UIBasic{Width: 400, Height: 300}})
		modal := ui.NewModal(ui.Modal{Content: panel, UIBasic: ui.UIBasic{Width: engo.WindowWidth() / engo.GetGlobalScale().X, Height: engo.WindowHeight() / engo.GetGlobalScale().Y}})
		w.AddEntity(modal)
	})
}

func (g *Game) Type() string {
	return "city game"
}
