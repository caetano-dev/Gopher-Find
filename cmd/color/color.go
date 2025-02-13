package color

import "runtime"

// Reset is a reset code for the terminal.
var Reset = "\033[0m"

// Red is a red code for the terminal.
var Red = "\033[31m"

// Green is a green code for the terminal.
var Green = "\033[32m"

var Yellow = "\033[33m"

func init() {
	if runtime.GOOS == "windows" {
		Reset = ""
		Red = ""
		Green = ""
		Yellow = ""
	}
}
