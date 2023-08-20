package config

import (
	"context"
	"net/http"
	"os"

	"golang.org/x/oauth2"
)

func OauthGithub(ctx context.Context) *http.Client {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: os.Getenv("GITHUB_TOKEN")},
	)
	tc := oauth2.NewClient(ctx, ts)
	return tc
}
