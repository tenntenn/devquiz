package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/pkg/errors"
)

type CLI struct {
	stdout io.Writer
	stderr io.Writer

	// args
	csvpath string

	// options
	urlRow      int
	nameRow     int
	categoryRow int
	withHeader  bool
}

func NewCLI(args []string) (*CLI, error) {
	var cli CLI

	if len(args) <= 0 {
		return nil, errors.New("path for CSV file must be specified")
	}

	flagset := flag.NewFlagSet("devquiz", flag.ContinueOnError)
	flagset.BoolVar(&cli.withHeader, "withheader", false, "with header")
	flagset.IntVar(&cli.urlRow, "urlrow", 6, "URL row")
	flagset.IntVar(&cli.nameRow, "namerow", 1, "name row")
	flagset.IntVar(&cli.categoryRow, "categoryrow", 0, "category row")
	if err := flagset.Parse(args); err != nil {
		return nil, errors.Wrap(err, "NewCLI")
	}
	cli.csvpath = flagset.Arg(0)

	return &cli, nil
}

func (cli *CLI) Stdout() io.Writer {
	if cli.stdout == nil {
		return os.Stdout
	}
	return cli.stdout
}

func (cli *CLI) Stderr() io.Writer {
	if cli.stderr == nil {
		return os.Stderr
	}
	return cli.stderr
}

func (cli *CLI) Run() error {
	r, err := os.Open(cli.csvpath)
	if err != nil {
		return errors.Wrap(err, "could not open csv file")
	}

	cr := csv.NewReader(r)
	if !cli.withHeader {
		cr.Read() // skip
	}

	v := NewValidator(cr, cli.urlRow, cli.nameRow, cli.categoryRow)

LOOP:
	for i := 1; v.Next(); i++ {

		if err := v.Err(); err != nil {
			fmt.Printf("%d %s failed(%q)\n", i, v.Name(), err)
			continue LOOP
		}

		if v.Result() == nil {
			fmt.Println(i, "skip")
			continue
		}

		if errs := v.Result().Errors; errs != "" {
			fmt.Printf("%d %s failed(%q)\n", i, v.Name(), v.Result().Errors)
			continue
		}

		for _, e := range v.Result().Events {
			if e.Kind == "stderr" {
				fmt.Printf("%d %s failed(%q)\n", i, v.Name(), e.Message)
				continue LOOP
			}
			if strings.Contains(e.Message, "FAIL") {
				fmt.Printf("%d %s failed(%q)\n", i, v.Name(), e.Message)
				continue LOOP
			}

		}

		fmt.Println(i, v.Name(), "ok")
	}

	if err := v.Err(); err != nil {
		return err
	}

	return nil
}
