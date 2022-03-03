package flows

import (
	"fmt"
	"regexp"
	"strings"
)

func Fuck(err error) {
	if err != nil {
		panic(err)
	}
}

func Slugify(s string) string {
	r := regexp.MustCompile(`(?m)[!@#$%^&\*\(\)\[\]\{\};:\,\./<>\?\|\x60~=_+ ]`)
	return r.ReplaceAllString(strings.ToLower(s), "-")
}

func ConfirmPrompt(format string, a ...interface{}) bool {
	fmt.Printf(format, a...)
	var confirm string
	fmt.Scanln(&confirm)
	if len(confirm) == 0 {
		return false
	} else {
		return strings.ToLower(confirm)[0] == 'y'
	}
}

func PrintJson(j []byte) {
	fmt.Println(string(j))
}
