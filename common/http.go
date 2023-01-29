package common

import (
	"errors"
	"log"
	"net/http"
	"strings"
	"sync"

	normalizeurl "github.com/sekimura/go-normalize-url"
)

var redirectTrap *http.Client
var redirectedUrls sync.Map

func init() {
	redirectTrap = &http.Client{}
	redirectTrap.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		redirectedTo := req.URL.String()
		for _, r := range via {
			normalized, err := normalizeurl.Normalize(r.URL.String())
			if err != nil {
				log.Println(err)
				continue
			}
			redirectedUrls.Store(normalized, redirectedTo)
		}
		return errors.New("redirect trapped")
	}
}

func CheckRedirect(req *http.Request) (string, error) {
	_, err := redirectTrap.Do(req)
	if err != nil && !strings.Contains(err.Error(), "redirect trapped") {
		return "", err
	}
	origin := req.URL.String()
	loaded, has := redirectedUrls.Load(origin)
	if !has {
		return origin, errors.New("there was no redirect")
	}
	return loaded.(string), nil
}
