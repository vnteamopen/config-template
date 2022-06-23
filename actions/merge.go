package actions

import (
	"errors"
	"fmt"
	"os"
)

func Merge(inputPath string, listOutputPath []string) error {
	// printInfo(inputPath, outputPath)
	if _, err := os.Stat(inputPath); errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("input '%s' doesn't exist", inputPath)
	}
	outputContent, err := parse(inputPath)
	if err != nil {
		return err
	}

	for _, outputPath := range listOutputPath {
		if err := write(outputPath, outputContent); err != nil {
			return err
		}
	}
	return nil
}

func printInfo(templatePath, outputPath string) {
	fmt.Printf("* Template path: %s\n* Output path: %s\n", templatePath, outputPath)
}

func write(path, content string) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.WriteString(content)
	if err != nil {
		return err
	}
	return nil
}
