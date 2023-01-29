package main

import (
	"flag"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/lithammer/fuzzysearch/fuzzy"
)

var iconVersion = flag.String("icon-version", "v2.5.0", "the version of RemixIcon to use")
var reset = flag.Bool("reset", false, "reset the cache")
var quiet = flag.Bool("quiet", false, "run quietly")

var output = flag.String("output", ".", "the directory to store the given icon")
var icon = flag.String("icon", "", "the icon to add, this takes no effect if --search is specified")
var outputFormat = flag.String("output-format", "svg", "the output format, supported formats: [svg, svelte]")
var search = flag.String("search", "", "search for an icon")
var maxResults = flag.Int("max-results", 10, "the maximum number of results")
var maxDistance = flag.Int("max-distance", 20, "the maximum levenshtein distance from a possible match")

func main() {
	log.SetFlags(log.Lmsgprefix)

	flag.Parse()

	if *output == "." {
		wd, err := os.Getwd()
		if err != nil {
			Fatal(err)
		}
		*output = wd
	}
	if *icon == "" && *search == "" {
		Fatal("you must provide an icon to add with --icon")
	}

	root, err := os.Executable()
	if err != nil {
		Fatal(err)
	}
	store := filepath.Join(root, "../", ".remixicon")
	Message("RemixIcon store is at %s\n", store)

	if *reset || !Exists(store) {
		Download(store, *iconVersion, Message, Fatal)
	}

	if *search != "" {
		ranks := fuzzy.RankFindNormalizedFold(*search, iconNames(store))
		i := 0
		for _, r := range ranks {
			if i == *maxResults {
				break
			}
			if r.Distance < *maxDistance {
				log.Printf("%d. %s", i+1, strings.Replace(r.Target, ".svg", "", 1))
				i++
			}
		}
		return
	}

	Message("icon destination %s", *output)
	formatMap := map[string]func(store, output, icon string){
		"svg":    svg,
		"svelte": svelte,
	}
	formatMap[*outputFormat](store, *output, *icon)
}
