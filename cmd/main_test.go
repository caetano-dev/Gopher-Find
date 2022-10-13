package main

import (
	"io/ioutil"
	"os"
	"strings"
	"testing"
)

func TestGenerateFileWithFoundAcconts(t *testing.T) {
	foundAccounts := []string{"https://www.google.com", "https://www.facebook.com"}
	generateFileWithFoundAcconts(foundAccounts, "test")
	expected := "https://www.google.com" + "\n" + "https://www.facebook.com"
	actual, err := ioutil.ReadFile("./test.txt")
	handleError(err)
	if !strings.Contains(string(actual), expected) {
		t.Errorf("File doesn't contain expected string %s", expected)
	}
	err = os.Remove("./test.txt")
	handleError(err)
}
