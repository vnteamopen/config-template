package actions

import (
	"bufio"
	"fmt"
	"github.com/pkg/errors"
	"io/ioutil"
	"os"
	"strings"
	"testing"
)

func TestParse(t *testing.T) {
	sample1Path := "../test_samples/sample1.txt"
	sample2Path := "../test_samples/sample2.txt"

	testCases := []struct {
		name        string
		pattern     []string
		template    string
		samplePath  string
		checkResult func(expected, received string)
		checkError  func(err error)
	}{
		{
			name:       "Matched pattern - Simple template",
			template:   fmt.Sprintf(`{{file "%s"}}`, sample1Path),
			samplePath: sample1Path,
			checkResult: func(sampleContent, received string) {
				if sampleContent != received {
					t.Errorf("Wrong parse: \nExpected: %s\nReceived: %s", sampleContent, received)
				}
			},
			checkError: func(err error) {},
		},
		{
			name:       "Matched pattern - Template contains 2 begin parts",
			template:   fmt.Sprintf(`{{file {{file "%s"}}`, sample1Path),
			samplePath: sample1Path,
			checkResult: func(sampleContent, received string) {
				expected := fmt.Sprintf("{{file %s", sampleContent)
				if expected != received {
					t.Errorf("Wrong parse: \nExpected: %s\nReceived: %s", expected, received)
				}
			},
			checkError: func(err error) {},
		},
		{
			name:       "Matched pattern - Template contains 2 end parts",
			template:   fmt.Sprintf(`{{file "%s"}}"}}`, sample1Path),
			samplePath: sample1Path,
			checkResult: func(sampleContent, received string) {
				expected := fmt.Sprintf(`%s"}}`, sampleContent)
				if expected != received {
					t.Errorf("Wrong parse: \nExpected: %s\nReceived: %s", expected, received)
				}
			},
			checkError: func(err error) {},
		},
		{
			name:       "Matched pattern - Template contains begin part at then end",
			template:   fmt.Sprintf(`{{file "%s"}}{{file`, sample1Path),
			samplePath: sample1Path,
			checkResult: func(sampleContent, received string) {
				expected := fmt.Sprintf(`%s{{file`, sampleContent)
				if expected != received {
					t.Errorf("Wrong parse: \nExpected: %s\nReceived: %s", expected, received)
				}
			},
			checkError: func(err error) {},
		},
		{
			name:       "Matched pattern - Template contains another parts except template",
			template:   fmt.Sprintf(`abc {{file "%s"}} def`, sample1Path),
			samplePath: sample1Path,
			checkResult: func(sampleContent, received string) {
				expected := fmt.Sprintf(`abc %s def`, sampleContent)
				if expected != received {
					t.Errorf("Wrong parse: \nExpected: %s\nReceived: %s", expected, received)
				}
			},
			checkError: func(err error) {},
		},
		{
			name:       "Not match pattern - Simple template",
			template:   fmt.Sprintf(`{{filezilla "%s"}}`, sample1Path),
			samplePath: sample1Path,
			checkResult: func(sampleContent, received string) {
				expected := fmt.Sprintf(`{{filezilla "%s"}}`, sample1Path)
				if expected != received {
					t.Errorf("Wrong parse: \nExpected: %s\nReceived: %s", expected, received)
				}
			},
			checkError: func(err error) {},
		},
		{
			name:       "Not match pattern - Invalid file name in template",
			template:   fmt.Sprintf(`{{file "'%s'"}}`, sample1Path),
			samplePath: sample1Path,
			checkResult: func(sampleContent, received string) {
				expected := fmt.Sprintf(`{{file "'%s'"}}`, sample1Path)
				if expected != received {
					t.Errorf("Wrong parse: \nExpected: %s\nReceived: %s", expected, received)
				}
			},
			checkError: func(err error) {},
		},
		{
			name:        "Not match pattern - Notfound file name in template",
			template:    fmt.Sprintf(`{{file "%s"}}`, sample2Path),
			samplePath:  sample1Path,
			checkResult: func(sampleContent, received string) {},
			checkError: func(err error) {
				if !strings.Contains(err.Error(), "open") {
					t.Errorf("Wrong error: \n Expected: %+v\nReceived: %+v", errors.Wrap(err, "open"), err)
				}
			},
		},
		{
			name:       "Matched custom pattern - Simple template",
			pattern:    []string{"%", "%"},
			template:   fmt.Sprintf(`%%file "%s"%%`, sample1Path),
			samplePath: sample1Path,
			checkResult: func(sampleContent, received string) {
				if sampleContent != received {
					t.Errorf("Wrong parse: \nExpected: %s\nReceived: %s", sampleContent, received)
				}
			},
			checkError: func(err error) {},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			reader := strings.NewReader(tc.template)
			in := bufio.NewReader(reader)
			parser := NewSeqParser(tc.pattern)

			writer := new(strings.Builder)
			out := bufio.NewWriter(writer)

			err := parseInputToOutput(parser, in, []*bufio.Writer{out})
			if err != nil {
				tc.checkError(err)
				return
			}

			sample, err := os.Open(tc.samplePath)
			if err != nil {
				t.Errorf("Failed to open sample file: %v+", err.Error())
				return
			}
			sampleContent, err := ioutil.ReadAll(sample)
			if err != nil {
				t.Errorf("Failed to read sample content: %v+", err.Error())
				return
			}
			tc.checkResult(string(sampleContent), writer.String())
		})
	}
}
