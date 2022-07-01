package actions

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

type sequenceParser struct {
	begin      string
	beginIndex int
	end        string
	endIndex   int
	filePath   string
	regexFile  *regexp.Regexp
}

func NewSeqParser() *sequenceParser {
	r, _ := regexp.Compile(fileNamePattern)
	return &sequenceParser{
		begin:      "{{file \"",
		beginIndex: 0,
		end:        "\"}}",
		endIndex:   0,
		regexFile:  r,
	}
}

func (p *sequenceParser) Transform(input byte) ([]byte, error) {
	switch true {
	case p.isMatchedBegin(input):
		p.beginIndex += 1
		return nil, nil
	case p.isMatchedFileName(input):
		p.filePath += string(input)
		return nil, nil
	case p.isMatchedEnd(input):
		p.endIndex += 1
		if p.endIndex == len(p.end) {
			output, err := p.getTemplateContent()
			p.Reset()
			return output, err
		}
		return nil, nil
	default:
		output := p.begin[:p.beginIndex] + p.filePath + p.end[:p.endIndex]
		p.Reset()
		if p.isMatchedBegin(input) {
			p.beginIndex += 1
			return []byte(output), nil
		} else {
			return []byte(output + string(input)), nil
		}
	}
}

func (p *sequenceParser) isMatchedBegin(input byte) bool {
	return !p.isEndBegin() && input == p.begin[p.beginIndex]
}

func (p *sequenceParser) isMatchedFileName(input byte) bool {
	if p.isEndBegin() && p.endIndex == 0 {
		return p.regexFile.MatchString(string(input))
	}

	return false
}

func (p *sequenceParser) isMatchedEnd(input byte) bool {
	if !p.isEndBegin() || len(p.filePath) == 0 || p.isEndEnd() {
		return false
	}

	return p.end[p.endIndex] == input
}

func (p *sequenceParser) isEndBegin() bool {
	return p.beginIndex == len(p.begin)
}

func (p *sequenceParser) isEndEnd() bool {
	return p.endIndex == len(p.end)
}

func (p *sequenceParser) Reset() {
	p.beginIndex = 0
	p.endIndex = 0
	p.filePath = ""
}

func (p *sequenceParser) Flush() []byte {
	output := p.begin[:p.beginIndex] + p.filePath + p.end[:p.endIndex]
	p.Reset()
	return []byte(output)
}

func (p *sequenceParser) getTemplateContent() ([]byte, error) {
	file, err := os.Open(p.filePath)
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
