package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

func main() {
	file, err := os.Open("./resources/data.json")
	handleError(err)
	var result map[string]interface{}
	err = json.NewDecoder(file).Decode(&result)
	handleError(err)
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
	handleError(err)
	defer resp.Body.Close()
	handleError(err)
	if resp.StatusCode == 200 {
		fmt.Println(key)
		fmt.Println("FOUND -", url.(string)+username)
	}
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
