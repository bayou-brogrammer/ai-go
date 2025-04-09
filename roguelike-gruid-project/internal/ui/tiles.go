//go:build js || sdl
// +build js sdl

package ui

import (
	"fmt"
	"image"

	"golang.org/x/image/font/gofont/gomono"
	"golang.org/x/image/font/opentype"

	"codeberg.org/anaseto/gruid"
	"codeberg.org/anaseto/gruid/tiles"
)

// TileDrawer implements TileManager from the gruid-sdl module.
type TileDrawer struct {
	drawer *tiles.Drawer
}

// GetImage implements TileManager.GetImage.
func (t *TileDrawer) GetImage(c gruid.Cell) image.Image {
	fgColor := ColorToRGBA(c.Style.Fg, true)
	bgColor := ColorToRGBA(c.Style.Bg, false)

	// Handle style attributes
	if c.Style.Attrs&AttrReverse != 0 {
		fgColor, bgColor = bgColor, fgColor
	}

	// We return an image with the given rune drawn using the previously
	// defined foreground and background colors.
	return t.drawer.Draw(c.Rune, image.NewUniform(fgColor), image.NewUniform(bgColor))
}

// TileSize implements TileManager.TileSize. It returns the tile size, in
// pixels. In this tutorial, it corresponds to the size of a character with the
// font we use.
func (t *TileDrawer) TileSize() gruid.Point {
	return t.drawer.Size()
}

// GetTileDrawer returns a TileDrawer that implements TileManager for the sdl
// driver, or an error if there were problems setting up the font face.
func GetTileDrawer() (*TileDrawer, error) {
	parsedFont, err := opentype.Parse(gomono.TTF)
	if err != nil {
		return nil, fmt.Errorf("failed to parse font: %w", err)
	}

	t := &TileDrawer{}

	face, err := opentype.NewFace(parsedFont, &opentype.FaceOptions{
		Size: 24,
		DPI:  72,
	})
	if err != nil {
		return nil, err
	}

	t.drawer, err = tiles.NewDrawer(face)
	if err != nil {
		return nil, err
	}
	return t, nil
}
