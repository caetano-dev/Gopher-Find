package main

import (
	"io/ioutil"
	"os"
	"testing"
)

func TestUrLWithUsername(t *testing.T) {
	url := "http://www.example.com/{}"
	username := "test"
	expected := "http://www.example.com/test"
	actual := urlWithUsername(url, username)
	if actual != expected {
		t.Errorf("Expected %s, got %s", expected, actual)
	}
}

func TestGenerateFileWithFoundAcconts(t *testing.T) {
	foundAccounts = []string{"https://www.google.com", "https://www.facebook.com"}
	generateFileWithFoundAcconts(foundAccounts, "test")
	expected := "https://www.google.com\nhttps://www.facebook.com\n"
	actual, err := ioutil.ReadFile("./test.txt")
	handleError(err)
	if string(actual) != expected {
		panic("expected: " + expected + " actual: " + string(actual))
	}
	err = os.Remove("./test.txt")
	handleError(err)
}
