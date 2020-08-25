package city

import (
	"github.com/EngoEngine/engo"
	"github.com/Z2Y/trpgo/city/scene"
)

func Start(width, height int) {
	opts := engo.RunOptions{
		Title:        "trpgo",
		Width:        width,
		Height:       height,
		MobileWidth:  width,
		MobileHeight: height,
	}
	engo.Run(opts, &scene.MainMenu{})
}
