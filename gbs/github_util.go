package gbs
import (
	"golang.org/x/oauth2"
	"github.com/google/go-github/github"
)
var githubClient *github.Client
func IsValidToken(token string) bool {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(oauth2.NoContext, ts)

	client := github.NewClient(tc)

	// list all repositories for the authenticated user
	req, err := client.NewRequest("GET", "/authorizations", nil)
	checkError(err)
	resp, err := client.Do(req, nil)
	if resp.StatusCode != 403 {
		return false
	}
	return true
}

func GetGithubClient() *github.Client {
	if githubClient != nil {
		return githubClient
	}
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: Config.GetToken()},
	)
	tc := oauth2.NewClient(oauth2.NoContext, ts)

	githubClient = github.NewClient(tc)
	return githubClient
}