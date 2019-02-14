package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"sync"
)

type Files []File
type File struct {
	Name string `json:"path"`
	Sum  string `json:"sum"`
}

var index = "https://update.beyond-horizon.fr/files.json"
var root = "https://update.beyond-horizon.fr/update"

var downloadsQueue = &sync.WaitGroup{}

func requestFiles(url string) (Files, error) {

	res, err := http.Get(url)
	if err != nil {
		return Files{}, err
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return Files{}, err
	}

	var files = Files{}
	err = json.Unmarshal(body, &files)
	if err != nil {
		return Files{}, err
	}

	return files, nil
}

func get(dir string, file string, hash string, queue *sync.WaitGroup) {

	defer queue.Done()
	var download Downloader

	download.Hash = hash
	download.File = strings.Replace(file, "/update", "", 1)
	download.Root = dir

	if !download.verify() {
		download.download()
	}
}

func main() {

	dir := flag.String("gamedir", "", "--gamedir Dossier du jeu")
	flag.Parse()

	fmt.Println(*dir)
	if *dir == "" {
		panic("Gamedir required")
	}

	files, err := requestFiles(index)
	if err != nil {
		panic(err.Error())
	}

	fmt.Println("Files: ", len(files))

	for i := range files {
		downloadsQueue.Add(1)
		go get(*dir, files[i].Name, files[i].Sum, downloadsQueue)
	}

	downloadsQueue.Wait()
}
