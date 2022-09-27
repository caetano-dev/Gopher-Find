// Package checkUsers contains the functions necessary for validating if a user exists in the websites.
package checkUsers

import (
	"errors"
	"fmt"
	"gopher-find/cmd/color"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

// FoundAccounts is an array with all of the found accounts.
var (
	FoundAccounts []string
	httpClient    = http.Client{Timeout: 30 * time.Second}
)

// Response is the http response structure with a code and a body.
type Response struct {
	code int
	body string
}

func doReq(url string) (Response, error) {
	resp, err := httpClient.Get(url)
	if err != nil {
		return Response{}, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return Response{}, err
	}

	return Response{body: string(body), code: resp.StatusCode}, nil
}

// CheckIfUserExistsByErrorMessage check if the user exists by error message.
func CheckIfUserExistsByErrorMessage(websiteName string, URLWithUsername string, errorMessage string, FalsePositive bool, falsePositiveMessage string) {
	if strings.Contains(websiteScrape(URLWithUsername), errorMessage) {
		fmt.Println(color.Red+"[-] NOT FOUND -", websiteName, color.Reset)
	} else {
		fmt.Println(color.Green+"[+] FOUND -", websiteName, color.Reset)
		fmt.Println(URLWithUsername)
		if FalsePositive {
			FoundAccounts = append(FoundAccounts, websiteName+" - "+URLWithUsername+falsePositiveMessage)
		} else {
			FoundAccounts = append(FoundAccounts, websiteName+" - "+URLWithUsername)
		}
	}
}

// CheckIfUserExistsByStatusCode check if the user exists by status code.
func CheckIfUserExistsByStatusCode(websiteName string, URLWithUsername string, FalsePositive bool, falsePositiveMessage string) {
	res, err := doReq(URLWithUsername)
	if err != nil {
		fmt.Println(err)
		return
	}

	if res.code == 200 {
		fmt.Println(color.Green+"[+] FOUND -", websiteName, color.Reset)
		fmt.Println(URLWithUsername)
		if FalsePositive {
			FoundAccounts = append(FoundAccounts, websiteName+" - "+URLWithUsername+falsePositiveMessage)
		} else {
			FoundAccounts = append(FoundAccounts, websiteName+" - "+URLWithUsername)
		}
	} else {
		fmt.Println(color.Red+"[-] NOT FOUND -", websiteName, color.Reset)
	}
}

// CheckIfUserExistsByRedirect check if the user exists by redirect.
func CheckIfUserExistsByRedirect(websiteName string, URLWithUsername string, FalsePositive bool, falsePositiveMessage string) {
	req, err := http.NewRequest("GET", URLWithUsername, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	client := httpClient
	client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return errors.New("redirect")
	}

	response, err := client.Do(req)
	if err == nil {
		if response.StatusCode == 302 {
			fmt.Println(color.Red+"[-] NOT FOUND -", websiteName, color.Reset)
		} else {
			fmt.Println(color.Green+"[+] FOUND -", websiteName, color.Reset)
			fmt.Println(URLWithUsername)
			if FalsePositive {
				FoundAccounts = append(FoundAccounts, websiteName+" - "+URLWithUsername+falsePositiveMessage)
			} else {
				FoundAccounts = append(FoundAccounts, websiteName+" - "+URLWithUsername)
			}
		}
	}
}

func websiteScrape(URLWithUsername string) string {
	res, err := doReq(URLWithUsername)
	if err != nil {
		fmt.Println(err)
		return ""
	}

	if res.code != 200 {
		fmt.Println("Unable to access website due to captcha/JavaScript/Cloudflare.")
	}

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(res.body))
	if err != nil {
		return ""
	}

	var websiteContent []string
	doc.Find("html").Each(func(index int, item *goquery.Selection) {
		websiteContent = append(websiteContent, item.Text())
	})

	return strings.Join(websiteContent, " ")
}

// URLWithUsername creates a URL with the username.
func URLWithUsername(websiteURL string, username string) string {
	return strings.Replace(websiteURL, "{}", username, -1)
}
