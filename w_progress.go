package yagk

import (
	"github.com/hajimehoshi/ebiten/v2"
	"image"
	"math/rand"
)

type Progress struct {
	id       WidgetId
	rect     image.Rectangle
	total    int
	progress *int
}

func newProgress(total, x, y, w, h int) (p *Progress) {
	p = &Progress{
		id:    WidgetId(rand.Int()),
		rect:  image.Rect(x, y, x+w, y+h),
		total: total,
	}
	return
}

func (p *Progress) Update(io *IO) {
	if p.progress != nil && *p.progress > p.total {
		*p.progress = p.total
	}
}

func (p *Progress) Draw(canvas *ebiten.Image, win *Window) {
	win.drawStyle(canvas, p.rect, 3, 0)
	win.drawStyle(canvas, p.rect, 3, 2)

	if p.progress != nil && *p.progress > 0 {
		xs := *p.progress * (p.rect.Dx() / p.total)
		c := image.Rect(p.rect.Min.X, p.rect.Min.Y, p.rect.Min.X+xs, p.rect.Max.Y)
		win.drawStyle(canvas, c, 3, 1)
	}
}

func (p *Progress) Translate(x, y int) {
}
