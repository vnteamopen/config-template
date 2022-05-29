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
		Action:               Action,
	}
	app.Run(os.Args)
}

func Action(c *cli.Context) error {
	c.App.Setup()
	if c.NArg() <= 1 {
		cli.ShowAppHelp(c)
		return cli.Exit("", 0)
	}
	templatePath := c.Args().Get(0)
	outputPath := c.Args().Get(1)
	if err := actions.Merge(templatePath, outputPath); err != nil {
		return cli.Exit(err.Error(), 1)
	}
	return cli.Exit("", 0)
}
