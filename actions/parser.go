package actions

import (
	"io/ioutil"
	"os"
	"regexp"
	"strings"

	"github.com/pkg/errors"
)

const (
	pattern = `{{file *"[a-zA-Z0-9_\.\-/\\ <>|:()&;]*" *}}`
)

// TODO: please convert it to use rune Walk and state checking
func parse(path string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", errors.Wrap(err, "open")
	}
	defer file.Close()

	b, err := ioutil.ReadAll(file)
	if err != nil {
		return "", errors.Wrap(err, "read file")
	}

	re := regexp.MustCompile(pattern)
	output := re.ReplaceAllStringFunc(string(b), func(pattern string) string {
		firstDoubleQuote := strings.Index(pattern, "\"")
		lastDoubleQuote := strings.LastIndex(pattern, "\"")
		if firstDoubleQuote == -1 || lastDoubleQuote == -1 || firstDoubleQuote+1 >= lastDoubleQuote {
			return pattern
		}
		path := pattern[firstDoubleQuote+1 : lastDoubleQuote]
		includedContent, err := parse(path)
		if err != nil {
			return pattern
		}

		// Remove end of file new line from POSIX standard: https://pubs.opengroup.org/onlinepubs/9699919799/basedefs/V1_chap03.html#tag_03_206
		if len(includedContent) > 0 || includedContent[len(includedContent)-1] == byte(10) {
			includedContent = includedContent[:len(includedContent)-1]
		}
		return includedContent
	})
	return output, nil
}
