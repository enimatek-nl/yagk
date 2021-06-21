package yagk

import (
	"github.com/hajimehoshi/ebiten/v2"
	"image"
	"log"
)

type Window struct {
	Width, Height int
	title         string
	panes         []*Pane
	io            *IO
}

func New(title string, width, height int) (win *Window) {
	win = &Window{
		Width:  width,
		Height: height,
		title:  title,
		io: &IO{
			state: WidgetState{
				active: -1,
				focus:  -1,
			},
		},
	}

	return
}

func (win *Window) PushPane(p *Pane) {
	p.parent = win
	p.id = len(win.panes)
	win.panes = append(win.panes, p)
}

func (win *Window) PopPane() bool {
	if len(win.panes) > 1 {
		win.panes = win.panes[0 : len(win.panes)-1]
		return true
	} else {
		return false
	}
}

func (win *Window) GetPane(i int) (p *Pane) {
	if i > len(win.panes)-1 {
		return nil
	}
	return win.panes[i]
}

func (win *Window) drawNinePatches(dst *ebiten.Image, dstRect image.Rectangle, srcRect image.Rectangle) {
	srcX := srcRect.Min.X
	srcY := srcRect.Min.Y
	srcW := srcRect.Dx()
	srcH := srcRect.Dy()

	dstX := dstRect.Min.X
	dstY := dstRect.Min.Y
	dstW := dstRect.Dx()
	dstH := dstRect.Dy()

	op := &ebiten.DrawImageOptions{}
	for j := 0; j < 3; j++ {
		for i := 0; i < 3; i++ {
			op.GeoM.Reset()

			sx := srcX
			sy := srcY
			sw := srcW / 4
			sh := srcH / 4
			dx := 0
			dy := 0
			dw := sw
			dh := sh
			switch i {
			case 1:
				sx = srcX + srcW/4
				sw = srcW / 2
				dx = srcW / 4
				dw = dstW - 2*srcW/4
			case 2:
				sx = srcX + 3*srcW/4
				dx = dstW - srcW/4
			}
			switch j {
			case 1:
				sy = srcY + srcH/4
				sh = srcH / 2
				dy = srcH / 4
				dh = dstH - 2*srcH/4
			case 2:
				sy = srcY + 3*srcH/4
				dy = dstH - srcH/4
			}

			op.GeoM.Scale(float64(dw)/float64(sw), float64(dh)/float64(sh))
			op.GeoM.Translate(float64(dx), float64(dy))
			op.GeoM.Translate(float64(dstX), float64(dstY))
			dst.DrawImage(win.peekPane().style.data.SubImage(image.Rect(sx, sy, sx+sw, sy+sh)).(*ebiten.Image), op)
		}
	}
}

func (win *Window) drawStyle(dst *ebiten.Image, rect image.Rectangle, row, index int) {
	y0 := row * win.peekPane().style.def.Size
	x0 := index * win.peekPane().style.def.Size
	win.drawNinePatches(dst, rect, image.Rect(x0, y0, x0+win.peekPane().style.def.Size, y0+win.peekPane().style.def.Size))
}

func (win *Window) Run() {
	if len(win.panes) > 0 {
		ebiten.SetWindowSize(win.Width, win.Height)
		ebiten.SetWindowTitle(win.title)
		if err := ebiten.RunGame(win); err != nil {
			log.Fatal(err)
		}
	} else {
		log.Fatal("use NewPane and win.PushPane to have at least 1 pane ready before win.Run")
	}
}

func (win *Window) Update() error {
	mpx, mpy := ebiten.CursorPosition()
	win.io.mouse = Mouse{
		pos: Pos{
			x: mpx,
			y: mpy,
		},
		button: Button{
			left:  ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft),
			right: ebiten.IsMouseButtonPressed(ebiten.MouseButtonRight),
		},
	}
	//for _, p := range win.panes {
	//	p.Update(io)
	//}
	win.peekPane().Update(win.io)
	return nil
}

func (win *Window) Draw(screen *ebiten.Image) {
	screen.Fill(
		win.peekPane().style.def.Background)
	win.peekPane().Draw(screen, win)
}

func (win *Window) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return win.Width, win.Height
}

func (win *Window) peekPane() *Pane {
	return win.panes[len(win.panes)-1]
}
