package widgets

import (
	"image"
	"math"
	"sync"

	"github.com/mum4k/termdash/private/canvas"
)

type ScrollState struct {
	Rows      int
	ScrollTop int
	Canvas    *canvas.Canvas
}

type ScrollManager struct {
	state        ScrollState
	lock         sync.Mutex
	ScrollMargin bool
}

func NewScrollManager(margin bool) *ScrollManager {
	return &ScrollManager{
		lock:         sync.Mutex{},
		ScrollMargin: margin,
	}
}

func (m *ScrollManager) State() ScrollState {
	return m.state
}

func (m *ScrollManager) Update(transform func(ScrollState) ScrollState) {
	defer m.lock.Unlock()
	m.lock.Lock()
	m.state = transform(m.state)
	h := 0
	if m.state.Canvas != nil {
		h = m.state.Canvas.Size().Y
	}
	if m.state.ScrollTop+h >= m.state.Rows {
		m.state.ScrollTop = m.state.Rows - m.state.Canvas.Size().Y
	} else if m.state.ScrollTop < 0 {
		m.state.ScrollTop = 0
	}
}

func (m *ScrollManager) ScrollBy(delta int) {
	m.Update(func(state ScrollState) ScrollState {
		state.ScrollTop += delta
		return state
	})
}

func (m *ScrollManager) Visible() (int, int) {
	start := m.state.ScrollTop
	end := m.state.ScrollTop + m.state.Canvas.Size().Y
	return start, end
}

func Distance(pt image.Point) float64 {
	return math.Sqrt(math.Pow(float64(pt.X), 2) + math.Pow(float64(pt.Y), 2))
}

type FitOperation = int

const (
	FIT_CONTAIN FitOperation = iota
	FIT_COVER
)

func FitRectangle(parent, child image.Rectangle, op FitOperation) image.Rectangle {
	height := int(float64(parent.Dx()) * float64(child.Dy()) / float64(child.Dx()))
	dimensions := image.Point{
		X: int(float64(parent.Dy()) * float64(child.Dx()) / float64(child.Dy())),
		Y: parent.Dy(),
	}
	if (op == FIT_CONTAIN && height <= parent.Dy()) ||
		(op == FIT_COVER && height >= parent.Dy()) {
		dimensions = image.Point{
			X: parent.Dx(),
			Y: height,
		}
	}
	return image.Rectangle{
		Min: image.Point{},
		Max: dimensions,
	}
}

type AlignAt = int

const (
	ALIGN_TOP AlignAt = iota
	ALIGN_CENTER
	ALIGN_BOTTOM
)

func AlignRectangle(parent, child image.Rectangle, x, y AlignAt) image.Rectangle {
	var offsetX, offsetY float64
	if x == ALIGN_CENTER {
		offsetX = 0.5
	}
	if x == ALIGN_BOTTOM {
		offsetX = 1
	}
	if y == ALIGN_CENTER {
		offsetY = 0.5
	}
	if y == ALIGN_BOTTOM {
		offsetY = 1
	}
	min := image.Pt(
		int(offsetX*float64(parent.Dx()))-int(offsetX*float64(child.Dx())),
		int(offsetY*float64(parent.Dy()))-int(offsetY*float64(child.Dy())),
	)
	return image.Rectangle{
		Min: min,
		Max: min.Add(child.Size()),
	}
}
