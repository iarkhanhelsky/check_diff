package app

import (
	"context"
	"github.com/iarkhanhelsky/check_diff/pkg/checker/k8s/kubelinter"
	"github.com/iarkhanhelsky/check_diff/pkg/core"
	"go.uber.org/fx"
	"io"
	"io/ioutil"
	"os"
	"strings"
)

func Main() {
	var formatter core.Formatter
	var checker core.Checker
	var options Options
	var config Config

	app := fx.New(fx.Options(
		Module,
		kubelinter.Module,
		fx.Populate(&formatter, &checker, &options, &config),
	))
	app.Start(context.Background())
	var reader io.Reader
	if len(options.InputFile) == 0 {
		reader = os.Stdin
	} else {
		file, err := os.Open(options.InputFile)
		if err != nil {
			panic(err)
		}
		defer file.Close()
		reader = file
	}
	diff, _ := ioutil.ReadAll(reader)
	parser := core.NewDiffParser()
	for _, line := range strings.Split(string(diff), "\n") {
		parser.ParseNextLine(line)
	}
	for _, download := range checker.Downloads() {
		if err := download.Download(config.VendorDir); err != nil {
			panic(err)
		}
	}
	issues, err := checker.Check(parser.Result())
	if err != nil {
		panic(err)
	}

	var writer io.Writer
	if len(config.OutputFile) == 0 {
		writer = os.Stdout
	} else {
		file, err := os.Open(config.OutputFile)
		defer file.Close()
		if err != nil {
			panic(err)
		}
		writer = file
	}

	if err := formatter.Print(issues, writer); err != nil {
		panic(err)
	}
	app.Stop(context.Background())
}

func Run(lint core.Checker, formatter core.Formatter) error {
	return nil
}
