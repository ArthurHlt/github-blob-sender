package gbs
import (
	"fmt"
	"os"
	"github.com/docker/docker/vendor/src/github.com/jfrazelle/go/canonical/json"
	"path"
	"io/ioutil"
	"strings"
)

type ConfigModel struct {
	Token     string `json:"token"`
	BlobFiles []BlobFile `json:"blob_files,omitempty"`
}
type BlobFile struct {
	InternalSha1 string `json:"internal_sha1"`
	GithubSha1   string `json:"github_sha1"`
	Name         string `json:"name"`
	Link         string `json:"string"`
}
var CONFIG_FILE = ".github-blob-sender"
var GITHUB_TOKEN_ENV = "GITHUB_TOKEN"
var Config *ConfigModel
var ConfigPath string
func init() {
	Config = &ConfigModel{
		Token: "",
		BlobFiles: make([]BlobFile, 0),
	}
	wd, err := os.Getwd()
	checkError(err)
	ConfigPath = path.Join(wd, CONFIG_FILE)
	if _, err := os.Stat(ConfigPath); os.IsNotExist(err) {
		loadFromEnvVar()
		Config.save()
		return
	}
	data, err := ioutil.ReadFile(ConfigPath)
	checkError(err)
	err = json.Unmarshal(data, Config)
	checkError(err)
}
func loadFromEnvVar() {
	token := os.Getenv(GITHUB_TOKEN_ENV)
	if token != "" {
		Config.SetToken(token)
	}
}
func (this *ConfigModel) GetToken() string {
	return this.Token
}
func (this *ConfigModel) GetBlobFiles() []BlobFile {
	return this.BlobFiles
}
func (this *ConfigModel) SetToken(token string) {
	this.Token = token
	this.save()
}
func (this *ConfigModel) AddOrReplaceFile(blobFile BlobFile) {
	index := this.FindIndexBlobFile(blobFile)
	if index == -1 {
		this.BlobFiles = append(this.BlobFiles, blobFile)
		this.save()
		return
	}
	this.BlobFiles[index].Name = blobFile.Name
	this.BlobFiles[index].Link = blobFile.Link
	this.BlobFiles[index].InternalSha1 = blobFile.InternalSha1
	this.save()
}

func (this *ConfigModel) FindIndexBlobFile(blobFileToFind BlobFile) int {
	for index, blobFile := range this.BlobFiles {
		if blobFile.Name == blobFileToFind.Name {
			return index
		}
	}
	return -1
}

func (this *ConfigModel) FindBlobFilesWithFileName(fileName string) []BlobFile {
	finalBlobFiles := make([]BlobFile, 0)
	for _, blobFile := range this.BlobFiles {
		info := strings.Split(blobFile.Name, "/")
		if info[2] == fileName {
			finalBlobFiles = append(finalBlobFiles, blobFile)
		}
	}
	return finalBlobFiles
}

func (this *ConfigModel) FindBlobFile(blobFileToFind BlobFile) *BlobFile {
	for _, blobFile := range this.BlobFiles {
		if blobFile.Name == blobFileToFind.Name {
			return &blobFile
		}
	}
	return nil
}
func (this *ConfigModel) save() {
	data, err := json.MarshalIndent(this, "", "  ")
	checkError(err)
	err = ioutil.WriteFile(ConfigPath, data, 0644)
	checkError(err)
}
func checkError(err error) {
	if err != nil {
		fmt.Println(fmt.Sprintf("%v", err))
		os.Exit(1)
	}
}