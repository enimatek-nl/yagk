package yagk

import (
	"bytes"
	"github.com/hajimehoshi/ebiten/v2"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
	"image"
	"image/color"
	_ "image/png"
	"log"
)

type Size struct {
	X, Y int
}

type Font struct {
	Color  color.Color
	Button color.Color
	Data   []byte
	Size   float64
	DPI    float64
}

type Icons struct {
	Data    []byte
	Size    int
	Scale   float64
	Columns int
}

type StyleDefinition struct {
	Offset     Size
	Sizes      Size
	Font       Font
	Icons      Icons
	Background color.Color
	Data       []byte
	Size       int
}

type Style struct {
	def                   *StyleDefinition
	data, icons           *ebiten.Image
	font                  font.Face
	fontWidth, fontHeight int
}

func NewStyle(def *StyleDefinition) (s *Style) {
	s = &Style{
		def: def,
	}

	i0, _, err := image.Decode(bytes.NewReader(s.def.Icons.Data))
	if err != nil {
		log.Fatal(err)
	}
	s.icons = ebiten.NewImageFromImage(i0)

	i1, _, err := image.Decode(bytes.NewReader(s.def.Data))
	if err != nil {
		log.Fatal(err)
	}
	s.data = ebiten.NewImageFromImage(i1)

	i2, err := opentype.Parse(s.def.Font.Data)
	if err != nil {
		log.Fatal(err)
	}

	s.font, err = opentype.NewFace(i2, &opentype.FaceOptions{
		Size:    s.def.Font.Size,
		DPI:     s.def.Font.DPI,
		Hinting: font.HintingFull,
	})
	if err != nil {
		return
	}

	// calculate the Height of this font based on the capital M
	b, _, _ := s.font.GlyphBounds('M')
	s.fontHeight = (b.Max.Y - b.Min.Y).Ceil()
	s.fontWidth = (b.Max.X - b.Min.X).Ceil()

	return
}
