package yagk

import (
	"github.com/hajimehoshi/ebiten/v2"
	"image"
	"math/rand"
)

type Slider struct {
	id              WidgetId
	rect            image.Rectangle
	orientation     Orientation
	total, selected int
	inView          *image.Rectangle
	hover, pressed  bool
	onchange        func(pos int)
}

func newSlider(total, x, y, w, h int, onchange func(pos int)) (s *Slider) {
	view := image.Rect(0, 0, w, h)
	s = &Slider{
		id:          WidgetId(rand.Int()),
		rect:        image.Rect(x, y, x+w, y+h),
		orientation: OrientVertical,
		inView:      &view,
		total:       total,
		selected:    0,
		onchange:    onchange,
	}
	return
}

func (s *Slider) calcRect(step int) (r image.Rectangle) {
	if s.orientation == OrientHorizontal {
		i := s.rect.Dx() / s.total
		r = image.Rect(s.rect.Min.X+(i*step), s.rect.Min.Y, s.rect.Min.X+i+(i*step), s.rect.Max.Y)
	} else {
		i := s.rect.Dy() / s.total
		r = image.Rect(s.rect.Min.X, s.rect.Min.Y+(i*step), s.rect.Max.X, s.rect.Min.Y+i+(i*step))
	}
	return
}

func (s *Slider) Update(io *IO) {
	if io.mousePosIn(s.calcRect(s.selected)) {
		s.hover = true
	} else {
		s.hover = false
	}

	if io.mouse.button.left {
		if s.hover && io.activate(s.id) {
			s.pressed = true
		}
		if s.pressed {
			o := s.selected
			for i := 0; i < s.total; i++ {
				if s.orientation == OrientVertical {
					if io.mousePosYIn(s.calcRect(i)) {
						s.selected = i
					}
				} else {
					if io.mousePosXIn(s.calcRect(i)) {
						s.selected = i
					}
				}
			}
			if s.selected != o && s.onchange != nil {
				s.onchange(s.selected)
			}
		}
	} else {
		io.deactivate(s.id)
		s.pressed = false
	}
}

func (s *Slider) Draw(canvas *ebiten.Image, win *Window) {
	if s.pressed {
		win.drawStyle(canvas, s.calcRect(s.selected), 4, 1)
	} else if s.hover {
		win.drawStyle(canvas, s.calcRect(s.selected), 4, 2)
	} else {
		win.drawStyle(canvas, s.calcRect(s.selected), 4, 0)
	}
}

func (s *Slider) Translate(x, y int) {
}
