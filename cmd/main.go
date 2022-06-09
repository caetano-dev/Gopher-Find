package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	u "sherlockgo/utils"
)

func main() {
	file, err := os.Open("./resources/data.json")
	u.HandleError(err)
	var result map[string]interface{}
	err = json.NewDecoder(file).Decode(&result)
	u.HandleError(err)
	fmt.Print("type the username: ")
	username := getInput()

	for key, value := range result {
		//fmt.Println(key, value)
		url := value.(map[string]interface{})["url"]
		checkURL(url, key, username)
	}
	defer file.Close()
}

func checkURL(url interface{}, key interface{}, username string) {
	resp, err := http.Get(url.(string))
	u.HandleError(err)
	defer resp.Body.Close()
	u.HandleError(err)
	if resp.StatusCode == 200 {
		fmt.Println(key)
		fmt.Println("FOUND -", url.(string)+username)
	}
}

func getInput() string {
	var input string
	fmt.Scanln(&input)
	return input
}
