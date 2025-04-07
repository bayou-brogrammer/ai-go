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

// TileDrawer implements TileManager from the gruid-sdl module. It is used to
// provide a mapping from virtual grid cells to images using the tiles package.
// In this tutorial, we just draw a font with a given foreground and
// background, but it would be possible to make a tiles version with custom
// drawings for cells.
type TileDrawer struct {
	drawer *tiles.Drawer
}

// GetImage implements TileManager.GetImage.
func (t *TileDrawer) GetImage(c gruid.Cell) image.Image {
	// Default colors for SDL/JS
	// fg := image.NewUniform(color.RGBA{0xad, 0xbc, 0xbc, 255}) // Light gray
	// bg := image.NewUniform(color.RGBA{0x10, 0x3c, 0x48, 255}) // Dark blue-gray

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
	// Parse the font bytes
	parsedFont, err := opentype.Parse(gomono.TTF)
	if err != nil {
		return nil, fmt.Errorf("failed to parse font: %w", err)
	}

	t := &TileDrawer{}

	// We retrieve a font face.
	face, err := opentype.NewFace(parsedFont, &opentype.FaceOptions{
		Size: 24,
		DPI:  72,
	})
	if err != nil {
		return nil, err
	}
	// We create a new drawer for tiles using the previous face. Note that
	// if more than one face is wanted (such as an italic or bold variant),
	// you would have to create drawers for thoses faces too, and then use
	// the relevant one accordingly in the GetImage method.
	t.drawer, err = tiles.NewDrawer(face)
	if err != nil {
		return nil, err
	}
	return t, nil
}
