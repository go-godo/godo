package gosu

import (
	"bytes"
)

type Asset struct {
	bytes.Buffer
	Info      FileAsset
	ReadPath  string
	WritePath string
	Accepts   string
}

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
