package yagk

import (
	"github.com/hajimehoshi/ebiten/v2"
	"image"
)

type Icon struct {
	rect   *image.Rectangle
	iconId int
	Scale  float64
	style  *StyleDefinition
	x, y   int
}

func newIcon(id, x, y int, style *StyleDefinition) (i *Icon) {
	i = &Icon{
		Scale:  style.Icons.Scale,
		x:      x,
		y:      y,
		style:  style,
	}
	i.setId(id)
	return
}

func (i *Icon) setId(id int) {
	i.iconId = id
	r := i.iconId / i.style.Icons.Columns
	c := i.iconId % i.style.Icons.Columns
	i.rect = &image.Rectangle{
		Min: image.Point{
			X: c * i.style.Icons.Size,
			Y: r * i.style.Icons.Size,
		},
		Max: image.Point{
			X: (c * i.style.Icons.Size) + i.style.Icons.Size,
			Y: (r * i.style.Icons.Size) + i.style.Icons.Size,
		},
	}
}

func (i *Icon) Update(io *IO) {
}

func (i *Icon) Draw(canvas *ebiten.Image, win *Window) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(i.Scale, i.Scale)
	op.GeoM.Translate(float64(i.x), float64(i.y))
	canvas.DrawImage(win.peekPane().style.icons.SubImage(*i.rect).(*ebiten.Image), op)
}

func (i *Icon) Translate(x, y int) {
}
