package main

import (
	"encoding/json"
	"flag"
	"github.com/ncw/swift"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

var (
	conn       swift.Connection
	sourceDir  = flag.String("s", "", "Directory where we can find the files that need to be send to swift")
	targetDir  = flag.String("t", "", "Target directory in container where we will store the files")
	configFile = flag.String("c", "", "Location of configuration file")
)

type connectionInfo struct {
	ApiKey        string `json:"api_key"`
	AuthURL       string `json:"auth_url"`
	Tenant        string `json:"tenant"`
	UserName      string `json:"username"`
	ContainerName string `json:"container"`
}

func loadConfig() connectionInfo {
	configJSON, err := ioutil.ReadFile(*configFile)

	if err != nil {
		log.Fatal(err)
	}

	var config connectionInfo
	err = json.Unmarshal(configJSON, &config)
	return config
}

type fileSearcher struct {
	latestFile    string
	latestModTime time.Time
}

func (f *fileSearcher) visit(path string, fi os.FileInfo, err error) error {

	if fi.IsDir() {
		return nil
	}

	modTime := fi.ModTime()

	diff := f.latestModTime.Sub(modTime)

	if diff.Seconds() < 0 {
		f.latestFile = path
		f.latestModTime = modTime
	}

	return nil
}

func (f fileSearcher) getNewestFile() string {
	f.latestModTime = time.Date(2000, time.January, 1, 0, 0, 0, 0, time.UTC)
	filepath.Walk(*sourceDir, f.visit)
	return f.latestFile
}

func main() {
	flag.Parse()
	config := loadConfig()

	conn = swift.Connection{
		UserName: config.UserName,
		ApiKey:   config.ApiKey,
		AuthUrl:  config.AuthURL,
		Tenant:   config.Tenant,
	}
	log.Println("Walking over " + *sourceDir)

	fs := fileSearcher{}
	newestFile := fs.getNewestFile()

	log.Println("Most recent file: " + newestFile)

	log.Println("Sending to directory in ObjectStore: " + *targetDir)

	data, err := ioutil.ReadFile(newestFile)

	if err != nil {
		log.Fatal(err)
	}

	newestFile = strings.Replace(newestFile, *sourceDir, "", 1)
	targetPath := *targetDir + "/" + newestFile

	err = conn.ObjectPutBytes(config.ContainerName, targetPath, data, http.DetectContentType(data))

	if err != nil {
		log.Fatal(err)
	}

	log.Println("Successfully finished")
}
