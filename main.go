package main

import (
	"golang.org/x/oauth2"
	"github.com/google/go-github/github"
	"fmt"
	"os"
	"github.com/ArthurHlt/github-blob-sender/gbs"
	"github.com/codegangsta/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "github-blob-sender"
	app.Usage = "Send file to github blob"
	app.Commands = []cli.Command{
		{
			Name:      "test",
			Aliases:     []string{},
			Usage:     "Just run a simple test",
			Action: test,
		},
	}
	app.Run(os.Args)
}

func test(c *cli.Context) {
	config := gbs.Config
	config.SetToken("9d64dca319640a791cf87fcef4409a274eddd790")
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: config.GetToken()},
	)
	tc := oauth2.NewClient(oauth2.NoContext, ts)

	client := github.NewClient(tc)

	// list all repositories for the authenticated user
	req, err := client.NewRequest("GET", "/authorizations", nil)
	checkError(err)
	resp, err := client.Do(req, nil)
	if resp.StatusCode != 403 {
		fmt.Println("not valid token")
		os.Exit(1)
	}
	_, _, err = client.Repositories.List("", nil)
	if err != nil {
		fmt.Println("%v", err)
	}
	//fmt.Println("%#v", repos)
}

func checkError(err error) {
	if err != nil {
		fmt.Println(fmt.Sprintf("%v", err))
		os.Exit(1)
	}
}