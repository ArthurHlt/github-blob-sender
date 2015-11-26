package gbs
import (
	"fmt"
	"os/user"
	"os"
	"github.com/docker/docker/vendor/src/github.com/jfrazelle/go/canonical/json"
	"path"
	"io/ioutil"
)

type ConfigModel struct {
	Token     string `json:"token"`
	BlobFiles []BlobFile `json:"blob_files,omitempty"`
}
type BlobFile struct {
	Sha1 string `json:"sha1"`
	Name string `json:"name"`
	Link string `json:"string"`
}
var CONFIG_FILE = ".github-blob-sender"
var Config *ConfigModel
var ConfigPath string
func init() {
	Config = &ConfigModel{
		Token: "",
		BlobFiles: make([]BlobFile, 0),
	}
	usr, err := user.Current()
	checkError(err)
	ConfigPath = path.Join(usr.HomeDir, CONFIG_FILE)
	if _, err := os.Stat(ConfigPath); os.IsNotExist(err) {
		Config.save()
		return
	}
	data, err := ioutil.ReadFile(ConfigPath)
	checkError(err)
	err = json.Unmarshal(data, Config)
	checkError(err)
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
func (this *ConfigModel) AddFile(blobFile BlobFile) {
	this.BlobFiles = append(this.BlobFiles, blobFile)
	this.save()
}
func (this *ConfigModel) save() {
	data, err := json.Marshal(this)
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