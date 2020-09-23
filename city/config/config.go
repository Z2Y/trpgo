package config

import (
	"github.com/EngoEngine/engo"
)

var (
	GameWidth  = float32(667)
	GameHeight = float32(375)

	windowScale = engo.Point{X: 1, Y: 1}
)

func GetSafeScale(width, height float32) float32 {
	gameRatio := float32(GameHeight) / float32(GameWidth)
	safeHeight := width * gameRatio
	if safeHeight > height {
		return height / float32(GameHeight)
	}
	return width / float32(GameWidth)
}

func UpdateWindowScale(engo.Message) {
	curW, curH := engo.GameWidth(), engo.GameHeight()
	windowScale.X = engo.WindowWidth() / curW
	windowScale.Y = engo.WindowHeight() / curH
}

func WindowScale() engo.Point {
	return windowScale
}
