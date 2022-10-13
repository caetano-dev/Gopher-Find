package checkUsers

import (
	"testing"
)

func TestUrLWithUsername(t *testing.T) {
	url := "http://www.example.com/{}"
	username := "test"
	expected := "http://www.example.com/test"
	actual := URLWithUsername(url, username)
	if actual != expected {
		t.Errorf("Expected %s, got %s", expected, actual)
	}
}
