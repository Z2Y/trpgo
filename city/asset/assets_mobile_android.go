// +build android

package asset

import (
	"fmt"
	"io"
	"log"
	"path/filepath"
	"strings"

	"github.com/EngoEngine/engo"
	"github.com/EngoEngine/engo/common"
	"golang.org/x/mobile/asset"
)

func openFile(url string) (io.ReadCloser, error) {
	usedUrl := url
	if strings.HasPrefix(url, "assets/") {
		usedUrl = usedUrl[7:]
	}

	return asset.Open(usedUrl)
}

func load(url string) error {
	f, err := openFile(filepath.Join(engo.Files.GetRoot(), url))

	if err != nil {
		log.Println(fmt.Sprintf("unable to open resource: %s", err))
		return fmt.Errorf("unable to open resource: %s", err)
	}

	defer f.Close()

	return engo.Files.LoadReaderData(url, f)
}

func LoadAsset(urls ...string) error {
	for _, url := range urls {
		err := load(url)
		if err != nil {
			return err
		}
	}
	return nil
}

func LoadedSprite(url string) *common.Texture {
	sprite, err := common.LoadedSprite(url)

	if err != nil {
		panic(err)
	}

	return sprite
}

func LoadedSubSprite(url string) *common.Texture {
	return LoadedSprite(subtextureURL(url))
}
