package display

import (
	"github.com/llgcode/draw2d"
	"github.com/BurntSushi/xgbutil"
	"github.com/BurntSushi/xgbutil/keybind"
	"github.com/BurntSushi/xgbutil/xevent"
	"github.com/BurntSushi/xgbutil/xgraphics"
	"github.com/BurntSushi/xgbutil/xwindow"
	"image"
	"image/color"
)

type Display struct {
	x *xgbutil.XUtil
	gc   *draw2d.ImageGraphicContext
	ximg *xgraphics.Image
	wid  *xwindow.Window
	w float64
	h float64
	bord float64
	head float64
}

func NewDisplay(width, height, border, heading int, name string) (*Display, error) {
	d := new(Display)
	d.w = float64(width)
	d.h = float64(height)
	d.bord = float64(border)
	d.head = float64(heading)
	X, err := xgbutil.NewConn()
	if err != nil {
		return nil, err
	}
	keybind.Initialize(X)
	d.ximg = xgraphics.New(X, image.Rect(
		0,
		0,
		border * 2 + width,
		border * 2 + heading + height))
	err = d.ximg.CreatePixmap()
	if err != nil {
		return nil, err
	}
	painter := NewXimgPainter(d.ximg)
	d.gc = draw2d.NewGraphicContextWithPainter(d.ximg, painter)
	d.gc.Save()
	d.gc.SetStrokeColor(color.White)
	d.gc.SetFillColor(color.White)
	d.gc.Clear()
	d.wid = d.ximg.XShowExtra(name, true)
	d.x = X
	go func() {
		xevent.Main(X)
	}()
	return d, nil
}

func (d *Display) NewKeyBinding(f func(), key string) error {
	err := keybind.KeyReleaseFun(
		func(X *xgbutil.XUtil, e xevent.KeyReleaseEvent) {
			f()
	}).Connect(d.x, d.x.RootWin(), key, true)
	if err != nil {
		return err
	}
	return nil
}

func (d *Display) SetHeadingText(text string) {
	d.gc.SetStrokeColor(color.White)
	d.gc.SetFillColor(color.White)
	draw2d.Rect(d.gc, d.bord, d.bord, d.bord + d.w, d.bord + d.head)
	d.gc.FillStroke()
	d.gc.SetStrokeColor(color.Black)
	d.gc.SetFillColor(color.Black)
	d.gc.FillStringAt(text, d.bord * 2, d.bord + d.head - 1)
}

func (d *Display) Draw(x, y, r float64, c color.Color) {
	d.gc.SetStrokeColor(c)
	d.gc.SetFillColor(c)
	draw2d.Circle(d.gc, x + d.bord, d.head + d.bord + d.h - y, r)
	d.gc.FillStroke()
	d.ximg.XDraw()
}

func (d *Display) Frame() {
	d.ximg.XPaint(d.wid.Id)
}

func (d *Display) NewParticle(x, y, r float64, c color.Color) *Particle {
	p := new(Particle)
	p.disp = d
	p.x = x
	p.y = y
	p.r = r
	p.c = c
	p.disp.Draw(p.x, p.y, p.r, p.c)
	return p
}

type Particle struct {
	disp *Display
	x    float64
	y    float64
	r    float64
	c    color.Color
}

func (p *Particle) Move(x, y float64) {
	p.disp.Draw(p.x, p.y, p.r+1, color.White)
	p.x = x
	p.y = y
	p.disp.Draw(p.x, p.y, p.r, p.c)
}

func (p *Particle) ChangeColor(c color.Color) {
	p.c = c
}
