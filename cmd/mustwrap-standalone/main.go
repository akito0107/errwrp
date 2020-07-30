package main

import (
	"log"
	"os"
	"strings"

	"path/filepath"

	"fmt"

	"io"

	"bufio"

	"github.com/akito0107/errwrp"
	. "github.com/logrusorgru/aurora"
	"github.com/pkg/errors"
	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "mustwrap-standalone"
	app.Usage = "check if return err with no errors.Wrap(f)"
	app.UsageText = "mustwrap-standalone [OPTIONS]"
	app.Action = run
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "path, p",
			Usage: "check file name (or directory), acceptable for comma separated (required)",
		},
		cli.StringFlag{
			Name:  "exclude, e",
			Usage: "exclude file name (or directory), acceptable for comma separated (default=vendor)",
			Value: "vendor",
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func run(ctx *cli.Context) error {
	p := ctx.String("path")
	e := ctx.String("exclude")
	if p == "" {
		return errors.New("path is required")
	}
	fnames := strings.Split(p, ",")
	excludes := strings.Split(e, ",")
	var filepaths []string
	for _, fname := range fnames {
		err := filepath.Walk(fname, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return errors.Wrap(err, "inner filepath/Walk")
			}
			if strings.HasPrefix(path, ".") || !strings.HasSuffix(path, ".go") {
				return nil
			}
			for _, e := range excludes {
				if strings.HasPrefix(path, e) {
					return nil
				}
			}
			filepaths = append(filepaths, path)
			return nil
		})
		if err != nil {
			return errors.Wrap(err, "filepath/Walk")
		}
	}

	var results []*errwrp.Result
	for _, f := range filepaths {
		src, err := os.Open(f)
		if err != nil {
			return errors.Wrap(err, "run os/Open")
		}
		decls, fset, cm, err := errwrp.Parse(src, f)
		if err != nil {
			return errors.Wrap(err, "run errwrp/Parse")
		}
		for _, d := range decls {
			r, err := errwrp.Check(d, fset, cm)
			if err != nil {
				return errors.Wrap(err, "run errwrp/Check")
			}
			results = append(results, r...)
		}
	}
	if err := out(results); err != nil {
		return err
	}

	if len(results) > 0 {
		return errors.New("may contain unwrapped errors")
	}

	return nil
}

func out(results []*errwrp.Result) error {
	w := os.Stdout
	for _, r := range results {
		fmt.Fprintf(w, Sprintf(Red("file: %s line no: %d\n"), Red(r.Fname), Red(r.Position.Line)))
		f, err := os.Open(r.Fname)
		if err != nil {
			return errors.Wrap(err, "out")
		}
		printSpecificLine(w, f, r.Position.Line)
		fmt.Fprintf(w, "\n")
		fmt.Fprintf(w, "\n")
	}
	return nil
}

func printSpecificLine(w io.Writer, r io.Reader, lnum int) {
	var line int
	sc := bufio.NewScanner(r)
	for sc.Scan() {
		line++
		if line == lnum-1 || line == lnum+1 || line == lnum-2 || line == lnum+2 {
			fmt.Fprint(w, Sprintf(Gray("%d| %s\n"), Gray(line), Gray(sc.Text())))
		}
		if line == lnum {
			fmt.Fprint(w, Sprintf(Red("%d| %s\n"), Red(line).Bold(), Red(sc.Text()).Bold()))
		}
	}
}
