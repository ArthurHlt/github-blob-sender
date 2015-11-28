package gbs
import (
	"path"
	"os"
	"github.com/cloudfoundry/cli/cf/errors"
	"fmt"
	"io/ioutil"
	"encoding/base64"
	"path/filepath"
	"github.com/google/go-github/github"
	"sync"
	"crypto/sha1"
	"encoding/hex"
	"github.com/daviddengcn/go-colortext"
	"strings"
)


func PushFiles(owner, repo string, files ...string) error {
	var wg sync.WaitGroup
	wg.Add(len(files))
	var err error
	for _, file := range files {
		go func(owner, repo, file string) {
			defer wg.Done()
			err = PushFile(owner, repo, file)
			if err != nil {
				ct.Foreground(ct.Red, false)
				fmt.Print("Error when sending file '")
				ct.Foreground(ct.Cyan, false)
				fmt.Print(file)
				ct.Foreground(ct.Red, false)
				fmt.Print("': ")
				ct.ResetColor()
				fmt.Printf("%v\n", err)
			}else {
				ct.Foreground(ct.Green, false)
				fmt.Print("File '")
				ct.Foreground(ct.Cyan, false)
				fmt.Print(file)
				ct.Foreground(ct.Green, false)
				fmt.Print("' sent to '")
				ct.Foreground(ct.Cyan, false)
				fmt.Printf("github.com/%s/%s", owner, repo)
				ct.Foreground(ct.Green, false)
				fmt.Println("'")
				ct.ResetColor()
			}
		}(owner, repo, file)
	}
	wg.Wait()
	return nil
}
func CreateAllFolders(folder string) error {
	_, err := os.Stat(folder)
	if err == nil {
		return nil
	}
	return os.MkdirAll(folder, 0755)
}
func DownloadAll(folder string, createFolder bool) error {
	var wg sync.WaitGroup
	var err error
	blobFiles := Config.BlobFiles

	if createFolder {
		CreateAllFolders(folder)
	}
	_, err = os.Stat(folder)
	if os.IsNotExist(err) {
		return errors.New("Folder '" + folder + "' doesn't exist, use flag --create-folders to create it automatically")
	}

	wg.Add(len(blobFiles))
	for _, blobFile := range blobFiles {
		go func(blobFile BlobFile, folder string) {
			defer wg.Done()
			info := strings.Split(blobFile.Name, "/")
			owner := info[0]
			repo := info[1]
			file := info[2]
			err = DownloadFileTo(owner, repo, file, path.Join(folder, file))
			if err != nil {
				ct.Foreground(ct.Red, false)
				fmt.Print("Error when downloading file '")
				ct.Foreground(ct.Cyan, false)
				fmt.Print(file)
				ct.Foreground(ct.Red, false)
				fmt.Print("': ")
				ct.ResetColor()
				fmt.Printf("%v\n", err)
			}else {
				ct.Foreground(ct.Green, false)
				fmt.Print("File '")
				ct.Foreground(ct.Cyan, false)
				fmt.Print(file)
				ct.Foreground(ct.Green, false)
				fmt.Print("' downloaded from '")
				ct.Foreground(ct.Cyan, false)
				fmt.Printf("github.com/%s/%s", owner, repo)
				ct.Foreground(ct.Green, false)
				fmt.Print("' and placed in '")
				ct.Foreground(ct.Cyan, false)
				fmt.Printf("%s", path.Join(folder, file))
				ct.Foreground(ct.Green, false)
				fmt.Println("'")
				ct.ResetColor()
			}
		}(blobFile, folder)
	}
	wg.Wait()
	return nil
}
func DownloadFileTo(owner, repo, file, output string) error {
	if output == "" {
		output = file
	}
	blobFile := Config.FindBlobFile(BlobFile{
		Name: owner + "/" + repo + "/" + file,
	})
	if blobFile == nil {
		return errors.New(fmt.Sprintf("Cannot find file '%s' in github.", file))
	}
	dataDownloaded, err := GetFile(owner, repo, file)
	if err != nil {
		return err
	}
	shaCalc := sha1.New()
	shaCalc.Write(dataDownloaded)
	shaEncoded := hex.EncodeToString(shaCalc.Sum(nil))
	if shaEncoded != blobFile.InternalSha1 {
		return errors.New("Checksum doesn't match between downloaded file and saved sha")
	}
	return ioutil.WriteFile(output, dataDownloaded, 0644)
}
func GetFile(owner, repo, file string) ([]byte, error) {
	client := GetGithubClient()
	blobFile := Config.FindBlobFile(BlobFile{
		Name: owner + "/" + repo + "/" + file,
	})
	if blobFile == nil {
		return nil, errors.New(fmt.Sprintf("Cannot find file '%s' in github.", file))
	}
	checkBlob, _, err := client.Git.GetBlob(owner, repo, blobFile.GithubSha1)
	if err != nil {
		return nil, err
	}
	base64dataDecoded, err := base64.StdEncoding.DecodeString(*checkBlob.Content)
	if err != nil {
		return nil, err
	}
	return base64dataDecoded, nil
}
func PushFile(owner, repo, file string) error {
	client := GetGithubClient()
	file, err := getRealFilePath(file)
	if err != nil {
		return err
	}
	dataFile, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}

	base64data := base64.StdEncoding.EncodeToString(dataFile)

	fileName := filepath.Base(file)
	encodingBlob := "base64"
	shaCalc := sha1.New()
	shaCalc.Write(dataFile)
	shaEncoded := hex.EncodeToString(shaCalc.Sum(nil))
	blob := &github.Blob{
		Content: &base64data,
		Encoding: &encodingBlob,
	}
	finalBlob, _, err := client.Git.CreateBlob(owner, repo, blob)
	if err != nil {
		return err
	}
	blobFile := BlobFile{
		InternalSha1: shaEncoded,
		GithubSha1: *finalBlob.SHA,
		Name: owner + "/" + repo + "/" + fileName,
		Link: *finalBlob.URL,
	}
	Config.AddOrReplaceFile(blobFile)
	return nil
}

func getRealFilePath(file string) (string, error) {
	wd, err := os.Getwd();
	if err != nil {
		return "", err
	}
	if _, err := os.Stat(file); os.IsNotExist(err) {
		file = path.Join(wd, file)
	}else {
		return file, nil
	}
	if _, err := os.Stat(file); os.IsNotExist(err) {
		return "", errors.New(fmt.Sprintf("File '%s' doesn't exist", file))
	}
	return file, nil
}