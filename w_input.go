package yagk

import (
	"fmt"
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
	interval           time.Time
	flicker            bool
	hover, drag, focus bool
	style              *Style
	last               []ebiten.Key
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
			if z := i.checkPos(c); z != -1 {
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
		i.resetCursor()
	}

	if i.focus {
		pk := inpututil.PressedKeys()
		for _, key := range pk {
			if key == ebiten.KeyLeft && !i.checkKey(ebiten.KeyLeft) {
				i.checkSelection()
				if i.cx > 0 {
					i.cx--
				}
				i.resetCursor()
			} else if key == ebiten.KeyRight && !i.checkKey(ebiten.KeyRight) {
				i.checkSelection()
				i.incCur()
				i.resetCursor()
			} else if !i.checkKey(key) {
				t := *i.text

				if i.ss != -1 {
					begin, end := calcBE(i.ss, i.cx)
					t = fmt.Sprintf("%s%s", t[0:begin], t[end:])
					i.cx = begin
					i.ss = -1
				}

				k := mapKey(key)
				if k != "" {
					t = fmt.Sprintf("%s%s%s", t[0:i.cx], k, t[i.cx:])
				} else if key == ebiten.KeyBackspace || key == ebiten.KeyDelete {
					if i.cx > 0 {
						i.cx--
						t = fmt.Sprintf("%s%s", t[0:i.cx], t[i.cx+1:])
					}
				}

				*i.text = t

				if k != "" {
					i.resetCursor()
					i.incCur()
				}
			}
		}
		i.last = pk
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
		begin, end := calcBE(i.ss, i.cx)
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

func (i *Input) incCur() {
	if i.cx < len(*i.text) {
		i.cx++
	}
}

func calcBE(s, e int) (begin, end int) {
	begin = s
	end = e
	if s > e {
		begin = e
		end = s
	}
	return
}

func (i *Input) checkSelection() {
	if ebiten.IsKeyPressed(ebiten.KeyShift) {
		if i.ss == -1 {
			i.ss = i.cx
		}
	}
}

func (i *Input) resetCursor() {
	i.flicker = true
	i.timer = time.Now()
	if !ebiten.IsKeyPressed(ebiten.KeyShift) {
		i.ss = -1
	}
}

func (i *Input) checkPos(nx int) int {
	t := *i.text
	if nx < 0 {
		return 0
	}
	b := text.BoundString(i.style.font, t[i.sx:])
	if nx > b.Max.X {
		return len(t)
	}
	for j := i.sx; j <= len(t); j++ {
		b := text.BoundString(i.style.font, t[i.sx:j])
		xtr := 0
		if j < len(t)-1 {
			c := text.BoundString(i.style.font, t[j:j+1])
			xtr = c.Dx() / 2 // add half the next letter to give a better 'feel'
		}
		if b.Max.X+xtr > nx && b.Min.X < nx {
			return j
		}
	}
	return -1
}

func (i *Input) checkKey(key ebiten.Key) bool {
	for _, k := range i.last {
		if k == key {
			return true
		}
	}
	return false
}

func
mapKey(key ebiten.Key) string {
	shift := ebiten.IsKeyPressed(ebiten.KeyShift)
	switch key {
	case ebiten.KeyA:
		if shift {
			return "A"
		}
		return "a"
	case ebiten.KeyB:
		if shift {
			return "B"
		}
		return "b"
	case ebiten.KeyC:
		if shift {
			return "C"
		}
		return "c"
	case ebiten.KeyD:
		if shift {
			return "D"
		}
		return "d"
	case ebiten.KeyE:
		if shift {
			return "E"
		}
		return "e"
	case ebiten.KeyF:
		if shift {
			return "F"
		}
		return "f"
	case ebiten.KeyG:
		if shift {
			return "G"
		}
		return "g"
	case ebiten.KeyH:
		if shift {
			return "H"
		}
		return "h"
	case ebiten.KeyI:
		if shift {
			return "I"
		}
		return "i"
	case ebiten.KeyJ:
		if shift {
			return "J"
		}
		return "j"
	case ebiten.KeyK:
		if shift {
			return "K"
		}
		return "k"
	case ebiten.KeyL:
		if shift {
			return "L"
		}
		return "l"
	case ebiten.KeyM:
		if shift {
			return "M"
		}
		return "m"
	case ebiten.KeyN:
		if shift {
			return "N"
		}
		return "n"
	case ebiten.KeyO:
		if shift {
			return "O"
		}
		return "o"
	case ebiten.KeyP:
		if shift {
			return "P"
		}
		return "p"
	case ebiten.KeyQ:
		if shift {
			return "Q"
		}
		return "q"
	case ebiten.KeyR:
		if shift {
			return "R"
		}
		return "r"
	case ebiten.KeyS:
		if shift {
			return "S"
		}
		return "s"
	case ebiten.KeyT:
		if shift {
			return "T"
		}
		return "t"
	case ebiten.KeyU:
		if shift {
			return "U"
		}
		return "u"
	case ebiten.KeyV:
		if shift {
			return "V"
		}
		return "v"
	case ebiten.KeyW:
		if shift {
			return "W"
		}
		return "w"
	case ebiten.KeyX:
		if shift {
			return "X"
		}
		return "x"
	case ebiten.KeyY:
		if shift {
			return "Y"
		}
		return "y"
	case ebiten.KeyZ:
		if shift {
			return "Z"
		}
		return "z"
	case ebiten.KeyDigit0:
		return "0"
	case ebiten.KeyDigit1:
		return "1"
	case ebiten.KeyDigit2:
		return "2"
	case ebiten.KeyDigit3:
		return "3"
	case ebiten.KeyDigit4:
		return "4"
	case ebiten.KeyDigit5:
		return "5"
	case ebiten.KeyDigit6:
		return "6"
	case ebiten.KeyDigit7:
		return "7"
	case ebiten.KeyDigit8:
		return "8"
	case ebiten.KeyDigit9:
		return "9"
	case ebiten.KeySpace:
		return " "
	case ebiten.KeyComma:
		return ","
	case ebiten.KeyPeriod:
		return "."
	case ebiten.KeySemicolon:
		return ";"
	case ebiten.KeyEqual:
		if shift {
			return "+"
		}
		return "="
	case ebiten.KeyMinus:
		if shift {
			return "_"
		}
		return "-"
	case ebiten.KeyQuote:
		return "\""
	default:
		return ""
	}
}
