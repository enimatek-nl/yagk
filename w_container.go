package yagk

import (
	"github.com/hajimehoshi/ebiten/v2"
	"image"
)

type Container struct {
	styleIndex            int
	rect                  image.Rectangle
	focus, hover, pressed bool
	onclick               func()
}

func newContainer(x, y, width, height, style int) (c *Container) {
	c = &Container{
		styleIndex: style,
		rect:       image.Rect(x, y, x+width, y+height),
	}
	return
}

func (c *Container) Update(io *IO) {
	if io.mousePosIn(c.rect) {
		c.hover = true
	} else {
		c.hover = false
	}

	if io.mouse.button.left {
		if c.hover {
			c.pressed = true
		}
	} else {
		if c.hover && c.pressed {
			if c.onclick != nil {
				c.onclick()
			}
		}
		c.pressed = false
	}
}

func (c *Container) Draw(canvas *ebiten.Image, win *Window) {
	if c.pressed && c.onclick != nil {
		win.drawStyle(canvas, c.rect, c.styleIndex, 1) // pressed
	} else if c.hover {
		win.drawStyle(canvas, c.rect, c.styleIndex, 2) // hover
	} else {
		win.drawStyle(canvas, c.rect, c.styleIndex, 0) // normal
	}
}

func (c *Container) Translate(x, y int) {
}
