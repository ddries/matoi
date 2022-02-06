package util

import (
	"fmt"
	"os"
)

func ThrowError(text string) {
	fmt.Println("[error] " + text)
	os.Exit(1)
}

func Verbose(prefix string, text string) {
	fmt.Println("[" + prefix + "] " + text)
}