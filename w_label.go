package yagk

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font"
	"image"
	"image/color"
)

type Label struct {
	rect      image.Rectangle
	text      string
	Color     color.Color
	Alignment Alignment
}

func newLabel(text string, x, y, w, h int) (l *Label) {
	l = &Label{
		rect:      image.Rect(x, y, x+w, y+h),
		text:      text,
		Color:     color.Black,
		Alignment: AlignLeft,
	}
	return
}

func (l *Label) Update(io *IO) {}

func (l *Label) Draw(canvas *ebiten.Image, win *Window) {
	bounds, _ := font.BoundString(win.peekPane().style.font, l.text)

	w := (bounds.Max.X - bounds.Min.X).Ceil()
	y := l.rect.Max.Y - (l.rect.Dy()-win.peekPane().style.fontHeight)/2
	x := l.rect.Min.X + win.peekPane().style.def.Offset.X

	if l.Alignment == AlignCenter {
		x = l.rect.Min.X + (l.rect.Dx()-w)/2
	} else if l.Alignment == AlignRight {
		x = (l.rect.Max.X - w) - (win.peekPane().style.def.Offset.X * 2)
	}

	text.Draw(canvas, l.text, win.peekPane().style.font, x, y, l.Color)
}

func (l *Label) Translate(x, y int) {
}
