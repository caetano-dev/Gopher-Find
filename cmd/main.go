package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"gopher-find/cmd/color"
	"gopher-find/cmd/models"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

var foundAccounts []string

func main() {
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

	wd, err := os.Getwd()
	handleError(err)

	file, err := os.Open(filepath.FromSlash(wd + "/cmd/resources/data.json"))
	handleError(err)
	defer file.Close()

	var endpoints map[string]interface{}
	err = json.NewDecoder(file).Decode(&endpoints)
	handleError(err)

	username := os.Args[1]

	for websiteName, parameter := range endpoints {
		var data models.Parameter
		d, _ := json.Marshal(parameter)
		err := json.Unmarshal(d, &data)
		handleError(err)

		urlWithName := urlWithUsername(data.URL, username)

		if data.ErrorType == "message" {
			checkIfUserExistsByErrorMessage(websiteName, urlWithName, data.ErrorMsg)
		} else if data.ErrorType == "response_url" {
			checkIfUserExistsByRedirect(websiteName, urlWithName)
		} else {
			checkIfUserExistsByStatusCode(websiteName, urlWithName)
		}
	}

	fmt.Printf("All websites checked! I created a file called %s.txt containing the links.🐹🔎", username)
	generateFileWithFoundAcconts(foundAccounts, username)
}

func checkIfUserExistsByErrorMessage(websiteName string, urlWithUsername string, errorMessage string) {
	if strings.Contains(websiteScrape(urlWithUsername), errorMessage) {
		fmt.Println(color.Red+"[-] NOT FOUND -", websiteName, color.Reset)
	} else {
		fmt.Println(color.Green+"[+] FOUND -", websiteName, color.Reset)
		fmt.Println(urlWithUsername)
		foundAccounts = append(foundAccounts, urlWithUsername)
	}
}

func checkIfUserExistsByStatusCode(websiteName string, urlWithUsername string) {
	resp, err := http.Get(urlWithUsername)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		fmt.Println(color.Green+"[+] FOUND -", websiteName, color.Reset)
		fmt.Println(urlWithUsername)
		foundAccounts = append(foundAccounts, urlWithUsername)
	} else {
		fmt.Println(color.Red+"[-] NOT FOUND -", websiteName, color.Reset)
	}
}

func checkIfUserExistsByRedirect(websiteName string, urlWithUsername string) {
	req, err := http.NewRequest("GET", urlWithUsername, nil)
	handleError(err)
	client := new(http.Client)
	client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return errors.New("redirect")
	}

	response, err := client.Do(req)
	if err == nil {
		if response.StatusCode == 302 {
			fmt.Println(color.Red+"[-] NOT FOUND -", websiteName, color.Reset)
		} else {
			fmt.Println(color.Green+"[+] FOUND -", websiteName, color.Reset)
			fmt.Println(urlWithUsername)
		}

	}
}

func websiteScrape(urlWithUsername string) string {
	res, err := http.Get(urlWithUsername)
	handleError(err)
	defer res.Body.Close()
	if res.StatusCode != 200 {
		fmt.Println("Unable to access website due to captcha/JavaScript/Cloudflare.")
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return ""
	}

	var websiteContent []string
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
		log.Fatal(err)
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
