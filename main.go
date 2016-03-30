package main

import (
	"bytes"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type TReq struct {
	NameTest   string
	Url        string
	Host       string
	StatusCode int
	SubString  string
}

var tests []TReq
var verbose bool = false
var exitCode int = 0

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func readfile(filename string) {
	file, err := ioutil.ReadFile(filename)
	check(err)
	//fmt.Printf("%s\n", string(file))
	err = json.Unmarshal(file, &tests)
	check(err)
}

func TempFileName(prefix, suffix string) string {
	prefix = strings.Replace(prefix, " ", "_", -1)
	suffix = strings.Replace(suffix, " ", "_", -1)
	randBytes := make([]byte, 4)
	rand.Read(randBytes)
	return filepath.Join(prefix + "_" + hex.EncodeToString(randBytes) + suffix)
}

func checkURL(ip string, uri TReq) bool {
	var buffer bytes.Buffer
	buffer.WriteString("http://")
	buffer.WriteString(ip)
	buffer.WriteString(uri.Url)
	client := &http.Client{}
	req, _ := http.NewRequest("GET", buffer.String(), nil)
	req.Host = uri.Host
	req.Header.Set("Host", uri.Host)
	res, _ := client.Do(req)
	defer res.Body.Close()
	contents, err := ioutil.ReadAll(res.Body)
	check(err)
	if verbose {
		err := ioutil.WriteFile(TempFileName(uri.NameTest, ".html"), contents, 0644)
		check(err)
	}
	if res.StatusCode == uri.StatusCode {
		if strings.Contains(string(contents), uri.SubString) {
			return true
		}
	}
	return false
}

func main() {
	ipPtr := flag.String("ip", "127.0.0.1", "webserver's ip address")
	filePtr := flag.String("f", "ssv4.json", "filename with test")
	verbosePtr := flag.Bool("v", false, "verbose mode: save html pages to files")
	flag.Parse()
	verbose = *verbosePtr
	readfile(*filePtr)
	for index, element := range tests {
		if checkURL(*ipPtr, element) {
			fmt.Println("Test#" + strconv.Itoa(index) + ": " + element.NameTest + "  ... Pass")
		} else {
			fmt.Println("Test#" + strconv.Itoa(index) + ": " + element.NameTest + "  ... Failure")
			exitCode = 1
			os.Exit(exitCode)
		}
	}
	os.Exit(exitCode)
}
