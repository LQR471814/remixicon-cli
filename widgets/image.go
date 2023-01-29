package widgets

import (
	"bytes"
	"image"
	"strings"
	"sync"
	"unicode/utf8"

	"github.com/glibsm/dots"
	"github.com/mum4k/termdash/private/area"
	"github.com/mum4k/termdash/private/canvas"
	"github.com/mum4k/termdash/private/draw"
	"github.com/mum4k/termdash/terminal/terminalapi"
	"github.com/mum4k/termdash/widgetapi"
	"github.com/nfnt/resize"
)

type ImageProps struct {
	Image image.Image
	//express as width divided by height
	CharAspectRatio float64
	AlignX          AlignAt
	AlignY          AlignAt
}

type Image struct {
	lock  sync.Mutex
	props ImageProps

	canvas         *canvas.Canvas
	lastDimensions image.Rectangle
	rendered       string
	renderOffset   image.Point
}

func NewImage() *Image {
	return &Image{
		lock: sync.Mutex{},
		props: ImageProps{
			CharAspectRatio: 0.8,
			AlignX:          ALIGN_CENTER,
			AlignY:          ALIGN_CENTER,
			Image:           image.NewRGBA(image.Rectangle{}),
		},
	}
}

func (img *Image) Props() ImageProps {
	defer img.lock.Unlock()
	img.lock.Lock()
	return img.props
}

func (img *Image) SetProps(transform func(ImageProps) ImageProps) {
	defer img.lock.Unlock()
	img.lock.Lock()
	img.props = transform(img.props)

	if img.props.Image == nil {
		return
	}
	imageSize := img.props.Image.Bounds().Size()
	width := float64(imageSize.X) / img.props.CharAspectRatio
	img.props.Image = resize.Resize(
		uint(width), uint(imageSize.Y),
		img.props.Image, resize.Bicubic,
	)
	img.renderImage()
}

func (img *Image) renderImage() {
	buffer := bytes.NewBuffer(nil)
	dots.Write(
		img.props.Image, dots.Writer(buffer),
		dots.Width(img.bounds().X),
	)
	img.rendered = buffer.String()

	if img.canvas == nil {
		img.renderOffset = image.Point{}
		return
	}

	firstLine, _, _ := strings.Cut(img.rendered, "\n")
	img.renderOffset = AlignRectangle(
		img.canvas.Area(), image.Rectangle{
			Max: image.Point{
				X: utf8.RuneCountInString(firstLine),
				Y: strings.Count(img.rendered, "\n"),
			},
		},
		img.props.AlignX, img.props.AlignY,
	).Min
}

func (img *Image) Bounds() image.Point {
	defer img.lock.Unlock()
	img.lock.Lock()
	return img.bounds()
}

func (img *Image) bounds() image.Point {
	if img.props.Image == nil {
		return image.Point{}
	}
	imageSize := img.props.Image.Bounds()
	if img.canvas == nil {
		return image.Point{}
	}
	// rendering space is actually twice the height since a
	// single braille character is 2x4
	renderingSpace := img.canvas.Area()
	renderingSpace.Max.Y = renderingSpace.Max.Y * 2
	return FitRectangle(
		renderingSpace, imageSize,
		FIT_CONTAIN,
	).Size()
}

func (img *Image) Draw(cvs *canvas.Canvas, meta *widgetapi.Meta) error {
	img.lock.Lock()
	defer img.lock.Unlock()

	img.canvas = cvs
	needAr, err := area.FromSize(image.Pt(1, 1))
	if err != nil {
		return err
	}
	if !needAr.In(cvs.Area()) {
		return draw.ResizeNeeded(cvs)
	}

	if !cvs.Area().Eq(img.lastDimensions) && img.props.Image != nil {
		img.renderImage()
	}

	y := 0
	x := 0
	for _, r := range img.rendered {
		if r == '\n' {
			y++
			x = 0
			continue
		}

		cvs.SetCell(image.Pt(x, y).Add(img.renderOffset), r)
		x++
	}
	img.lastDimensions = cvs.Area()

	// draw.Text(
	// 	cvs,
	// 	fmt.Sprintf(
	// 		"%v[%v] = %v rendered: (%d, %d)",
	// 		cvs.Area().Size(),
	// 		img.props.Image.Bounds().Size(),
	// 		img.bounds(),
	// 		utf8.RuneCountInString(firstLine), y,
	// 	),
	// 	image.Pt(0, 0),
	// )

	return nil
}

func (img *Image) Keyboard(*terminalapi.Keyboard, *widgetapi.EventMeta) error {
	return nil
}

func (img *Image) Mouse(*terminalapi.Mouse, *widgetapi.EventMeta) error {
	return nil
}

func (img *Image) Options() widgetapi.Options {
	return widgetapi.Options{
		MinimumSize: image.Pt(1, 1),
	}
}
