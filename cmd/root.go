package cmd

import (
	"bytes"
	"context"
	"icon-cli/common"
	"icon-cli/library"
	"icon-cli/widgets"
	"image"
	"log"
	"strings"

	_ "image/jpeg"

	"github.com/mum4k/termdash"
	"github.com/mum4k/termdash/container"
	"github.com/mum4k/termdash/linestyle"
	"github.com/mum4k/termdash/terminal/tcell"
	"github.com/mum4k/termdash/terminal/terminalapi"
	"github.com/mum4k/termdash/widgetapi"
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

		iconIndexIds := make([]string, len(iconLibrary.Data.Index))
		list := widgets.NewList()
		list.SetProps(func(p widgets.ListProps) widgets.ListProps {
			p.KeyboardScope = widgetapi.KeyScopeGlobal
			p.MouseScope = widgetapi.MouseScopeGlobal
			i := 0
			for k := range iconLibrary.Data.Index {
				name := strings.Join(strings.Split(k, "-"), " ")
				p.Rows = append(p.Rows, name)
				iconIndexIds[i] = k
				i++
			}
			return p
		})

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

		root, err := container.New(
			term,
			container.Border(linestyle.None),
			container.SplitVertical(
				container.Left(
					container.Border(linestyle.Round),
					container.PlaceWidget(img),
				),
				container.Right(
					container.Border(linestyle.Round),
					container.PlaceWidget(list),
				),
			),
		)
		if err != nil {
			log.Fatal(err)
		}

		ctx, cancel := context.WithCancel(context.Background())
		handler := termdash.KeyboardSubscriber(func(k *terminalapi.Keyboard) {
			if k.Key == 'q' || k.Key == 'Q' {
				cancel()
			}
		})
		err = termdash.Run(ctx, term, root, handler)
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
