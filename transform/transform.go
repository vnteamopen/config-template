package transform

import (
	"github.com/pkg/errors"
	"io/ioutil"
	"os"
	"regexp"
)

var fileNamePattern = `[a-zA-Z0-9_\.\-/\\ <>|:()&;]`

type Pattern struct {
	start rune
	end   rune
}

type transformer struct {
	begin      string
	beginIndex int
	end        string
	endIndex   int
	filePath   string
	regexFile  *regexp.Regexp
}

func NewTransformer() *transformer {
	r, _ := regexp.Compile(fileNamePattern)
	return &transformer{
		begin:      "{{file \"",
		beginIndex: 0,
		end:        "\"}}",
		endIndex:   0,
		regexFile:  r,
	}
}

func (t *transformer) Transform(input byte) ([]byte, error) {
	switch true {
	case t.isMatchedBegin(input):
		t.beginIndex += 1
		return nil, nil
	case t.isMatchedFileName(input):
		t.filePath += string(input)
		return nil, nil
	case t.isMatchedEnd(input):
		t.endIndex += 1
		if t.endIndex == len(t.end) {
			output, err := t.getTemplateContent()
			t.Reset()
			return output, err
		}
		return nil, nil
	default:
		output := t.begin[:t.beginIndex] + t.filePath + t.end[:t.endIndex]
		t.Reset()
		if t.isMatchedBegin(input) {
			t.beginIndex += 1
			return []byte(output), nil
		} else {
			return []byte(output + string(input)), nil
		}
	}
}

func (t *transformer) isMatchedBegin(input byte) bool {
	return !t.isEndBegin() && input == t.begin[t.beginIndex]
}

func (t *transformer) isMatchedFileName(input byte) bool {
	if t.isEndBegin() && t.endIndex == 0 {
		return t.regexFile.MatchString(string(input))
	}

	return false
}

func (t *transformer) isMatchedEnd(input byte) bool {
	if !t.isEndBegin() || len(t.filePath) == 0 || t.isEndEnd() {
		return false
	}

	return t.end[t.endIndex] == input
}

func (t *transformer) isEndBegin() bool {
	return t.beginIndex == len(t.begin)
}

func (t *transformer) isEndEnd() bool {
	return t.endIndex == len(t.end)
}

func (t *transformer) Reset() {
	t.beginIndex = 0
	t.endIndex = 0
	t.filePath = ""
}

func (t *transformer) getTemplateContent() ([]byte, error) {
	file, err := os.Open(t.filePath)
	if err != nil {
		return nil, errors.Wrap(err, "open")
	}
	defer file.Close()
	fileContent, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, errors.Wrap(err, "read file")
	}
	return fileContent, nil
}
