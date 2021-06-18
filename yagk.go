package yagk

import (
	"github.com/hajimehoshi/ebiten/v2"
)

/*
Window -> Panes <- ⬇️Widgets
*/


type Alignment int
type Orientation int

const (
	AlignLeft   Alignment = 0
	AlignCenter Alignment = 1
	AlignRight  Alignment = 2

	OrientHorizontal Orientation = 0
	OrientVertical   Orientation = 1
)

type Widget interface {
	Update(io *IO)
	Draw(canvas *ebiten.Image, win *Window)
	Translate(x, y int)
}
