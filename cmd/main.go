// Package main is the main package.
package main

import (
	"encoding/json"
	"fmt"
	c "gopher-find/cmd/checkUsers"
	"gopher-find/cmd/color"
	"gopher-find/cmd/models"
	"gopher-find/cmd/utils"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"sync"
	"sync/atomic"
	"time"
)

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

	var endpoints map[string]models.Parameter
	err = json.NewDecoder(file).Decode(&endpoints)
	handleError(err)

	username := os.Args[1]

	rateLimiter := utils.NewRateLimiter()
	var wg sync.WaitGroup
	wg.Add(len(endpoints))

	count := int64(len(endpoints))

	for websiteName, parameter := range endpoints {
		w := websiteName
		p := parameter

		go func() {
			defer func() {
				wg.Done()
				atomic.AddInt64(&count, -1)
			}()

			urlWithName := c.URLWithUsername(p.URL, username)

			parsedURL, err := url.Parse(urlWithName)
			if err != nil {
				fmt.Printf("Error parsing URL for %s: %v\n", w, err)
				return
			}

			for !rateLimiter.Allow(parsedURL.Host) {
				time.Sleep(100 * time.Millisecond)
			}

			if p.ErrorType == "message" {
				c.CheckIfUserExistsByErrorMessage(w, urlWithName, p)
			} else if p.ErrorType == "response_url" {
				c.CheckIfUserExistsByRedirect(w, urlWithName, p)
			} else {
				c.CheckIfUserExistsByStatusCode(w, urlWithName, p)
			}
		}()
	}

	wg.Wait()

	fmt.Printf("All websites checked! I created a file called %s.txt containing the links.🐹🔎", username)
	generateFileWithFoundAcconts(c.FoundAccounts, username)
}

func generateFileWithFoundAcconts(foundAccounts []string, fileName string) {
	file, err := os.Create(fmt.Sprintf("./%s.txt", fileName))
	handleError(err)
	defer file.Close()
	file.WriteString("WARNING!\n Websites that return false positives are included with a warn. They are added in the file because we believe that it is better to assume these accounts exist and manually check them instead of possibly missing results. We are working to solve this inconvenience and reduce the amount of bad entries.\n------------------------------------------------------------\n")
	for _, account := range foundAccounts {
		file.WriteString(account + "\n")
	}
}
func handleError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
