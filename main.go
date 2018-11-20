package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
)

const GITHUB_API = `https://api.github.com/repos/%s/%s/releases/latest`

type Body struct {
	Assets []Asset `json:"assets"`
}

type Asset struct {
	BrowserDownloadUrl string `json:"browser_download_url"`
}

func Download(url string) error {
	cmd := exec.Command("curl", url)
	log.Println("\nRunning command and waiting for it to finish...")
	err := cmd.Run()
	log.Printf("Command finished with error: %v", err)
	return err
}

func main() {
	length := len(os.Args)

	if length < 3 {
		fmt.Println("Not found username or repository name")
		os.Exit(1)
	}

	Url := fmt.Sprintf(GITHUB_API, os.Args[1], os.Args[2])
	res, err := http.Get(Url)
	if err != nil {
		log.Fatal(err)
	}
	data, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		log.Fatal(err)
	}
	var body Body
	json.Unmarshal(data, &body)

	err = Download(body.Assets[0].BrowserDownloadUrl)
	if err != nil {
		println(err.Error)
		os.Exit(1)
	}
}
