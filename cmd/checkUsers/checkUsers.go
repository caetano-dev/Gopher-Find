package checkUsers

import (
	"errors"
	"fmt"
	"gopher-find/cmd/color"
	"gopher-find/cmd/models"
	"io"
	"math/rand"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

var (
	FoundAccounts []string
	httpClient    = http.Client{Timeout: 30 * time.Second}
	userAgents    = []string{
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36",
		"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:89.0) Gecko/20100101 Firefox/89.0",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10.15; rv:89.0) Gecko/20100101 Firefox/89.0",
	}
	errorIndicators = []string{
		"404",
		"not found",
		"doesn't exist",
		"page not found",
		"user not found",
		"profile not found",
		"account not found",
	}
)

type Response struct {
	code int
	body string
}

func getRandomUserAgent() string {
	return userAgents[rand.Intn(len(userAgents))]
}

func validateResponse(res Response, parameter models.Parameter) bool {
	if res.code != 200 {
		return false
	}

	// Check minimum content length
	if len(res.body) < 100 {
		return false
	}

	content := strings.ToLower(res.body)

	// Check for error indicators
	for _, indicator := range errorIndicators {
		if strings.Contains(content, indicator) {
			return false
		}
	}

	// Check for delayed redirects
	if checkForDelayedRedirect(res) {
		return false
	}

	// If the site has a claimed username example, check if the response matches that pattern
	if parameter.UsernameClaimed != "" {
		claimedRes, err := doReq(strings.Replace(parameter.URL, "{}", parameter.UsernameClaimed, -1))
		if err == nil && claimedRes.code == 200 {
			// If the current response is significantly different from a known valid profile,
			// it might be a false positive
			if len(res.body) < len(claimedRes.body)/2 {
				return false
			}
		}
	}

	return true
}

func doReq(urlStr string) (Response, error) {
	_, err := url.Parse(urlStr)
	if err != nil {
		return Response{}, err
	}

	req, err := http.NewRequest("GET", urlStr, nil)
	if err != nil {
		return Response{}, err
	}

	req.Header.Set("User-Agent", getRandomUserAgent())

	resp, err := httpClient.Do(req)
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

func CheckIfUserExistsByErrorMessage(websiteName string, URLWithUsername string, parameter models.Parameter) {
	res, err := doReq(URLWithUsername)
	if err != nil {
		fmt.Println(color.Red+"[-] ERROR -", websiteName, "-", err, color.Reset)
		return
	}

	// Check for any of the error messages
	found := false
	for _, errorMsg := range parameter.GetErrorMessages() {
		if strings.Contains(res.body, errorMsg) {
			found = true
			break
		}
	}

	if found {
		fmt.Println(color.Red+"[-] NOT FOUND -", websiteName, color.Reset)
	} else {
		fmt.Println(color.Green+"[+] FOUND -", websiteName, color.Reset)
		fmt.Println(URLWithUsername)
		FoundAccounts = append(FoundAccounts, websiteName+" - "+URLWithUsername)
	}
}

func CheckIfUserExistsByStatusCode(websiteName string, URLWithUsername string, parameter models.Parameter) {
	res, err := doReq(URLWithUsername)
	if err != nil {
		fmt.Println(color.Red+"[-] ERROR -", websiteName, "-", err, color.Reset)
		return
	}

	if !validateResponse(res, parameter) {
		fmt.Println(color.Red+"[-] NOT FOUND -", websiteName, color.Reset)
		return
	}

	// Check if it might be a false positive
	if checkForDelayedRedirect(res) {
		fmt.Println(color.Yellow+"[?] POSSIBLE FALSE POSITIVE -", websiteName, "(Delayed Redirect)", color.Reset)
		fmt.Println(URLWithUsername)
		FoundAccounts = append(FoundAccounts, websiteName+" - "+URLWithUsername+" - Possible false positive (Delayed Redirect)")
		return
	}

	fmt.Println(color.Green+"[+] FOUND -", websiteName, color.Reset)
	fmt.Println(URLWithUsername)
	FoundAccounts = append(FoundAccounts, websiteName+" - "+URLWithUsername)
}

// CheckIfUserExistsByRedirect check if the user exists by redirect.
func CheckIfUserExistsByRedirect(websiteName string, URLWithUsername string, parameter models.Parameter) {
	// First make a normal request to check for delayed redirects in the content
	res, err := doReq(URLWithUsername)
	if err != nil {
		fmt.Println(color.Red+"[-] ERROR -", websiteName, "-", err, color.Reset)
		return
	}

	// Check for client-side redirects
	if checkForDelayedRedirect(res) {
		fmt.Println(color.Yellow+"[?] POSSIBLE FALSE POSITIVE -", websiteName, "(Delayed Redirect)", color.Reset)
		fmt.Println(URLWithUsername)
		FoundAccounts = append(FoundAccounts, websiteName+" - "+URLWithUsername+" - Possible false positive (Delayed Redirect)")
		return
	}

	// Now check for immediate server-side redirects
	req, err := http.NewRequest("GET", URLWithUsername, nil)
	if err != nil {
		fmt.Println(color.Red+"[-] ERROR -", websiteName, "-", err, color.Reset)
		return
	}

	client := httpClient
	client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return errors.New("redirect")
	}

	response, err := client.Do(req)
	if err != nil {
		// Some websites might return an error when we prevent redirects
		if !strings.Contains(err.Error(), "redirect") {
			fmt.Println(color.Red+"[-] ERROR -", websiteName, "-", err, color.Reset)
		}
		return
	}
	defer response.Body.Close()

	// Check the status code
	if response.StatusCode == 302 || response.StatusCode == 301 {
		fmt.Println(color.Red+"[-] NOT FOUND -", websiteName, color.Reset)
	} else {
		// Additional validation
		body, err := io.ReadAll(response.Body)
		if err != nil {
			fmt.Println(color.Red+"[-] ERROR -", websiteName, "-", err, color.Reset)
			return
		}

		res := Response{
			code: response.StatusCode,
			body: string(body),
		}

		if !validateResponse(res, parameter) {
			fmt.Println(color.Red+"[-] NOT FOUND -", websiteName, color.Reset)
			return
		}

		fmt.Println(color.Green+"[+] FOUND -", websiteName, color.Reset)
		fmt.Println(URLWithUsername)
		FoundAccounts = append(FoundAccounts, websiteName+" - "+URLWithUsername)
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

func checkForDelayedRedirect(res Response) bool {
	content := strings.ToLower(res.body)

	// Common redirect patterns
	redirectPatterns := []string{
		"window.location",
		"document.location",
		"setTimeout",
		"<meta http-equiv=\"refresh\"",
		"window.navigate",
		".href",
		"history.pushState",
		"history.replaceState",
	}

	for _, pattern := range redirectPatterns {
		if strings.Contains(content, pattern) {
			return true
		}
	}

	return false
}
