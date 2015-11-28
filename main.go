package main

import (
	"fmt"
	"os"
	"github.com/ArthurHlt/github-blob-sender/gbs"
	"github.com/codegangsta/cli"
	"github.com/cloudfoundry/cli/cf/errors"
	"github.com/daviddengcn/go-colortext"
	"github.com/olekukonko/tablewriter"
	"strings"
)
var githubToken string
var owner string
var repo string
var outputDownload string

var showGithubSha1 bool
var showInternalSha1 bool
var showLink bool
var createFolder bool

func main() {
	app := cli.NewApp()
	app.Name = "github-blob-sender"
	app.Usage = "Store and restore file from github blob api"
	app.Version = "1.0.0"
	app.Commands = []cli.Command{
		{
			Name:      "upload",
			Aliases:     []string{"u"},
			Usage:     "Upload file to github blob",
			ArgsUsage: "[file1 file2 ...]",
			Action: upload,
			Flags: defaultFlags(false),
		},
		{
			Name:      "cat",
			Aliases:     []string{"c"},
			Usage:     "Cat file from github blob",
			ArgsUsage: "[file-name] (Note: file-name can be listed with list command)",
			Action: cat,
			Flags: defaultFlags(true),
		},
		{
			Name:      "download",
			Aliases:     []string{"d"},
			Usage:     "Download file from github blob",
			ArgsUsage: "[file-name] (Note: file-name can be listed with list command)",
			Action: download,
			Flags: append(defaultFlags(true), cli.StringFlag{
				Name:        "output",
				Value:       "",
				Usage:       "Set where to write the downloaded file",
				Destination: &outputDownload,
			}),
		},
		{
			Name:      "download-all",
			Aliases:     []string{"a"},
			ArgsUsage: "[folder/to/put/downloaded/file]",
			Usage:     "Download all registered files in folder from github blob",
			Action: downloadAll,
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name: "create-folders, create-folder, c",
					Usage: "Create folders if not exist",
					Destination: &createFolder,
				},
			},
		},
		{
			Name:      "list",
			Aliases:     []string{"l"},
			Usage:     "List registered files",
			Action: list,
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name: "show-github-sha1, g",
					Usage: "Show github checksum (in sha1)",
					Destination: &showGithubSha1,
				},
				cli.BoolFlag{
					Name: "show-registered-sha1, r",
					Usage: "Show registered checksum (in sha1)",
					Destination: &showInternalSha1,
				},
				cli.BoolFlag{
					Name: "show-link, l",
					Usage: "Show github link",
					Destination: &showLink,
				},
			},
		},
	}
	app.Run(os.Args)
}
func defaultFlags(optional bool) []cli.Flag {
	var optionnalText string = ""
	if optional {
		optionnalText = " (optional)"
	}
	return []cli.Flag{
		cli.StringFlag{
			Name:        "github-token, gt",
			Value:       "",
			Usage:       "Set your github token (optional if already set or " + gbs.GITHUB_TOKEN_ENV + " env var set)",
			Destination: &githubToken,
		},
		cli.StringFlag{
			Name:        "owner, o",
			Value:       "",
			Usage:       "Which org or user you own to send file" + optionnalText,
			Destination: &owner,
		},
		cli.StringFlag{
			Name:        "repo, r",
			Value:       "",
			Usage:       "Which repo you own to send file" + optionnalText,
			Destination: &repo,
		},
	}
}
func downloadAll(c *cli.Context) {
	checkToken()
	if len(c.Args()) != 1 {
		checkError(errors.New("You need to pass a folder as first arg."))
	}
	err := gbs.DownloadAll(c.Args().First(), createFolder)
	checkError(err)
}
func download(c *cli.Context) {
	checkToken()
	if len(c.Args()) != 1 {
		checkError(errors.New("You need to pass one file name."))
	}
	tryFoundNotAmbiguousFile(c.Args().First())
	err := gbs.DownloadFileTo(owner, repo, c.Args().First(), outputDownload)
	checkError(err)
	ct.Foreground(ct.Green, false)
	fmt.Println("File has been downloaded")
	ct.ResetColor()
}
func tryFoundNotAmbiguousFile(arg string) {
	blobFilesFound := gbs.Config.FindBlobFilesWithFileName(arg)
	if len(blobFilesFound) == 0 {
		ct.Foreground(ct.Red, false)
		fmt.Println("This file is not registered.")
		ct.ResetColor()
		os.Exit(1)
	}
	if len(blobFilesFound) > 1 {
		checkRequiredInformation()
	}else {
		info := strings.Split(blobFilesFound[0].Name, "/")
		owner = info[0]
		repo = info[1]
	}
}
func cat(c *cli.Context) {
	checkToken()
	if len(c.Args()) != 1 {
		checkError(errors.New("You need to pass one file name."))
	}
	tryFoundNotAmbiguousFile(c.Args().First())
	data, err := gbs.GetFile(owner, repo, c.Args().First())
	checkError(err)
	fmt.Printf("%s\n", data)
}

func list(c *cli.Context) {
	data := [][]string{}
	for _, blobFile := range gbs.Config.BlobFiles {
		unsplitInfo := strings.Split(blobFile.Name, "/")
		content := []string{
			unsplitInfo[2],
			unsplitInfo[0],
			unsplitInfo[1],
		}
		if showGithubSha1 {
			content = append(content, blobFile.GithubSha1)
		}
		if showInternalSha1 {
			content = append(content, blobFile.InternalSha1)
		}
		if showLink {
			content = append(content, blobFile.Link)
		}
		data = append(data, content)
	}
	table := tablewriter.NewWriter(os.Stdout)
	header := []string{"Name", "Owner", "Repo"}
	if showGithubSha1 {
		header = append(header, "Github checksum")
	}
	if showInternalSha1 {
		header = append(header, "Registered checksum")
	}
	if showLink {
		header = append(header, "Github url")
	}
	table.SetHeader(header)

	for _, v := range data {
		table.Append(v)
	}
	table.Render()
}

func upload(c *cli.Context) {
	checkToken()
	checkRequiredInformation()
	if len(c.Args()) == 0 {
		checkError(errors.New("You need to pass one or multiple file(s)."))
	}
	err := gbs.PushFiles(owner, repo, c.Args()...)
	checkError(err)
}
func checkToken() {
	for gbs.Config.GetToken() == "" || !gbs.IsValidToken(gbs.Config.GetToken()) {
		if githubToken != "" {
			gbs.Config.SetToken(githubToken)
			continue
		}
		ct.Foreground(ct.Cyan, false)
		fmt.Println("You need to provide valid github token")
		fmt.Println("Generate a default one here: https://github.com/settings/tokens/new\n")
		ct.ResetColor()
		fmt.Print("Your token: ")
		_, err := fmt.Scanf("%s\n", &githubToken)
		if err != nil {
			fmt.Println(fmt.Sprintf("%v", err))
			continue
		}
		gbs.Config.SetToken(githubToken)
	}
}
func checkRequiredInformation() {
	if owner == "" {
		fmt.Print("Owner (github user or org): ")
		_, err := fmt.Scanf("%s\n", &owner)
		checkError(err)
	}
	if repo == "" {
		fmt.Print("Github repo: ")
		_, err := fmt.Scanf("%s\n", &repo)
		checkError(err)
	}
}

func checkError(err error) {
	if err != nil {
		ct.Foreground(ct.Red, false)
		fmt.Println(fmt.Sprintf("%v", err))
		ct.ResetColor()
		os.Exit(1)
	}
}