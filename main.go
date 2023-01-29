package main

import (
	"icon-cli/cmd"
	"icon-cli/common"
	"icon-cli/library"
	"log"
	"os"
)

func main() {
	f, err := os.Create(library.NewPath(
		common.RootFolder, "latest.log",
	).String())
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	log.Default().SetOutput(f)
	log.Default().SetFlags(log.Lshortfile | log.Ltime)

	cmd.Execute()
}
