// +build darwin
// +build ios

package asset

import (
	"encoding/xml"
	"fmt"
	"io"
	"path"
	"strings"

	"github.com/EngoEngine/engo"
	"github.com/EngoEngine/engo/common"
)

// fix origin textureAtlasLoader couldn't work in mobile
var texAtlasLoader *textureAtlasLoader

// TextureAtlasResource contains reference to a loaded TextureAtlas and the texture of the main image
type TextureAtlasResource struct {
	// url is the location of the xml file
	url string
	// Atlas is the TextureAtlas filled with data from the parsed XML file
	Atlas *common.TextureAtlas
}

// URL retrieves the url to the .xml file
func (r TextureAtlasResource) URL() string {
	return r.url
}

// textureAtlasLoader is reponsible for managing '.xml' files exported from TexturePacker (https://www.codeandweb.com/texturepacker)
type textureAtlasLoader struct {
	atlases map[string]*TextureAtlasResource
	images  map[string]common.TextureResource
}

// Load will load the xml file and the main image as well as add references
// for sub textures/images in engo.Files, subtextures keep their path url (with appended extension from main image path if it does not exist),
// the main image is loaded in reference to the directory of the xml file
// For example this sub texture:
//  <SubTexture name="subimg" x="10" y="10" width="50" height="50"/>
// can be retrieved with this go code
//  texture, err := common.LoadedSprite("subimg.png")
func (t *textureAtlasLoader) Load(url string, data io.Reader) error {
	atlas, err := t.createAtlasFromXML(data, url)
	if err != nil {
		return err
	}

	t.atlases[url] = atlas
	return nil
}

// Unload removes the preloaded atlass from the cache and clears
// references to all SubTextures from the image loader
func (t *textureAtlasLoader) Unload(url string) error {
	imgURL := path.Join(path.Dir(url), t.atlases[url].Atlas.ImagePath)
	if err := engo.Files.Unload(imgURL); err != nil {
		return err
	}
	for _, subTexture := range t.atlases[url].Atlas.SubTextures {
		delete(t.images, subTexture.Name)
	}

	delete(t.atlases, url)
	return nil
}

// Resource retrieves and returns the texture atlas of type TextureAtlasResource
func (t *textureAtlasLoader) Resource(url string) (engo.Resource, error) {
	ext := path.Ext(url)

	if ext == ".subtexture" {
		return t.SubTextureResource(url)
	}

	atlas, ok := t.atlases[url]
	if !ok {
		return nil, fmt.Errorf("resource not loaded by `FileLoader`: %q", url)
	}

	return atlas, nil
}

func (t *textureAtlasLoader) SubTextureResource(url string) (engo.Resource, error) {
	texture, ok := t.images[url]
	if !ok {
		return nil, fmt.Errorf("resource not loaded by `FileLoader`: %q", url)
	}
	return texture, nil
}

// createAtlasFromXML unmarshals and unpacks the xml data into a TextureAtlas
// it also adds the main image and subtextures to the imageLoader
// if the subtexture doesn't have an extension in it's Name field,
// it will append the main image's extension to it
func (t *textureAtlasLoader) createAtlasFromXML(r io.Reader, url string) (*TextureAtlasResource, error) {
	var atlas *common.TextureAtlas
	err := xml.NewDecoder(r).Decode(&atlas)
	if err != nil {
		return nil, err
	}

	imgURL := path.Join(path.Dir(url), atlas.ImagePath)
	if err := load(imgURL); err != nil {
		return nil, fmt.Errorf("failed load texture atlas image: %v", err)
	}

	res, err := engo.Files.Resource(imgURL)
	if err != nil {
		return nil, err
	}

	img, ok := res.(common.TextureResource)
	if !ok {
		return nil, fmt.Errorf("resource not of type `TextureResource`: %v", url)
	}

	for i, subTexture := range atlas.SubTextures {
		viewport := engo.AABB{
			Min: engo.Point{
				X: subTexture.X / img.Width,
				Y: subTexture.Y / img.Height,
			},
			Max: engo.Point{
				X: (subTexture.X + subTexture.Width) / img.Width,
				Y: (subTexture.Y + subTexture.Height) / img.Height,
			},
		}

		subURL := subtextureURL(subTexture.Name)

		atlas.SubTextures[i].Name = subURL

		t.images[subURL] = common.TextureResource{Texture: img.Texture, Width: subTexture.Width, Height: subTexture.Height, Viewport: &viewport}
	}

	return &TextureAtlasResource{
		Atlas: atlas,
		url:   url,
	}, nil

}

func subtextureURL(url string) string {
	subURL := url
	subExt := path.Ext(url)
	if subExt == "" {
		subURL += ".subtexture"
	} else {
		subURL = strings.Replace(subURL, subExt, ".subtexture", 1)
	}
	return subURL
}

func init() {
	texAtlasLoader = &textureAtlasLoader{atlases: make(map[string]*TextureAtlasResource), images: make(map[string]common.TextureResource)}
	engo.Files.Register(".xml", texAtlasLoader)
	engo.Files.Register(".subtexture", texAtlasLoader)
}
