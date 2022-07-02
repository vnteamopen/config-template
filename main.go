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
	FlagOverwrite      FlagName = "overwrite"
	FlagOutputToScreen FlagName = "out-screen"
)

var Flags = []cli.Flag{
	&cli.BoolFlag{
		Name:     string(FlagOverwrite),
		Usage:    "-w",
		Required: false,
		Value:    false,
		Aliases:  []string{"w"},
	},
	&cli.BoolFlag{
		Name:     string(FlagOutputToScreen),
		Usage:    "-out-screen",
		Required: false,
		Value:    false,
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
	isOutputToScreen := c.Bool(string(FlagOutputToScreen))
	if valid := validArgs(c.NArg(), isOverwrite, isOutputToScreen); !valid {
		cli.ShowAppHelp(c)
		return cli.Exit("", 1)
	}
	templatePath, outputPaths := getPaths(c.Args(), isOverwrite)

	if err := actions.CharByCharMerge(templatePath, outputPaths, isOutputToScreen); err != nil {
		return cli.Exit(err.Error(), 1)
	}

	if isOverwrite {
		if err := actions.OverwriteInput(templatePath); err != nil {
			return cli.Exit(err.Error(), 1)
		}
	}

	return cli.Exit("", 0)
}

func validArgs(totalArgs int, isOverwrite, isOutputToScreen bool) bool {
	requiredArgs := 2
	if isOverwrite || isOutputToScreen {
		requiredArgs = 1
	}

	return totalArgs >= requiredArgs
}

func getPaths(args cli.Args, isOverwrite bool) (templatePath string, listOutputPath []string) {
	templatePath = args.Get(0)
	noOutputs := args.Len() - 1

	listOutputPath = make([]string, 0, args.Len())
	if isOverwrite {
		listOutputPath = append(listOutputPath, actions.CreateTmpFile(templatePath))
	}
	for i := 0; i < noOutputs; i++ {
		listOutputPath = append(listOutputPath, args.Get(i+1))
	}
	return templatePath, listOutputPath
}
