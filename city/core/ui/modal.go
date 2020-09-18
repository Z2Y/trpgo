package ui

import (
	"image/color"

	"github.com/EngoEngine/ecs"
	"github.com/EngoEngine/engo"
	"github.com/EngoEngine/engo/common"
	"github.com/Z2Y/trpgo/city/asset"
	"github.com/Z2Y/trpgo/city/core/ui/layout"
)

type Modal struct {
	UIBasic

	closeBtn *Button
	Content  UIFace
}

func NewModal(modal Modal) *Modal {
	modal.BasicEntity = ecs.NewBasic()
	modal.RenderComponent = common.RenderComponent{
		Drawable:    common.Rectangle{},
		Color:       color.RGBA{0x00, 0x00, 0x00, 0x94},
		StartZIndex: UIPopLayerIndex,
	}
	modal.SpaceComponent = common.SpaceComponent{
		Width:  modal.Width,
		Height: modal.Height,
	}

	contentEntity := modal.Content.GetComponentUI()
	contentEntity.SetZIndex(UIPopLayerIndex + 1)
	contentEntity.SetPosition(layout.AlignCenter(modal.AABB(), contentEntity.AABB()))

	modal.closeBtn = NewButton(Button{UIBasic: UIBasic{Width: 20, Height: 20}, Image: asset.LoadedSubSprite("blue_boxCross.png")})
	modal.closeBtn.SetPosition(layout.AlignRightTop(contentEntity.AABB(), modal.closeBtn.AABB(), 10, 10))
	modal.closeBtn.SetZIndex(UIPopLayerIndex + 1)
	modal.closeBtn.OnClick(modal.Close)

	modal.AddSubEntity(modal.closeBtn)
	modal.AddSubEntity(modal.Content)
	modal.SetShader(common.HUDShader)
	modal.EventResponder = true
	return &modal
}

func (modal *Modal) Close() {
	engo.Mailbox.Dispatch(UICloseEvent{Target: &modal.UIBasic})
}

func (modal *Modal) OnClick(f func()) {
	if modal.MessageListener == nil {
		modal.MessageListener = &engo.MessageManager{}
	}
	modal.MessageListener.Listen("UIMouseEvent", func(engo.Message) {
		modal.Close()
	})
}
