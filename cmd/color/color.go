package color

import "runtime"

var Reset = "\033[0m"
var Red = "\033[31m"
var Green = "\033[32m"

func init() {
	if runtime.GOOS == "windows" {
		Reset = ""
		Red = ""
		Green = ""
	}
}
