package main

import (
	"fmt"
	"github.com/pierrec/xxHash/xxHash32"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"strings"
)

type Downloader struct {
	Root string
	File string
	Hash string
}

func (downloader Downloader) createPath(dirPath string) {
	err := os.MkdirAll(path.Dir(dirPath), os.ModePerm)
	if err != nil {
		panic(err)
	}
}

func (downloader Downloader) download() bool {

	var url strings.Builder
	url.WriteString(root)
	url.WriteString(downloader.File)

	file := path.Join(downloader.Root, downloader.File)

	if _, err := os.Stat(path.Dir(file)); os.IsNotExist(err) {
		downloader.createPath(file)
	}

	output, err := os.Create(file)
	if err != nil {
		fmt.Println("Error while creating:", file, "-", err)
		return false
	}
	defer output.Close()

	response, err := http.Get(url.String())
	if err != nil {
		fmt.Println("Error while downloading:", url.String(), "-", err)
		return false
	}
	defer response.Body.Close()

	_, err = io.Copy(output, response.Body)
	if err != nil {
		fmt.Println("Error while downloading:", url.String(), "-", err)
		return false
	}

	fmt.Println("Downloaded: ", downloader.File)

	return true
}

func (downloader Downloader) verify() bool {

	fmt.Println("Verifying: ", downloader.File)
	file := path.Join(downloader.Root, downloader.File)

	buffer, err := ioutil.ReadFile(file)
	if err != nil {
		return false
	}

	x := xxHash32.Checksum(buffer, 0)
	return fmt.Sprintf("%x", x) == downloader.Hash
}
