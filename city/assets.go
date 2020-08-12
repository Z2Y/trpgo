//+build !mobilebind

package city

import (
	"github.com/EngoEngine/engo"
)

func LoadAsset(urls ...string) error {
	return engo.Files.Load(urls...)
}
