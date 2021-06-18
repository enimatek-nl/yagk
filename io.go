package yagk

import "image"

type IO struct {
	mouse Mouse
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
