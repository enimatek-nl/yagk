package yagk

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text"
	"image"
	"image/color"
	"math/rand"
	"strings"
	"time"
)

type Input struct {
	id                 WidgetId
	rect               image.Rectangle
	sx, cx, ss         int // start x, cursor x, start selection
	text               *string
	color              color.Color
	timer              time.Time
	flicker            bool
	hover, drag, focus bool
	style              *Style
}

func newInput(text *string, x, y, w, h int, style *Style) (i *Input) {
	i = &Input{
		id:    WidgetId(rand.Int()),
		rect:  image.Rect(x, y, x+w, y+h),
		sx:    0,
		cx:    0,
		ss:    -1,
		text:  text,
		color: color.Black,
		style: style,
	}
	return
}

func (i *Input) Update(io *IO) {
	n := time.Now()
	d := n.Sub(i.timer)
	if d.Milliseconds() > 500 {
		i.flicker = !i.flicker
		i.timer = n
	}

	if io.mousePosIn(i.rect) {
		i.hover = true
	} else {
		i.hover = false
	}

	if io.mouse.button.left {
		if i.hover {
			io.focus(i.id)
			c := io.mouse.pos.x - i.rect.Min.X - i.style.def.Offset.X
			if z := i.check(c); z != -1 {
				i.cx = z
			}
			if !i.drag {
				i.ss = i.cx
			}
			i.drag = true
		}
	} else {
		i.drag = false
	}

	if io.state.focus == i.id {
		i.focus = true
	} else {
		i.focus = false
	}

	if i.focus {
		if ebiten.IsKeyPressed(ebiten.KeyShift) {
			if i.ss == -1 {
				i.ss = i.cx
			}
		}
		if inpututil.IsKeyJustReleased(ebiten.KeyRight) {
			i.cx++
			i.resetCursor()
		}
		if inpututil.IsKeyJustReleased(ebiten.KeyLeft) {
			i.cx--
			i.resetCursor()
		}
	}
}

func (i *Input) Draw(canvas *ebiten.Image, win *Window) {
	if i.focus {
		win.drawStyle(canvas, i.rect, 6, 1) // focus
	} else {
		win.drawStyle(canvas, i.rect, 6, 0) // normal
	}

	t := *i.text
	t = strings.ReplaceAll(t, " ", "_") // TODO hack... spaces are empty in measurements

	bb := text.BoundString(i.style.font, t[i.sx:i.cx])

	w := bb.Dx()

	y := i.rect.Max.Y - (i.rect.Dy()-win.peekPane().style.fontHeight)/2
	x := i.rect.Min.X + win.peekPane().style.def.Offset.X

	if i.ss != -1 && i.ss != i.cx {
		begin := i.ss
		end := i.cx
		if i.ss > i.cx {
			begin = i.cx
			end = i.ss
		}
		bb := text.BoundString(i.style.font, t[i.sx:begin])
		eb := text.BoundString(i.style.font, t[begin:end])

		r := image.Rect(x+bb.Max.X, i.rect.Min.Y, x+bb.Max.X+eb.Max.X, i.rect.Max.Y)

		win.drawStyle(canvas, r, 6, 2) // hover
	}

	text.Draw(canvas, *i.text, i.style.font, x, y, i.color)

	if i.flicker && i.focus {
		for t := 0; t < 2; t++ {
			ebitenutil.DrawLine(canvas, float64(x+w+t), float64(i.rect.Min.Y)+4, float64(x+w+t), float64(i.rect.Max.Y)-4, color.Black)
		}
	}

}

func (i *Input) Translate(x, y int) {
}

func (i *Input) resetCursor() {
	i.flicker = true
	i.timer = time.Now()
	if !ebiten.IsKeyPressed(ebiten.KeyShift) {
		i.ss = -1
	}
}

func (i *Input) check(nx int) (r int) {
	t := *i.text
	r = -1
	for j := len(t); j >= i.sx; j-- {
		b := text.BoundString(i.style.font, t[i.sx:j])
		if b.Max.X >= nx {
			r = j
		}
	}
	return
}
