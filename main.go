package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"log/slog"

	"github.com/google/go-github/v54/github"
	"github.com/komisan19/gm/config"
)

const (
	name    = "gm"
	version = "0.0.1"
)

type GistDesciption struct {
	IsPublic    bool
	Description string
	Context     []byte
}

func healthCheck(client *github.Client) {
	ctx := context.Background()
	_, _, err := client.Repositories.List(ctx, "", nil)
	if _, ok := err.(*github.RateLimitError); ok {
		slog.Warn("Limit Over", "Mes", "hit rate limit")
		return
	}
}

func createGist(ctx context.Context, client *github.Client, gd *GistDesciption, uploadFile string) error {

	context, err := os.ReadFile(uploadFile)
	if err != nil {
		return err
	}

	fileName := filepath.Base(uploadFile)

	gist := &github.Gist{
		Description: github.String(gd.Description),
		Public:      github.Bool(gd.IsPublic),
		Files: map[github.GistFilename]github.GistFile{
			github.GistFilename(fileName): {
				Content: github.String(string(context)),
			},
		},
	}
	gist, _, err = client.Gists.Create(ctx, gist)
	if err != nil {
		return err
	}
	url := fmt.Sprintf("https://gist.github.com/%v/%s", gist.GetOwner().GetLogin(), gist.GetID())
	slog.Info("Create Gist🎉", "URL", url)
	return nil
}

func main() {
	var d string
	var file string
	var p bool
	var ver bool
	var hc bool

	flag.StringVar(&d, "d", "", "gist description")
	flag.StringVar(&file, "f", "", "upload file")
	flag.BoolVar(&p, "p", false, "public check(default: false)")
	flag.BoolVar(&ver, "v", false, "show version")
	flag.BoolVar(&hc, "hc", false, "show health")
	flag.Parse()

	if ver {
		fmt.Println(version)
		os.Exit(0)
	}

	ctx := context.Background()

	client := github.NewClient(config.OauthGithub(ctx))

	if hc {
		healthCheck(client)
		os.Exit(0)
	}

	err := createGist(ctx, client, &GistDesciption{IsPublic: p, Description: d}, file)
	if err != nil {
		slog.Error("Faile createMemo.", "Mes", err)
		return
	}
}
