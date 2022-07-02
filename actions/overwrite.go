package actions

import (
	"io"
	"os"

	"github.com/pkg/errors"
)

func OverwriteInput(input string) error {
	inputFile, err := os.Open(CreateTmpFile(input))
	if err != nil {
		return errors.Wrap(err, "open")
	}
	defer inputFile.Close()

	outputFile, err := os.Create(input)
	if err != nil {
		return errors.Wrap(err, "open")
	}
	defer outputFile.Close()

	_, err = io.Copy(outputFile, inputFile)
	if err != nil {
		return err
	}

	return os.Remove(CreateTmpFile(input))
}

func CreateTmpFile(filePath string) string {
	return filePath + ".tmp"
}
