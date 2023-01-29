package library

import (
	"fmt"
	"log"
	"testing"
)

func TestGithub(t *testing.T) {
	github := Github{
		RepoPath: "/Remix-Design/RemixIcon",
	}

	var latest string
	t.Run("latest", func(t *testing.T) {
		var err error
		latest, err = github.Latest()
		if err != nil {
			t.Error(err)
			return
		}
		log.Println(latest)
	})

	t.Run("pull", func(t *testing.T) {
		pulled, err := github.Pull(latest)
		if err != nil {
			t.Error(err)
			return
		}
		i := 0
		for k := range pulled {
			if i == 10 {
				break
			}
			fmt.Println(k)
			i++
		}
	})
}
