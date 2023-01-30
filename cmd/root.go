package cmd

import (
	"bytes"
	"context"
	"icon-cli/common"
	"icon-cli/library"
	"icon-cli/widgets"
	"image"
	"log"

	_ "image/jpeg"

	"github.com/lithammer/fuzzysearch/fuzzy"
	"github.com/mum4k/termdash"
	"github.com/mum4k/termdash/cell"
	"github.com/mum4k/termdash/container"
	"github.com/mum4k/termdash/keyboard"
	"github.com/mum4k/termdash/linestyle"
	"github.com/mum4k/termdash/terminal/tcell"
	"github.com/mum4k/termdash/terminal/terminalapi"
	"github.com/mum4k/termdash/widgetapi"
	"github.com/mum4k/termdash/widgets/textinput"
	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
	"github.com/srwiley/oksvg"
	"github.com/srwiley/rasterx"
)

var imageRes = 200

func renderSVG(id string) image.Image {
	buffer := bytes.NewBuffer(
		iconLibrary.Data.Index[id],
	)
	icon, err := oksvg.ReadIconStream(buffer)
	if err != nil {
		log.Println(err)
		return nil
	}
	icon.SetTarget(0, 0, float64(imageRes), float64(imageRes))
	rgba := image.NewRGBA(image.Rect(0, 0, imageRes, imageRes))
	dasher := rasterx.NewDasher(
		imageRes, imageRes,
		rasterx.NewScannerGV(
			imageRes, imageRes, rgba, rgba.Bounds(),
		),
	)
	icon.Draw(dasher, 1)
	return rgba
}

const (
	FOCUS_SEARCH container.FocusGroup = iota
	FOCUS_LIST
)

var rootCmd = &cobra.Command{
	Use:   "icon",
	Short: "an SVG icon cli",
	Long:  "a framework/toolchain agnostic method of incorporating SVG icons into projects.",
	Run: func(cmd *cobra.Command, args []string) {
		err := Update(false)
		if err != nil {
			log.Fatal(err)
		}

		term, err := tcell.New()
		if err != nil {
			log.Fatal(err)
		}
		defer term.Close()

		img := widgets.NewImage()
		if err != nil {
			log.Fatal(err)
		}

		list := widgets.NewList()
		list.SetProps(func(lp widgets.ListProps) widgets.ListProps {
			lp.KeyboardScope = widgetapi.KeyScopeGlobal
			lp.MouseScope = widgetapi.MouseScopeGlobal
			return lp
		})
		list.Prefix = widgets.NumberPrefix

		indexIds := make([]string, len(iconLibrary.Data.Index))
		i := 0
		for k := range iconLibrary.Data.Index {
			indexIds[i] = k
			i++
		}

		var iconIndexIds []string

		// * if nil, will use the full index
		updateListItems := func(items []string) {
			iconIndexIds = items
			list.SetProps(func(lp widgets.ListProps) widgets.ListProps {
				lp.Rows = items
				return lp
			})
		}

		updateListItems(indexIds)

		img.SetProps(func(ip widgets.ImageProps) widgets.ImageProps {
			ip.Image = renderSVG(iconIndexIds[0])
			return ip
		})

		list.OnHover = func(i int) {
			rendered := renderSVG(iconIndexIds[i])
			img.SetProps(func(ip widgets.ImageProps) widgets.ImageProps {
				ip.Image = rendered
				return ip
			})
		}

		input, err := textinput.New(
			textinput.FillColor(cell.ColorBlack),
			textinput.PlaceHolder("Search"),
			textinput.PlaceHolderColor(cell.ColorLime),
			textinput.OnChange(func(data string) {
				if data == "" {
					updateListItems(indexIds)
					return
				}
				ranks := fuzzy.RankFindFold(data, indexIds)
				if ranks.Len() > 0 {
					results := make([]string, ranks.Len())
					for i, r := range ranks {
						results[i] = r.Target
					}
					updateListItems(results)
				}
			}),
		)
		if err != nil {
			log.Fatal(err)
		}

		root, err := container.New(
			term,
			// container.KeyFocusNext(keyboard.KeyTab),
			container.Border(linestyle.None),
			container.SplitVertical(
				container.Left(
					container.Border(linestyle.Round),
					container.KeyFocusSkip(),
					container.PlaceWidget(img),
				),
				container.Right(
					container.Border(linestyle.Round),
					container.SplitHorizontal(
						container.Top(
							container.PaddingLeft(1),
							container.PaddingRight(1),
							container.Focused(),
							container.PlaceWidget(input),
						),
						container.Bottom(
							container.PlaceWidget(list),
						),
						container.SplitFixed(1),
					),
				),
			),
		)
		if err != nil {
			log.Fatal(err)
		}

		ctx, cancel := context.WithCancel(context.Background())
		handler := termdash.KeyboardSubscriber(func(k *terminalapi.Keyboard) {
			switch k.Key {
			case keyboard.KeyCtrlX:
				cancel()
			case keyboard.KeyEsc:
				input.ReadAndClear()
			}
		})
		errorHandler := termdash.ErrorHandler(func(err error) {
			log.Println(err)
		})
		err = termdash.Run(ctx, term, root, handler, errorHandler)
		if err != nil {
			log.Fatal(err)
		}
	},
}

// var verbose *bool
var libPath *string
var configPath *string

func GenerateDocs(dir string) error {
	return doc.GenMarkdownTree(rootCmd, dir)
}

func Execute() {
	rootCmd.Execute()
}

func init() {
	// verbose = rootCmd.PersistentFlags().BoolP("verbose", "v", false, "verbose mode")
	libPath = rootCmd.PersistentFlags().StringP(
		"library", "l", library.NewPath(common.RootFolder, "icons.bin").String(),
		"specify where the library should be stored",
	)
	configPath = rootCmd.PersistentFlags().StringP(
		"config", "c", library.NewPath(common.RootFolder, "config.bin").String(),
		"specify where the config should be stored",
	)
}
