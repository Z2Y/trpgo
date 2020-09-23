package city

import (
	"github.com/EngoEngine/engo"
	"github.com/Z2Y/trpgo/city/config"
	"github.com/Z2Y/trpgo/city/scene"
)

func Start(width, height int) {
	gameScale := config.GetSafeScale(float32(width), float32(height))
	opts := engo.RunOptions{
		Title:         "trpgo",
		Width:         width,
		Height:        height,
		MobileWidth:   width,
		MobileHeight:  height,
		FPSLimit:      120,
		ScaleOnResize: true,
		GlobalScale:   engo.Point{X: gameScale, Y: gameScale},
	}

	engo.Run(opts, &scene.MainMenu{})
}
