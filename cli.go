package main

import (
	"fmt"
	"io"
	"strings"

	"github.com/urfave/cli"
)

const (
	// ExitCodeOK ...
	ExitCodeOK    int = 0
	ExitCodeError int = 1
)

// CLI ...
type CLI struct {
	outStream io.Writer
	errStream io.Writer
}

// Run ...
func (c *CLI) Run(args []string) int {
	app := cli.NewApp()
	app.Name = "glr"
	app.Usage = "gitlab releaser"
	app.Action = func(ctx *cli.Context) error {
		if ctx.NArg() < 3 {
			return fmt.Errorf("Require three options. [project name] [tag] [dist dir]")
		}
		p := ctx.Args().Get(0)
		t := ctx.Args().Get(1)
		d := ctx.Args().Get(2)
		r := "HEAD"

		tag, err := getTag(ctx, p, t)
		if err != nil {
			return err
		}

		if tag != nil {
			err = deleteTag(ctx, p, t)
			if err != nil {
				return err
			}
		}

		files := []string{}
		resps, err := multiUploads(ctx, p, d)
		if err != nil {
			return err
		}
		for _, res := range resps {
			files = append(files, fmt.Sprintf("- %s", res.Markdown))
		}

		msg := fmt.Sprintf("%s Release", t)
		desc := strings.Join(files, "\n")
		tag, err = createTag(ctx, p, t, r, msg, desc)
		if err != nil {
			return err
		}

		fmt.Fprintf(c.outStream, "%s %s Released.\n", p, t)

		return nil
	}
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "token, t",
			Usage: "access token",
		},
	}

	err := app.Run(args)
	if err != nil {
		fmt.Fprintln(c.errStream, err)
		return ExitCodeError
	}

	return ExitCodeOK
}
