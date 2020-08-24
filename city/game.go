package city

import (
	"github.com/EngoEngine/engo"
	"github.com/Z2Y/trpgo/city/scene"
)

func Start(width, height int) {
	opts := engo.RunOptions{
		Title:        "trpgo",
		Width:        375,
		Height:       667,
		MobileWidth:  width,
		MobileHeight: height,
	}
	engo.Run(opts, &scene.MainMenu{})
}
