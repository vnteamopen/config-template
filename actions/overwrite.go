package actions

import (
	"io"
	"os"

	"github.com/pkg/errors"
)

type WriteTempRequest struct {
	InputPath      string
	Overwrite      bool
	OverwritePath  string
	OutputToScreen bool
}

func WriteTemp(r WriteTempRequest) error {
	outs := []io.Writer{}

	if r.Overwrite {
		outputFile, err := os.Create(r.OverwritePath)
		if err != nil {
			return errors.Wrap(err, "open")
		}
		defer outputFile.Close()
		outs = append(outs, outputFile)
	}

	if r.OutputToScreen {
		outs = append(outs, os.Stdout)
	}

	if len(outs) == 0 {
		return nil
	}

	tmpFile := CreateTmpFile(r.InputPath)
	in, err := os.Open(tmpFile)
	if err != nil {
		return errors.Wrap(err, "open")
	}
	defer in.Close()

	if _, err := io.Copy(io.MultiWriter(outs...), in); err != nil {
		return err
	}
	return os.Remove(tmpFile)
}

func CreateTmpFile(filePath string) string {
	return filePath + ".tmp"
}
