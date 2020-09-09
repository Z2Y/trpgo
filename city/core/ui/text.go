package ui

import (
	"fmt"
	"image/color"

	"github.com/EngoEngine/ecs"
	"github.com/EngoEngine/engo"
	"github.com/EngoEngine/engo/common"
)

var (
	DefaultFont    *common.Font
	TextLayerIndex = float32(100)
)

type Text struct {
	UIBasic
	Font     *common.Font
	Position engo.Point
	Value    string
}

func NewText(text Text) *Text {

	text.BasicEntity = ecs.NewBasic()

	if text.Font == nil {
		text.Font = DefaultFont
	}

	width, height, _ := text.Font.TextDimensions(text.Value)

	text.SpaceComponent = common.SpaceComponent{
		Width:    float32(width),
		Height:   float32(height),
		Position: text.Position,
	}

	text.RenderComponent.Drawable = text.Font.Render(text.Value)

	text.SetShader(common.HUDShader)
	text.SetZIndex(TextLayerIndex)

	return &text
}

func (t *Text) SetText(value string) error {
	if t.Value == value {
		return nil
	}
	t.Value = value
	return t.update()
}

func (t *Text) SetFont(font *common.Font) error {
	t.Font = font

	return t.update()
}

func (t *Text) update() error {

	if t.Font == nil {
		return fmt.Errorf("Text update without setting Font")
	}

	t.RenderComponent.Drawable = t.Font.Render(t.Value)

	width, height, _ := t.Font.TextDimensions(t.Value)

	t.SpaceComponent.Width = float32(width)
	t.SpaceComponent.Height = float32(height)

	return nil
}

func SetDefaultFont(URL string) {
	fnt := &common.Font{
		URL:  URL,
		FG:   color.White,
		Size: 24,
	}
	err := fnt.CreatePreloaded()

	if err != nil {
		panic(err)
	}

	DefaultFont = fnt
}
