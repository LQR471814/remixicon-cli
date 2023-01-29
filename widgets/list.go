package widgets

import (
	"fmt"
	"icon-cli/common"
	"image"
	"sync"

	"github.com/mum4k/termdash/cell"
	"github.com/mum4k/termdash/keyboard"
	"github.com/mum4k/termdash/mouse"
	"github.com/mum4k/termdash/private/area"
	"github.com/mum4k/termdash/private/canvas"
	"github.com/mum4k/termdash/private/draw"
	"github.com/mum4k/termdash/terminal/terminalapi"
	"github.com/mum4k/termdash/widgetapi"
)

func NumberPrefix(id int) string {
	return fmt.Sprintf(" %d. ", id)
}

func DotPrefix(_ int) string {
	return " • "
}

type ListProps struct {
	Hovered int
	Rows    []string
	// a number from 0-1, defines where to start scrolling the viewport
	ScrollMargin  float64
	KeyboardScope widgetapi.KeyScope
	MouseScope    widgetapi.MouseScope
}

type List struct {
	scroll *ScrollManager
	props  ListProps
	lock   sync.Mutex

	OnSelect func(int)
	OnHover  func(int)
	Prefix   func(int) string
}

func NewList(rows ...string) *List {
	return &List{
		scroll: NewScrollManager(true),
		props: ListProps{
			Rows:          rows,
			ScrollMargin:  0.5,
			KeyboardScope: widgetapi.KeyScopeFocused,
			MouseScope:    widgetapi.MouseScopeWidget,
		},
		Prefix: NumberPrefix,
		lock:   sync.Mutex{},
	}
}

func (l *List) Props() ListProps {
	return l.props
}

func (l *List) SetProps(transform func(ListProps) ListProps) {
	defer l.lock.Unlock()
	l.lock.Lock()
	l.props = transform(l.props)
	l.scroll.Update(func(ss ScrollState) ScrollState {
		ss.Rows = len(l.props.Rows)
		return ss
	})
}

func (l *List) scrollBy(delta int) {
	l.props.Hovered = common.Clamp(l.props.Hovered+delta, 0, len(l.props.Rows)-1)
	if l.OnHover != nil {
		l.OnHover(l.props.Hovered)
	}

	start, end := l.scroll.Visible()
	position := l.props.Hovered - start

	switch {
	case delta < 0:
		threshold := int((1 - l.props.ScrollMargin) * float64(end-start))
		if position <= threshold {
			l.scroll.ScrollBy(delta)
		}
	case delta > 0:
		threshold := int(l.props.ScrollMargin * float64(end-start))
		if position >= threshold {
			l.scroll.ScrollBy(delta)
		}
	}
}

func (l *List) Draw(cvs *canvas.Canvas, meta *widgetapi.Meta) error {
	needAr, err := area.FromSize(image.Pt(5, 1))
	if err != nil {
		return err
	}
	if !needAr.In(cvs.Area()) {
		return draw.ResizeNeeded(cvs)
	}

	l.scroll.Update(func(state ScrollState) ScrollState {
		state.Canvas = cvs
		return state
	})
	start, end := l.scroll.Visible()
	for i, row := range l.props.Rows[start:end] {
		id := i + start

		opt := []draw.TextOption{}
		if id == l.props.Hovered {
			opt = append(
				opt, draw.TextCellOpts(
					cell.FgColor(cell.ColorBlack),
					cell.BgColor(cell.ColorWhite),
				),
			)
		}

		var offset int
		if l.Prefix != nil {
			prefix := l.Prefix(id)
			offset = len(prefix)
			draw.Text(cvs, prefix, image.Pt(0, i))
		}

		err := draw.Text(cvs, row, image.Pt(offset, i), opt...)
		if err != nil {
			panic(err)
		}
	}

	return nil
}

func (l *List) Keyboard(k *terminalapi.Keyboard, meta *widgetapi.EventMeta) error {
	defer l.lock.Unlock()
	l.lock.Lock()

	switch k.Key {
	case 'j', keyboard.KeyArrowDown:
		l.scrollBy(1)
	case 'k', keyboard.KeyArrowUp:
		l.scrollBy(-1)
	case keyboard.KeyPgDn:
		l.scrollBy(16)
	case keyboard.KeyPgUp:
		l.scrollBy(-16)
	case keyboard.KeyEnter:
		if l.OnSelect != nil {
			l.OnSelect(l.props.Hovered)
		}
	}
	return nil
}

func (l *List) Mouse(m *terminalapi.Mouse, meta *widgetapi.EventMeta) error {
	defer l.lock.Unlock()
	l.lock.Lock()

	switch m.Button {
	case mouse.ButtonWheelUp:
		l.scrollBy(-1)
	case mouse.ButtonWheelDown:
		l.scrollBy(1)
	case mouse.ButtonLeft:
		if l.OnSelect != nil {
			l.OnSelect(l.props.Hovered)
		}
	}
	return nil
}

func (l *List) Options() widgetapi.Options {
	return widgetapi.Options{
		MinimumSize:  image.Pt(5, 1),
		WantKeyboard: l.props.KeyboardScope,
		WantMouse:    l.props.MouseScope,
	}
}
