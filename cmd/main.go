package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"sherlockgo/cmd/color"
	"strings"
)

var foundAccounts []string

func main() {
	file, err := os.Open("./resources/data.json")
	handleError(err)
	var endpoints map[string]interface{}
	err = json.NewDecoder(file).Decode(&endpoints)
	handleError(err)
	fmt.Print(`   ______            __                 _____           __
  / ____/___  ____  / /_  ___  _____   / __(_)___  ____/ /
 / / __/ __ \/ __ \/ __ \/ _ \/ ___/  / /_/ / __ \/ __  / 
/ /_/ / /_/ / /_/ / / / /  __/ /     / __/ / / / / /_/ /  
\____/\____/ .___/_/ /_/\___/_/     /_/ /_/_/ /_/\__,_/   
          /_/                                             

`)
	fmt.Print("🐹🔎Who are you looking for? ")
	username := getInput()

	for websiteName, parameter := range endpoints {
		websiteURL := parameter.(map[string]interface{})["url"]
		checkURL(websiteURL, websiteName, username)
	}
	fmt.Printf("All websites checked! I created a file called %s.txt containing the links.🐹🔎", username)
	generateFileWithFoundAcconts(foundAccounts, username)
	defer file.Close()
}

func checkURL(websiteURL interface{}, websiteName interface{}, username string) {
	resp, err := http.Get(urlWithUsername(websiteURL.(string), username))
	// if timeout, skip
	if err != nil {
		return
	}
	defer resp.Body.Close()
	handleError(err)
	checkStatusCode(resp, websiteName, urlWithUsername(websiteURL.(string), username))
}

func checkStatusCode(resp *http.Response, websiteName interface{}, url string) {
	if resp.StatusCode == 200 {
		fmt.Println(color.Green+"[+] FOUND -", websiteName, color.Reset)
		fmt.Println(url)
		foundAccounts = append(foundAccounts, url)
	} else {
		fmt.Println(color.Red+"[-] NOT FOUND -", websiteName, color.Reset)
	}
}

func urlWithUsername(url string, username string) string {
	return strings.Replace(url, "{}", username, -1)
}

func handleError(err error) {
	if err != nil {
		panic(err)
	}
}

func generateFileWithFoundAcconts(foundAccounts []string, fileName string) {
	file, err := os.Create(fmt.Sprintf("./%s.txt", fileName))
	handleError(err)
	defer file.Close()
	for _, account := range foundAccounts {
		file.WriteString(account + "\n")
	}
}

func getInput() string {
	var input string
	fmt.Scanln(&input)
	return input
}
