package main

import (
	"os"
	"time"

	"github.com/urfave/cli/v2"
	"github.com/vnteamopen/config-template/actions"
)

const (
	//  Help template: cli.AppHelpTemplate
	name    = "Config template"
	version = "1.0.0"
)

type FlagName string

const (
	FlagOverwrite FlagName = "overwrite"
)

var Flags = []cli.Flag{
	&cli.BoolFlag{
		Name:     string(FlagOverwrite),
		Usage:    "-w",
		Required: false,
		Value:    false,
		Aliases:  []string{"w"},
	},
}

func main() {
	app := &cli.App{
		Name:     name,
		Version:  version,
		Compiled: time.Now(),
		Authors:  []*cli.Author{&cli.Author{Name: "https://vnteamopen.com"}},
		HelpName: "config-template",
		Usage:    "A tool to merge file's contents to a template. Embedded pattern is {{file \"\"}}",
		UsageText: `config-template /path/to/input/file /path/to/output/file
config-template -w /path/to/input/file
config-template help`,
		EnableBashCompletion: true,
		Flags:                Flags,
		Action:               Action,
	}
	app.Run(os.Args)
}

func Action(c *cli.Context) error {
	c.App.Setup()

	isOverwrite := c.Bool(string(FlagOverwrite))
	valid, requiredArgs := validArgs(c.NArg(), isOverwrite)
	if !valid {
		cli.ShowAppHelp(c)
		return cli.Exit("", 1)
	}
	templatePath, outputPaths := getPaths(c.Args(), requiredArgs, isOverwrite)

	if err := actions.CharByCharMerge(templatePath, outputPaths); err != nil {
		return cli.Exit(err.Error(), 1)
	}
	if isOverwrite {
		if err := actions.OverwriteInput(templatePath); err != nil {
			return cli.Exit(err.Error(), 1)
		}
	}

	return cli.Exit("", 0)
}

func validArgs(totalArgs int, isOverwrite bool) (valid bool, requiredArgs int) {
	requiredArgs = 2
	if isOverwrite {
		requiredArgs = 1
	}

	return totalArgs >= requiredArgs, requiredArgs
}

func getPaths(args cli.Args, requiredArgs int, isOverwrite bool) (templatePath string, listOutputPath []string) {
	totalArgs := args.Len()
	templatePath = args.Get(0)
	numberOfAdditionOutput := totalArgs - requiredArgs

	listOutputPath = make([]string, 0, totalArgs)
	if isOverwrite {
		listOutputPath = append(listOutputPath, actions.CreateTmpFile(templatePath))
	} else {
		firstOutPut := args.Get(requiredArgs - 1)
		listOutputPath = append(listOutputPath, firstOutPut)
	}

	for i := 0; i < numberOfAdditionOutput; i++ {
		listOutputPath = append(listOutputPath, args.Get(requiredArgs+i))
	}

	return templatePath, listOutputPath
}
