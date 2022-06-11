package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"sherlockgo/cmd/color"
	"strings"
)

func main() {
	file, err := os.Open("./resources/data.json")
	handleError(err)
	var endpoints map[string]interface{}
	err = json.NewDecoder(file).Decode(&endpoints)
	handleError(err)
	fmt.Print("Type the username: ")
	username := getInput()

	for websiteName, parameter := range endpoints {
		websiteURL := parameter.(map[string]interface{})["url"]
		checkURL(websiteURL, websiteName, username)
	}
	defer file.Close()
}

func checkURL(websiteURL interface{}, websiteName interface{}, username string) {
	url := URLWithUsername(websiteURL.(string), username)
	resp, err := http.Get(url)
	// if timeout, skip
	if err != nil {
		return
	}
	handleConnectionError(err)
	defer resp.Body.Close()
	handleError(err)
	checkStatusCode(resp, websiteName, url)
}

func checkStatusCode(resp *http.Response, websiteName interface{}, url string) {
	if resp.StatusCode == 200 {
		fmt.Println(color.Green+"[+] FOUND -", websiteName, color.Reset)
		fmt.Println(url)
	} else {
		fmt.Println(color.Red+"[-] NOT FOUND -", websiteName, color.Reset)
	}
}

func URLWithUsername(url string, username string) string {
	return strings.Replace(url, "{}", username, -1)
}

func handleError(err error) {
	if err != nil {
		panic(err)
	}
}

func getInput() string {
	var input string
	fmt.Scanln(&input)
	return input
}
