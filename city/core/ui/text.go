package ui

import (
	"fmt"
	"image/color"

	"github.com/EngoEngine/ecs"
	"github.com/EngoEngine/engo"
	"github.com/EngoEngine/engo/common"
)

var (
	DefaultFont *common.Font
)

type Text struct {
	UIBasic
	Font     *common.Font
	Position engo.Point
	Value    string
}

func NewText(position engo.Point, value string) *Text {
	text := Text{Position: position, Value: value, Font: DefaultFont}

	text.BasicEntity = ecs.NewBasic()

	width, height, _ := text.Font.TextDimensions(value)

	text.SpaceComponent = common.SpaceComponent{
		Width:    float32(width),
		Height:   float32(height),
		Position: position,
	}

	text.RenderComponent.Drawable = common.Text{
		Font: text.Font,
		Text: text.Value,
	}

	text.SetShader(common.TextHUDShader)

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

	t.RenderComponent.Drawable = common.Text{
		Font: t.Font,
		Text: t.Value,
	}

	width, height, _ := t.Font.TextDimensions(t.Value)

	t.SpaceComponent.Width = float32(width)
	t.SpaceComponent.Height = float32(height)

	return nil
}

func SetDefaultFont(URL string) {
	fnt := &common.Font{
		URL:  URL,
		FG:   color.Black,
		Size: 64,
	}
	err := fnt.CreatePreloaded()

	if err != nil {
		panic(err)
	}

	DefaultFont = fnt
}
