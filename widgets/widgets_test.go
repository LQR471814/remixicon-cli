package widgets

import (
	"context"
	"fmt"
	"image"
	"log"
	"os"
	"testing"

	"github.com/mum4k/termdash"
	"github.com/mum4k/termdash/container"
	"github.com/mum4k/termdash/linestyle"
	"github.com/mum4k/termdash/terminal/tcell"
	"github.com/mum4k/termdash/terminal/terminalapi"
	"github.com/mum4k/termdash/widgetapi"
)

func testWidget(w widgetapi.Widget) error {
	f, err := os.Create("test.log")
	if err != nil {
		return err
	}
	defer f.Close()

	log.Default().SetOutput(f)
	log.Default().SetFlags(log.Lshortfile | log.Lmicroseconds)

	term, err := tcell.New()
	if err != nil {
		return err
	}
	defer term.Close()
	ctx, cancel := context.WithCancel(context.Background())
	handler := termdash.KeyboardSubscriber(func(k *terminalapi.Keyboard) {
		if k.Key == 'q' || k.Key == 'Q' {
			cancel()
		}
	})
	root, err := container.New(
		term,
		container.BorderTitle("Press Q to exit"),
		container.Border(linestyle.Round),
		container.PlaceWidget(w),
	)
	if err != nil {
		return err
	}
	return termdash.Run(ctx, term, root, handler)
}

func TestList(t *testing.T) {
	l := NewList()

	itemCount := 100
	items := make([]string, itemCount)
	for i := 0; i < itemCount; i++ {
		items[i] = fmt.Sprintf("row %d", i)
	}
	l.SetProps(func(state ListProps) ListProps {
		state.Rows = items
		return state
	})
	l.OnSelect = func(selected int) {
		l.SetProps(func(s ListProps) ListProps {
			s.Rows = append(
				s.Rows[:selected],
				append(
					[]string{s.Rows[selected]},
					s.Rows[selected:]...,
				)...,
			)
			return s
		})
	}

	err := testWidget(l)
	if err != nil {
		t.Error(err)
	}
}

func TestImage(t *testing.T) {
	f, err := os.Open("google.png")
	if err != nil {
		t.Error(err)
	}
	defer f.Close()

	img, _, err := image.Decode(f)
	if err != nil {
		t.Error(err)
		return
	}

	i := NewImage()
	i.SetProps(func(ip ImageProps) ImageProps {
		ip.Image = img
		return ip
	})
	err = testWidget(i)
	if err != nil {
		t.Error(err)
	}
}

func TestFit(t *testing.T) {
	w := &fitTest{Size: image.Pt(50, 35)}
	err := testWidget(w)
	if err != nil {
		t.Error(err)
	}
}
