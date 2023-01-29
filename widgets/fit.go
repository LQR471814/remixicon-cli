package widgets

import (
	"image"

	"github.com/mum4k/termdash/private/canvas"
	"github.com/mum4k/termdash/private/draw"
	"github.com/mum4k/termdash/terminal/terminalapi"
	"github.com/mum4k/termdash/widgetapi"
)

type fitTest struct {
	Size image.Point
}

func (f *fitTest) Draw(cvs *canvas.Canvas, meta *widgetapi.Meta) error {
	rectangle := FitRectangle(cvs.Area(), image.Rectangle{Max: f.Size}, FIT_CONTAIN)
	rectangle = AlignRectangle(
		cvs.Area(), rectangle,
		ALIGN_CENTER, ALIGN_CENTER,
	)
	return draw.Border(cvs, rectangle)
}

func (f *fitTest) Keyboard(*terminalapi.Keyboard, *widgetapi.EventMeta) error {
	return nil
}

func (f *fitTest) Mouse(*terminalapi.Mouse, *widgetapi.EventMeta) error {
	return nil
}

func (f *fitTest) Options() widgetapi.Options {
	return widgetapi.Options{}
}
