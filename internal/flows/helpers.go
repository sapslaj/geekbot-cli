package flows

import (
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
