package library

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"io/ioutil"
	"log"
	"net/url"
	"sync"
	"time"

	"icon-cli/common"

	"github.com/carlmjohnson/requests"
	"github.com/iancoleman/strcase"
	"golang.org/x/net/html"
)

type Library struct {
	Index      map[KebabCase][]byte
	Version    string
	LastUpdate time.Time
}

// a string in KebabCase
type KebabCase = string

type Provider interface {
	Latest() (string, error)
	Pull(string) (map[KebabCase][]byte, error)
}

type HTTP struct {
	Url    *url.URL
	Buffer *bytes.Buffer
	lock   sync.Mutex
}

func (h *HTTP) Latest() (string, error) {
	defer h.lock.Unlock()
	h.lock.Lock()

	h.Buffer = bytes.NewBuffer(nil)
	err := requests.
		URL(h.Url.String()).
		ToBytesBuffer(h.Buffer).
		Fetch(context.Background())
	if err != nil {
		return "", err
	}
	hasher := sha256.New()
	hasher.Write(h.Buffer.Bytes())

	return hex.EncodeToString(hasher.Sum(nil)), nil
}

func (h *HTTP) Pull(_ string) (map[KebabCase][]byte, error) {
	return GenerateFromTarGz(h.Url)
}

type Github struct {
	RepoPath string
}

func (g Github) Latest() (string, error) {
	root := &html.Node{}

	builder := requests.
		URL((&url.URL{
			Scheme: "https",
			Host:   "github.com",
			Path:   NewPath(g.RepoPath, "/releases/latest").String(),
		}).String()).
		Handle(requests.ToHTML(root))

	err := builder.Fetch(context.Background())
	if err != nil {
		return "", err
	}

	request, err := builder.Request(context.Background())
	if err != nil {
		return "", err
	}

	redirected, err := common.CheckRedirect(request)
	if err != nil {
		return "", err
	}
	parsed, err := url.Parse(redirected)
	if err != nil {
		return "", err
	}
	splitPath := NewPath(parsed.Path)

	return splitPath[len(splitPath)-1], nil
}

func (g Github) Pull(tag string) (map[KebabCase][]byte, error) {
	return GenerateFromTarGz(&url.URL{
		Scheme: "https",
		Host:   "github.com",
		Path:   NewPath(g.RepoPath, "/archive/refs/tags", tag+".tar.gz").String(),
	})
}

func GenerateFromTarGz(loc *url.URL) (map[KebabCase][]byte, error) {
	buffer := bytes.NewBuffer(nil)

	err := requests.
		URL(loc.String()).
		ToBytesBuffer(buffer).
		Fetch(context.Background())
	if err != nil {
		return nil, err
	}

	gzipReader, err := gzip.NewReader(buffer)
	if err != nil {
		return nil, err
	}
	uncompressed, err := ioutil.ReadAll(gzipReader)
	if err != nil {
		return nil, err
	}

	lib := map[KebabCase][]byte{}

	tarReader := tar.NewReader(bytes.NewBuffer(uncompressed))
	for {
		header, err := tarReader.Next()
		if err != nil {
			break
		}

		switch header.Typeflag {
		case tar.TypeReg:
			buffer := bytes.NewBuffer(nil)
			_, err := io.Copy(buffer, tarReader)
			if err != nil {
				log.Println(err)
				continue
			}

			name := NewPath(header.Name).Basename()
			name, _ = SplitExtension(name)
			name = strcase.ToKebab(name)
			if name == "" {
				break
			}

			lib[name] = buffer.Bytes()
		}
	}

	return lib, nil
}
