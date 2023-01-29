package cmd

import (
	"icon-cli/common"
	"icon-cli/library"
	"log"
	"net/url"
	"time"

	"github.com/spf13/cobra"
)

const (
	SOURCE_HTTP   = "http"
	SOURCE_GITHUB = "github"
)

var iconLibrary *common.Store[library.Library]
var cfg *common.Store[Config]

func init() {
	rootCmd.AddCommand(updateCmd)
}

func Update(force bool) error {
	log.Println("checking for icon library updates...")

	cfg = common.NewStore(*configPath, Config{
		Source:   SOURCE_GITHUB,
		Location: "/Remix-Design/RemixIcon",
	})
	err := cfg.Load()
	if err != nil {
		return err
	}

	iconLibrary = common.NewStore(*libPath, library.Library{})
	err = iconLibrary.Load()
	if err != nil {
		return err
	}

	// * update every week
	if !force && time.Since(iconLibrary.Data.LastUpdate) < time.Hour*168 {
		return nil
	}

	var provider library.Provider
	switch cfg.Data.Source {
	case SOURCE_HTTP:
		parsed, err := url.Parse(cfg.Data.Location)
		if err != nil {
			return err
		}
		provider = &library.HTTP{
			Url: parsed,
		}
	case SOURCE_GITHUB:
		provider = &library.Github{
			RepoPath: cfg.Data.Location,
		}
	}
	provider.Latest()

	latest, err := provider.Latest()
	if err != nil {
		return err
	}
	if iconLibrary.Data.Version == latest {
		return nil
	}

	log.Printf("found new version %s, updating...", latest)
	data, err := provider.Pull(latest)
	if err != nil {
		return nil
	}

	iconLibrary.Data.Index = data
	iconLibrary.Data.Version = latest
	iconLibrary.Data.LastUpdate = time.Now()
	return iconLibrary.Write()
}

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "update the icon library",
	Run: func(cmd *cobra.Command, args []string) {
		err := Update(true)
		if err != nil {
			log.Fatal(err)
		}
	},
}
