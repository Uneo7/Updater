package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"sync"
)

type Files []File
type File struct {
	Name string `json:"path"`
	Sum  string `json:"sum"`
}

type Config struct {
	Index string
	Root  string
	Game  string
}

var config Config

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
	download.File = file
	download.Root = dir

	if !download.verify() {
		download.download()
	}
}

func parseConf() {

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Updater parameters:\n")
		flag.PrintDefaults()
	}

	dir := flag.String("g", "", "Game directory path (absolute)")
	root := flag.String("r", "", "Files location URL")
	index := flag.String("i", "", "Files index")

	flag.Parse()

	config.Game = *dir
	config.Index = *index
	config.Root = strings.TrimRight(*root, "/")
}

func main() {

	parseConf()

	files, err := requestFiles(config.Index)
	if err != nil {
		panic(err.Error())
	}

	fmt.Println("Files: ", len(files))

	downloadsQueue := &sync.WaitGroup{}

	for i := range files {
		downloadsQueue.Add(1)
		go get(config.Game, files[i].Name, files[i].Sum, downloadsQueue)
	}

	downloadsQueue.Wait()
}
