package main

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func Fatal(s any) {
	if !*quiet {
		log.Fatal(s)
	}
}

func Message(s string, args ...any) {
	if !*quiet {
		log.Printf(s, args...)
	}
}

func Exists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func location(version string) string {
	return fmt.Sprintf(
		"https://github.com/Remix-Design/RemixIcon/releases/download/%s/RemixIcon_SVG_%s.zip",
		version, version,
	)
}

func Download(store, version string, Message func(string, ...any), Fatal func(any)) {
	url := location(version)
	Message("downloading icon files from %s", url)

	client := &http.Client{}
	response, err := client.Get(url)
	if err != nil {
		Fatal(err)
	}
	defer response.Body.Close()

	buffer := bytes.NewBuffer([]byte{})
	s, err := io.Copy(buffer, response.Body)
	if err != nil {
		Fatal(err)
	}
	reader := bytes.NewReader(buffer.Bytes())

	unzip, err := zip.NewReader(reader, s)
	if err != nil {
		Fatal(err)
	}

	err = os.Mkdir(store, 0777)
	if err != nil {
		Fatal(err)
	}
	for _, f := range unzip.File {
		if !f.FileInfo().IsDir() && strings.HasPrefix(f.Name, "icons") {
			path := filepath.Join(store, filepath.Base(f.Name))

			src, err := f.Open()
			if err != nil {
				Fatal(err)
			}
			defer src.Close()
			dest, err := os.Create(path)
			if err != nil {
				Fatal(err)
			}
			defer dest.Close()

			_, err = io.Copy(dest, src)
			if err != nil {
				Fatal(err)
			}
		}
	}
}

func ToComponentName(name string) string {
	segments := []string{}
	for _, s := range strings.Split(name, "-") {
		segments = append(segments, strings.ToUpper(s[0:1])+s[1:])
	}
	return strings.Join(segments, "")
}

func SwapExtension(filename string, newExt string) string {
	segments := strings.Split(filename, ".")
	return strings.Join(append(segments[:len(segments)-1], newExt), ".")
}
