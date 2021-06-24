package yagk

import (
	"github.com/hajimehoshi/ebiten/v2"
	"image"
)

type Pane struct {
	id      int
	style   *Style
	rect    image.Rectangle
	widgets []Widget
	parent  *Window
}

func NewPane(x, y, w int, style *Style) (p *Pane) {
	p = &Pane{
		rect:  image.Rect(x, y, w, 0),
		style: style,
	}
	return
}

func (p *Pane) Height() int {
	return p.rect.Dy()
}

func (p *Pane) Update(io *IO) {
	for i := len(p.widgets) - 1; i >= 0; i-- {
		p.widgets[i].Update(io)
	}
}

func (p *Pane) Draw(canvas *ebiten.Image, win *Window) {
	for _, w := range p.widgets {
		w.Draw(canvas, win)
	}
}

type ListItem struct {
	id            int
	Text, BtnText string
	IconId        int
}

func (p *Pane) List(listItems []ListItem, inView int, action func(index int)) *Pane {

	s := p.style.def.Sizes.Y + (p.style.def.Offset.Y * 2) // per item size

	sw := p.style.def.Sizes.X // scrollbar Width
	v := inView               // items in view

	if v > len(listItems)-1 {
		v = len(listItems) - 1
	}

	for i, _ := range listItems {
		listItems[i].id = i
	}

	x := p.rect.Min.X
	y := p.rect.Max.Y
	w := p.rect.Dx()
	h := (v * s) + s
	t := p.style

	// initial view of items
	type item struct {
		lbl *Label
		ctn *Container
		icn *Icon
	}
	var items []item

	for c := 0; c <= v; c++ {

		i := item{}

		i.ctn = newContainer(x, y+(c*s), w-sw, s, 0)

		i.icn = newIcon(listItems[c].IconId, x+t.def.Offset.X, y+(c*s)+p.style.def.Offset.Y, p.style.def)
		o := t.def.Offset.X + i.icn.rect.Dx()

		i.lbl = newLabel(listItems[c].Text, x+t.def.Offset.X+o, y+(c*s)+p.style.def.Offset.Y, w-(t.def.Offset.X*2)-sw, p.style.def.Sizes.Y)
		i.lbl.Color = p.style.def.Font.Color

		p.widgets = append(p.widgets, i.ctn, i.lbl, i.icn)
		items = append(items, i)
	}

	// v-scrollbar
	ctn := newContainer(x+w-sw, y, sw, h, 0)
	sld := newSlider(len(listItems)-v, x+w-sw, y, sw, h, func(pos int) {
		for i := 0; i < len(items); i++ {
			items[i].lbl.text = listItems[pos+i].Text
			items[i].icn.setId(listItems[pos+i].IconId)
			index := pos + i
			items[i].ctn.onclick = func() {
				action(index)
			}
		}
	})

	p.widgets = append(p.widgets, ctn, sld)

	p.rect.Max.Y += h

	return p
}

func (p *Pane) Button(label string, action func()) *Pane {
	x := p.rect.Min.X
	y := p.rect.Max.Y
	w := p.rect.Max.X
	h := p.style.def.Sizes.Y
	t := p.style

	ctn := newContainer(x, y, w, h+(t.def.Offset.Y*2), 0)

	btn := newContainer(x+t.def.Offset.X, y+t.def.Offset.Y, w-(t.def.Offset.X*2), h, 1)
	btn.onclick = action
	lbl := newLabel(label, x+t.def.Offset.X, y+t.def.Offset.Y, w-(t.def.Offset.X*2), h)
	lbl.Color = p.style.def.Font.Button
	lbl.Alignment = AlignCenter

	p.widgets = append(p.widgets, ctn, btn, lbl)

	p.rect.Max.Y += h
	p.rect.Max.Y += t.def.Offset.Y * 2

	return p
}

func (p *Pane) Confirm(label1 string, action1 func(), label2 string, action2 func()) *Pane {
	x := p.rect.Min.X
	y := p.rect.Max.Y
	w := p.rect.Dx()
	hw := w / 2
	h := p.style.def.Sizes.Y
	t := p.style

	ctn := newContainer(x, y, w, h+(t.def.Offset.Y*2), 0)

	btn1 := newContainer(x+t.def.Offset.X, y+t.def.Offset.Y, hw-(t.def.Offset.X*2), h, 1)
	btn1.onclick = action1
	lbl1 := newLabel(label1, x+t.def.Offset.X, y+t.def.Offset.Y, hw-(t.def.Offset.X*2), h)
	lbl1.Color = p.style.def.Font.Button
	lbl1.Alignment = AlignCenter

	btn2 := newContainer(x+t.def.Offset.X+hw, y+t.def.Offset.Y, hw-(t.def.Offset.X*2), h, 2)
	btn2.onclick = action2
	lbl2 := newLabel(label2, x+t.def.Offset.X+hw, y+t.def.Offset.Y, hw-(t.def.Offset.X*2), h)
	lbl2.Color = p.style.def.Font.Button
	lbl2.Alignment = AlignCenter

	p.widgets = append(p.widgets, ctn, btn1, lbl1, btn2, lbl2)

	p.rect.Max.Y += h
	p.rect.Max.Y += t.def.Offset.Y * 2

	return p
}

func (p *Pane) Selection(label string, selected *int, items ...string) *Pane {
	x := p.rect.Min.X
	y := p.rect.Max.Y
	w := p.rect.Dx()
	hw := w / 2
	h := p.style.def.Sizes.Y
	t := p.style

	ctn := newContainer(x, y, w, h+(t.def.Offset.Y*2), 0)

	lblSelect := newLabel(items[0], x+hw+(t.def.Offset.X*4), y+t.def.Offset.Y, hw-(t.def.Offset.X*8), h)
	lblSelect.Alignment = AlignCenter
	lblSelect.Color = p.style.def.Font.Color

	btnLeft := newContainer(x+hw+t.def.Offset.X, y+t.def.Offset.Y, t.def.Offset.X*2, h, 1)
	btnLeft.onclick = func() {
		*selected--
		if *selected < 0 {
			*selected = len(items) - 1
		}
		lblSelect.text = items[*selected]
	}
	z := "<"
	lblLeft := newLabel(z, x+hw+t.def.Offset.X, y+t.def.Offset.Y, t.def.Offset.X*2, h)
	lblLeft.Color = p.style.def.Font.Button
	lblLeft.Alignment = AlignCenter

	btnRight := newContainer(x+w-(t.def.Offset.X*3), y+t.def.Offset.Y, t.def.Offset.X*2, h, 1)
	btnRight.onclick = func() {
		*selected++
		if *selected >= len(items) {
			*selected = 0
		}
		lblSelect.text = items[*selected]
	}
	z = ">"
	lblRight := newLabel(z, x+w-(t.def.Offset.X*3), y+t.def.Offset.Y, t.def.Offset.X*2, h)
	lblRight.Color = p.style.def.Font.Button
	lblRight.Alignment = AlignCenter

	lbl := newLabel(label, x+t.def.Offset.X, y+t.def.Offset.Y, hw-t.def.Offset.X, h)
	lbl.Color = p.style.def.Font.Color

	p.widgets = append(p.widgets, ctn, lbl, lblSelect, btnLeft, btnRight, lblLeft, lblRight)

	p.rect.Max.Y += h
	p.rect.Max.Y += t.def.Offset.Y * 2

	return p
}

func (p *Pane) Text(text string) *Pane {
	x := p.rect.Min.X
	y := p.rect.Max.Y
	w := p.rect.Max.X
	h := p.style.def.Sizes.Y
	t := p.style

	ctn := newContainer(x, y, w, h+(t.def.Offset.Y*2), 0)

	lbl := newLabel(text, x+t.def.Offset.X, y+t.def.Offset.Y, w-(t.def.Offset.X*2), h)
	lbl.Alignment = AlignCenter
	lbl.Color = p.style.def.Font.Color

	p.widgets = append(p.widgets, ctn, lbl)

	p.rect.Max.Y += h
	p.rect.Max.Y += t.def.Offset.Y * 2

	return p
}

func (p *Pane) IconText(iconId int, text string) *Pane {
	x := p.rect.Min.X
	y := p.rect.Max.Y
	w := p.rect.Max.X
	h := p.style.def.Sizes.Y
	t := p.style

	ctn := newContainer(x, y, w, h+(t.def.Offset.Y*2), 0)

	icn := newIcon(iconId, x+t.def.Offset.X, y+t.def.Offset.Y, p.style.def)
	lbl := newLabel(text, x+(t.def.Offset.X*2)+icn.rect.Dx(), y+t.def.Offset.Y, w-(t.def.Offset.X*2), h)
	lbl.Color = p.style.def.Font.Color

	p.widgets = append(p.widgets, ctn, icn, lbl)

	p.rect.Max.Y += h
	p.rect.Max.Y += t.def.Offset.Y * 2

	return p
}

func (p *Pane) ProgressBar(total int, progress *int) *Pane {
	x := p.rect.Min.X
	y := p.rect.Max.Y
	w := p.rect.Max.X
	h := p.style.def.Sizes.Y
	t := p.style

	ctn := newContainer(x, y, w, h+(t.def.Offset.Y*2), 0)

	prg := newProgress(total, x+t.def.Offset.X, y+t.def.Offset.Y, w-(t.def.Offset.X*2), h)
	prg.progress = progress

	p.widgets = append(p.widgets, ctn, prg)

	p.rect.Max.Y += h
	p.rect.Max.Y += t.def.Offset.Y * 2
	return p
}

func (p *Pane) Checkbox(label string, checked *bool) *Pane {
	x := p.rect.Min.X
	y := p.rect.Max.Y
	w := p.rect.Max.X
	h := p.style.def.Sizes.Y
	t := p.style

	ctn := newContainer(x, y, w, h+(t.def.Offset.Y*2), 0)

	chk := newCheckbox(x+t.def.Offset.X, y+t.def.Offset.Y, h, h)
	chk.checked = checked

	lbl := newLabel(label, x+t.def.Offset.X+chk.rect.Dx(), y+t.def.Offset.Y, w-(t.def.Offset.X*2), h)
	lbl.Color = p.style.def.Font.Color

	p.widgets = append(p.widgets, ctn, chk, lbl)

	p.rect.Max.Y += h
	p.rect.Max.Y += t.def.Offset.Y * 2
	return p
}

func (p *Pane) Input(text *string) *Pane {
	x := p.rect.Min.X
	y := p.rect.Max.Y
	w := p.rect.Max.X
	h := p.style.def.Sizes.Y
	t := p.style

	ctn := newContainer(x, y, w, h+(t.def.Offset.Y*2), 0)
	ipt := newInput(text, x+t.def.Offset.X, y+t.def.Offset.Y, w-(t.def.Offset.X*2), h, p.style)
	p.widgets = append(p.widgets, ctn, ipt)

	p.rect.Max.Y += h
	p.rect.Max.Y += t.def.Offset.Y * 2
	return p
}
