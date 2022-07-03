package main

import (
	"encoding/json"
	"fmt"
	"gopher-find/cmd/color"
	"net/http"
	"os"
	"strings"

	"github.com/PuerkitoBio/goquery"
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
	if len(os.Args[1:]) == 0 {
		fmt.Println(color.Red + "[!] Username is empty" + color.Reset)
		fmt.Println(color.Red + "[!] Usage: ./gopher-find <username>" + color.Reset)
		os.Exit(1)
	}
	username := os.Args[1]
	for websiteName, parameter := range endpoints {
		websiteURL := parameter.(map[string]interface{})["url"]
		errorType := parameter.(map[string]interface{})["errorType"]
		errorMessage := parameter.(map[string]interface{})["errorMsg"]
		if errorType == "message" {
			checkIfUserExistsByErrorMessage(websiteName, urlWithUsername(websiteURL.(string), username), errorMessage.(string))
		} else {
			checkIfUserExistsByStatusCode(getStatuscode(websiteURL, username), websiteName, urlWithUsername(websiteURL.(string), username))
		}

	}
	fmt.Printf("All websites checked! I created a file called %s.txt containing the links.🐹🔎", username)
	generateFileWithFoundAcconts(foundAccounts, username)
	defer file.Close()
}
func checkIfUserExistsByErrorMessage(websiteName interface{}, urlWithUsername string, errorMessage string) {
	if strings.Contains(websiteScrape(urlWithUsername), errorMessage) {
		fmt.Println(color.Red+"[-] NOT FOUND -", websiteName.(string), color.Reset)
	} else {
		fmt.Println(color.Green+"[+] FOUND -", websiteName.(string), color.Reset)
		fmt.Println(urlWithUsername)
		foundAccounts = append(foundAccounts, urlWithUsername)
	}
}

func checkIfUserExistsByStatusCode(statusCode int, websiteName interface{}, urlWithUsername string) {
	if statusCode == 200 {
		fmt.Println(color.Green+"[+] FOUND -", websiteName, color.Reset)
		fmt.Println(urlWithUsername)
		foundAccounts = append(foundAccounts, urlWithUsername)
	} else {
		fmt.Println(color.Red+"[-] NOT FOUND -", websiteName, color.Reset)
	}
}

func getStatuscode(websiteURL interface{}, username string) int {
	resp, err := http.Get(urlWithUsername(websiteURL.(string), username))
	if err != nil {
		return 0
	}
	defer resp.Body.Close()
	handleError(err)
	return resp.StatusCode
}

func websiteScrape(urlWithUsername string) string {
	var websiteContent []string
	doc, err := goquery.NewDocument(urlWithUsername)
	if err != nil {
		return ""
	}
	doc.Find("html").Each(func(index int, item *goquery.Selection) {
		websiteContent = append(websiteContent, item.Text())
	})
	return strings.Join(websiteContent, " ")
}

func urlWithUsername(websiteURL string, username string) string {
	return strings.Replace(websiteURL, "{}", username, -1)
}

func handleError(err error) {
	if err != nil {
		fmt.Println("error")
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
