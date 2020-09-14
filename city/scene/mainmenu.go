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
	asset.LoadAsset("background.jpg")
	asset.LoadAsset("font/CN.ttf")
}

func (s *MainMenu) Setup(u engo.Updater) {
	w := u.(*ecs.World)

	common.SetBackground(background)
	engo.Window.SetAspectRatio(667, 375)
	ui.SetDefaultFont("font/CN.ttf", 24)
	w.AddSystem(&common.RenderSystem{})
	w.AddSystemInterface(&ui.UISystem{}, ui.UIEntityFace, nil)

	buttonBg := asset.LoadedSubSprite("blue_button00.png")
	text := ui.NewText(ui.Text{Value: "开始游戏"})
	button := ui.NewButton(ui.Button{Text: text, Image: buttonBg, Width: 100, Height: 36})
	background := ui.NewImage(ui.Image{Texture: asset.LoadedSprite("background.jpg"), Width: engo.WindowWidth() / engo.GetGlobalScale().X, Height: engo.WindowHeight() / engo.GetGlobalScale().Y})

	button.SetPosition(layout.AlignToWorldCenter(button.SpaceComponent.AABB()))
	background.SetPosition(layout.AlignToWorldCenter(background.SpaceComponent.AABB()))
	background.SetZIndex(ui.UILayerIndex - 1)

	button.OnClick(func() {
		engo.SetScene(&Game{}, true)
	})

	w.AddEntity(background)
	w.AddEntity(button)
	w.AddEntity(text)

}
