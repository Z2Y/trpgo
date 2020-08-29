package city

import (
	"github.com/EngoEngine/engo"
	"github.com/Z2Y/trpgo/city/scene"
)

var (
	GameWidth  = 667
	GameHeight = 375
)

func getSafeScale(width, height float32) float32 {
	gameRatio := float32(GameHeight) / float32(GameWidth)
	safeHeight := width * gameRatio
	if safeHeight > height {
		return height / float32(GameHeight)
	}
	return width / float32(GameWidth)
}

func Start(width, height int) {
	gameScale := getSafeScale(float32(width), float32(height))
	opts := engo.RunOptions{
		Title:        "trpgo",
		Width:        width,
		Height:       height,
		MobileWidth:  width,
		MobileHeight: height,
		FPSLimit:     120,
		GlobalScale:  engo.Point{X: gameScale, Y: gameScale},
	}

	engo.Run(opts, &scene.MainMenu{})
}
