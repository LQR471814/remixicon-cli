package main

import "io/ioutil"

func iconNames(store string) []string {
	files, err := ioutil.ReadDir(store)
	if err != nil {
		Fatal(err)
	}
	list := make([]string, len(files))
	for i, f := range files {
		list[i] = SwapExtension(f.Name(), "")
	}
	return list
}
