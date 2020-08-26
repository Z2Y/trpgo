package scene

import (
	"github.com/Z2Y/trpgo/city/asset"
	"github.com/Z2Y/trpgo/city/core/ui"
	"github.com/Z2Y/trpgo/city/core/ui/layout"

	"github.com/EngoEngine/ecs"
	"github.com/EngoEngine/engo"
	"github.com/EngoEngine/engo/common"
)

type MainMenu struct {
}

func (s *MainMenu) Type() string {
	return "MainMenu"
}

func (s *MainMenu) Preload() {
	asset.LoadAsset("ui/blueSheet.xml")
	asset.LoadAsset("font/CN.ttf")
}

func (s *MainMenu) Setup(u engo.Updater) {
	w := u.(*ecs.World)

	common.SetBackground(background)
	ui.SetDefaultFont("font/CN.ttf")
	w.AddSystem(&common.RenderSystem{})
	w.AddSystemInterface(&ui.UISystem{}, ui.UIEntityFace, nil)

	buttonBg := asset.LoadedSubSprite("blue_button00.png")
	text := ui.NewText(ui.Text{Value: "开始游戏"})
	button := ui.NewButton(ui.Button{Text: text, Image: buttonBg, Width: 150, Height: 50})

	button.SetPosition(layout.AlignToWorldCenter(button.SpaceComponent.AABB()))

	button.OnClick(func() {
		engo.SetScene(&Game{}, true)
	})

	w.AddEntity(button)
	w.AddEntity(text)
}
