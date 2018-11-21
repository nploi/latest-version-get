package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"
)

const GITHUB_API = `https://api.github.com/repos/%s/%s/releases/latest`

type Body struct {
	Assets []Asset `json:"assets"`
}

type Asset struct {
	BrowserDownloadUrl string `json:"browser_download_url"`
}

func getNameFromURL(url string) string {
	items := strings.Split(url, `/`)
	leng := len(items)
	return items[leng-1]
}

func download(url, path string) error {
	name := getNameFromURL(url)
	if len(path) == 0 {
		path = "./" + name
	} else if path[len(path)-1] == '/' {
		path += name
	} else {
		path += "/" + name
	}

	cmd := exec.Command("wget", url, "-O", path)
	log.Println("Downloading " + name + ", waiting for it to finish...")
	err := cmd.Run()
	log.Printf("Download with error: %v", err)
	return err
}

func main() {
	length := len(os.Args)
	if length < 3 {
		fmt.Println("Not found username or repository name\ngo main.go [username] [repository_name] [path]")
		os.Exit(1)
	}

	// Send request
	Url := fmt.Sprintf(GITHUB_API, os.Args[1], os.Args[2])
	res, err := http.Get(Url)
	if err != nil {
		log.Fatal(err)
	}

	// Read data
	data, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		log.Fatal(err)
	}
	// Parse json
	var body Body
	json.Unmarshal(data, &body)

	if len(body.Assets) == 0 {
		println("Not found release version!!")
		os.Exit(1)
	}

	if length > 3 {
		err = download(body.Assets[0].BrowserDownloadUrl, os.Args[3])
	} else {
		err = download(body.Assets[0].BrowserDownloadUrl, ".")
	}

	if err != nil {
		println(err.Error)
		os.Exit(1)
	}
}
