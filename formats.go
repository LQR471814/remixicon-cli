package main

import (
	"bytes"
	_ "embed"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/antchfx/htmlquery"
	"golang.org/x/net/html"
)

//go:embed svelte/Icon.svelte
var svelteIconTemplate string

//go:embed svelte/icon-context.ts
var svelteIconContext string

func svg(store, output, icon string) {
	src, err := os.Open(filepath.Join(store, icon+".svg"))
	if err != nil {
		Fatal(err)
	}
	dest, err := os.Create(filepath.Join(output, icon+".svg"))
	if err != nil {
		Fatal(err)
	}
	_, err = io.Copy(dest, src)
	if err != nil {
		Fatal(err)
	}
}

func svelte(store, output, icon string) {
	src, err := os.Open(filepath.Join(store, icon+".svg"))
	if err != nil {
		Fatal(err)
	}
	defer src.Close()

	doc, err := html.Parse(src)
	if err != nil {
		Fatal(err)
	}
	paths, err := htmlquery.QueryAll(doc, "//path")
	if err != nil {
		Fatal(err)
	}

	svgContent := ""
	for _, p := range paths {
		buff := bytes.NewBuffer([]byte{})
		err = html.Render(buff, p)
		if err != nil {
			Fatal(err)
		}
		svgContent += buff.String()
	}

	iconContext, err := os.Create(filepath.Join(output, "icon-context.ts"))
	if err != nil {
		Fatal(err)
	}
	defer iconContext.Close()
	iconContext.WriteString(svelteIconContext)

	dest, err := os.Create(filepath.Join(output, ToComponentName(icon)+".svelte"))
	if err != nil {
		Fatal(err)
	}
	defer dest.Close()
	dest.WriteString(strings.Replace(
		svelteIconTemplate, "[SVG_CONTENT]",
		svgContent, 1,
	))
}
