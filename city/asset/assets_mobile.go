// +build darwin
// +build ios

package asset

/*
#cgo CFLAGS: -x objective-c
#cgo LDFLAGS: -framework Foundation
#import <Foundation/Foundation.h>

const char* mainBundlePath(void);
*/
import "C"

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/EngoEngine/engo"
)

func mainBundlePath() string {
	return C.GoString(C.mainBundlePath())
}

func openFile(url string) (io.ReadCloser, error) {
	usedUrl := url
	if strings.HasPrefix(url, "assets/") {
		usedUrl = usedUrl[7:]
	}

	if !filepath.IsAbs(usedUrl) {
		bundlePath := mainBundlePath()
		usedUrl = filepath.Join(bundlePath, "assets", usedUrl)
	}

	f, err := os.Open(usedUrl)

	if err != nil {
		return nil, err
	}

	return f, nil
}

func load(url string) error {
	f, err := openFile(filepath.Join(engo.Files.GetRoot(), url))

	if err != nil {
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
