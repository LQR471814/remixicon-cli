package common

import (
	"log"
	"os"
	"path/filepath"
)

var RootFolder string

func init() {
	executableAt, err := os.Executable()
	if err != nil {
		log.Fatal(err)
	}
	RootFolder = filepath.Dir(executableAt)
}
