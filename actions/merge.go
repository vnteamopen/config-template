package actions

import (
	"bufio"
	"errors"
	"fmt"
	"io"
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

type MergeRequest struct {
	InputPath        string
	ListOutputPath   []string
	IsOutputToScreen bool
	Pattern          []string
}

func CharByCharMerge(req MergeRequest) error {
	inFile, err := getInputFile(req.InputPath)
	if err != nil {
		return err
	}
	defer closeFile(inFile)

	outFiles, err := getOutputFiles(req.ListOutputPath, req.IsOutputToScreen)
	defer closeFile(outFiles...)
	if err != nil {
		return err
	}

	inReader, outWriters := getInOutStreams(inFile, outFiles)

	parser := NewSeqParser(req.Pattern)
	return parseInputToOutput(parser, inReader, outWriters)
}

func getInputFile(path string) (*os.File, error) {
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		return nil, fmt.Errorf("input '%s' doesn't exist", path)
	}

	input, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("cannot open input file '%s': %s", path, err.Error())
	}

	return input, nil
}

func getOutputFiles(listPath []string, isOutputToScreen bool) ([]*os.File, error) {
	outputFiles := make([]*os.File, len(listPath))
	for i, path := range listPath {
		output, err := os.Create(path)
		if err != nil {
			return nil, fmt.Errorf("cannot create output file '%s': %s", path, err.Error())
		}
		outputFiles[i] = output
	}
	if isOutputToScreen {
		outputFiles = append(outputFiles, os.Stdout)
	}

	return outputFiles, nil
}

func getInOutStreams(inputFile *os.File, listOutputFile []*os.File) (*bufio.Reader, []*bufio.Writer) {
	inReader := bufio.NewReader(inputFile)
	outWriters := make([]*bufio.Writer, len(listOutputFile))
	for i, output := range listOutputFile {
		outWriters[i] = bufio.NewWriter(output)
	}
	return inReader, outWriters
}

func closeFile(files ...*os.File) {
	for _, file := range files {
		if err := file.Close(); err != nil {
			fmt.Printf("cannot close file: '%s'", err.Error())
		}
	}
}

func parseInputToOutput(parser Parser, inReader *bufio.Reader, outWriters []*bufio.Writer) error {
	buf := make([]byte, 1)
	for {
		n, err := inReader.Read(buf)
		if err != nil && err != io.EOF {
			return fmt.Errorf("cannot read file: %s", err.Error())
		}

		if err == io.EOF {
			break
		}

		if n == 0 {
			continue
		}

		transformerBuf, templateFilePath, err := parser.Transform(buf[0])
		if err != nil {
			return err
		}
		if len(templateFilePath) > 0 {
			templateFile, err := getInputFile(templateFilePath)
			if err != nil {
				return err
			}
			defer closeFile(templateFile)
			templateReader := bufio.NewReader(templateFile)
			if err := parseInputToOutput(parser, templateReader, outWriters); err != nil {
				return err
			}
		}
		if transformerBuf == nil {
			continue
		}

		for i := range outWriters {
			if _, err := outWriters[i].Write(transformerBuf); err != nil {
				return fmt.Errorf("cannot write file: %s", err.Error())
			}
		}
	}

	remainBuf := parser.Flush()
	for i := range outWriters {
		if _, err := outWriters[i].Write(remainBuf); err != nil {
			return fmt.Errorf("cannot write file: %s", err.Error())
		}
		if err := outWriters[i].Flush(); err != nil {
			return fmt.Errorf("cannot write file: %s", err.Error())
		}
	}

	return nil
}
