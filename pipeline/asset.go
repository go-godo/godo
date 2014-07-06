package pipeline

import (
	"bytes"

	"github.com/mgutz/gosu"
)

// Asset is any file which can be loaded into memory for processing in the
// pipeline.
type Asset struct {
	bytes.Buffer
	Info      gosu.FileAsset
	ReadPath  string
	WritePath string
	Accepts   string
}

// Assets is a collection of Asset
type Assets struct {
	Assets []Asset
	Src    string
	Dest   string
}

// NewAssets creates a new instance and loads all sources.
func NewAssets(src string, dest string) *Assets {
	assets := &Assets{Src: src, Dest: dest}
	return assets
}
