//+build !mobilebind

package asset

import (
	"github.com/EngoEngine/engo"
	"github.com/EngoEngine/engo/common"
)

func LoadAsset(urls ...string) error {
	return engo.Files.Load(urls...)
}

func LoadedSprite(url string) *common.Texture {
	sprite, err := common.LoadedSprite(url)

	if err != nil {
		panic(err)
	}

	return sprite
}

func LoadedSubSprite(url string) *common.Texture {
	return LoadedSprite(url)
}
