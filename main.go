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
config-template help`,
		EnableBashCompletion: true,
		Flags:                Flags,
		Action:               Action,
	}
	app.Run(os.Args)
}

func Action(c *cli.Context) error {
	c.App.Setup()

	//requiredArgs := 2
	//isOverwrite := c.Bool(string(FlagOverwrite))
	//if isOverwrite {
	//	requiredArgs = 1
	//}
	//
	//if c.NArg() < requiredArgs {
	//	cli.ShowAppHelp(c)
	//	return cli.Exit("", 0)
	//}
	//
	//templatePath := c.Args().Get(0)
	//
	//numberOfAdditionOutput := c.NArg() - requiredArgs
	//listOutputPath := []string{c.Args().Get(requiredArgs - 1)}
	//for i := 0; i < numberOfAdditionOutput; i++ {
	//	listOutputPath = append(listOutputPath, c.Args().Get(requiredArgs+i))
	//}

	templatePath := "./output/test"
	listOutputPath := []string{"./output/test1"}
	if err := actions.CharByCharMerge(templatePath, listOutputPath[0]); err != nil {
		return cli.Exit(err.Error(), 1)
	}

	return cli.Exit("", 0)
}
