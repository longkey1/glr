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
	app.Version = "0.0.1"
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

		isForce := ctx.Bool("force")
		if tag != nil && isForce {
			err = deleteTag(ctx, p, t)
			if err != nil {
				return err
			}
			tag = nil
		}
		if tag != nil {
			fmt.Fprintf(c.outStream, "%s %s already released.\n", p, t)
			return nil
		}

		files := []string{}
		resps, err := multiUploads(ctx, p, d)
		if err != nil {
			return err
		}
		for _, res := range resps {
			files = append(files, fmt.Sprintf("- %s", res.Markdown))
		}

		msg := fmt.Sprintf("Release %s", t)
		desc := fmt.Sprintf("# Downlodas\n%s", strings.Join(files, "\n"))
		tag, err = createTag(ctx, p, t, r, msg, desc)
		if err != nil {
			return err
		}

		fmt.Fprintf(c.outStream, "%s %s released.\n", p, t)

		return nil
	}
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "token, t",
			Usage: "access token",
		},
		cli.BoolFlag{
			Name:  "force, f",
			Usage: "delete if it exists and release",
		},
	}

	err := app.Run(args)
	if err != nil {
		fmt.Fprintln(c.errStream, err)
		return ExitCodeError
	}

	return ExitCodeOK
}
