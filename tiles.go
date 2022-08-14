package main
import (
	"image"
	"image/color"

	"golang.org/x/image/font/gofont/gomono"
	"golang.org/x/image/font/opentype"

	"github.com/anaseto/gruid"
	"github.com/anaseto/gruid/tiles"
)

// Tile implements TileManager.
type TileDrawer struct {
	drawer *tiles.Drawer
}

func GetTileDrawer() (*TileDrawer, error) {
	t := &TileDrawer{}
	var err error
	// We get a monospace font TTF.
	font, err := opentype.Parse(gomono.TTF)
	if err != nil {
		return nil, err
	}
	// We retrieve a font face.
	face, err := opentype.NewFace(font, &opentype.FaceOptions{
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

func (t *TileDrawer) GetImage(c gruid.Cell) image.Image {
	// we use some selenized colors
	fg := image.NewUniform(color.RGBA{0xad, 0xbc, 0xbc, 255})
	bg := image.NewUniform(color.RGBA{0x18, 0x49, 0x56, 255})
	switch c.Style.Bg {
	case colorFOV:
		bg = image.NewUniform(color.RGBA{0x18, 0x49, 0x56, 255})
	case colorPlayer:
		fg = image.NewUniform(color.RGBA{0x46, 0x95, 0xf7, 255})
	case colorEnemy:
		fg = image.NewUniform(color.RGBA{0xfa, 0x57, 0x50, 255})
    case colorLogPlayerAttack, colorStatusHealthy:
        fg = image.NewUniform(color.RGBA{0x75, 0xb9, 0x38, 255})
    case colorLogEnemyAttack, colorStatusWounded:
        fg = image.NewUniform(color.RGBA{0xed, 0x86, 0x49, 255})
    case colorLogSpecial:
        fg = image.NewUniform(color.RGBA{0xf2, 0x75, 0xbe, 255})
	}
	
	return t.drawer.Draw(c.Rune, fg, bg)
}

func (t *TileDrawer) TileSize() gruid.Point {
	return t.drawer.Size()
}
