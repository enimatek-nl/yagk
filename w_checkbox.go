package yagk

import (
	"github.com/hajimehoshi/ebiten/v2"
	"image"
	"math/rand"
)

type Checkbox struct {
	id             WidgetId
	rect           image.Rectangle
	checked        *bool
	hover, pressed bool
}

func newCheckbox(x, y, w, h int) (c *Checkbox) {
	c = &Checkbox{
		id:   WidgetId(rand.Int()),
		rect: image.Rect(x, y, x+w, y+h),
	}
	return
}

func (c *Checkbox) Update(io *IO) {
	if io.mousePosIn(c.rect) {
		c.hover = true
	} else {
		c.hover = false
	}

	if io.mouse.button.left {
		if c.hover {
			if io.activate(c.id) {
				c.pressed = true
			}
		}
	} else {
		if c.hover && c.pressed {
			*c.checked = !*c.checked
		}
		io.deactivate(c.id)
		c.pressed = false
	}
}

func (c *Checkbox) Draw(canvas *ebiten.Image, win *Window) {
	if c.hover {
		win.drawStyle(canvas, c.rect, 5, 2) // hover
	} else {
		win.drawStyle(canvas, c.rect, 5, 0) // normal
	}

	if c.checked != nil && *c.checked {
		win.drawStyle(canvas, c.rect, 5, 1) // checked
	}
}

func (c *Checkbox) Translate(x, y int) {
}
