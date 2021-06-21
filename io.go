package yagk

import "image"

type IO struct {
	state WidgetState
	mouse Mouse
}

type WidgetState struct {
	active WidgetId
	focus  WidgetId
}

type Mouse struct {
	pos    Pos
	button Button
}

type Button struct {
	left, right bool
}

type Pos struct {
	x, y int
}

func (io *IO) mousePosXIn(rect image.Rectangle) bool {
	return rect.Min.X <= io.mouse.pos.x && io.mouse.pos.x < rect.Max.X
}

func (io *IO) mousePosYIn(rect image.Rectangle) bool {
	return rect.Min.Y <= io.mouse.pos.y && io.mouse.pos.y < rect.Max.Y
}

func (io *IO) mousePosIn(rect image.Rectangle) bool {
	return io.mousePosXIn(rect) && io.mousePosYIn(rect)
}

func (io *IO) activate(c WidgetId) bool {
	if io.state.active == -1 || io.state.active == c {
		io.state.active = c
		return true
	}
	return false
}

func (io *IO) deactivate(c WidgetId) {
	if io.state.active == c {
		io.state.active = -1
	}
}
