package app

import (
	"fmt"
	"github.com/iarkhanhelsky/check_diff/pkg/app/command"
	"github.com/iarkhanhelsky/check_diff/pkg/checker"
	"github.com/iarkhanhelsky/check_diff/pkg/core"
	"github.com/iarkhanhelsky/check_diff/pkg/formatter"
	"go.uber.org/config"
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
	"os"
)

var Module = fx.Options(
	checker.Module,
	formatter.Module,
	fx.WithLogger(NewLogger),
	fx.Provide(
		NewCliOptions,
		NewConfig,
		NewYaml,
		NewDiff,
		NewFormatterOptions,
		command.NewCheck,
	),
)

func NewFormatterOptions(config core.Config) formatter.Options {
	return formatter.Options{Format: config.OutputFormat}
}

func NewConfig(cliOpts CliOptions, yaml *config.YAML) (core.Config, error) {
	cfg := core.NewDefaultConfig()
	err := yaml.Get("CheckDiff").Populate(&cfg)
	if err != nil {
		err = fmt.Errorf("can't parse application cfg: %v", err)
	}
	if len(cliOpts.VendorDir) != 0 {
		cfg.VendorDir = cliOpts.VendorDir
	}
	if len(cliOpts.OutputFile) != 0 {
		cfg.OutputFile = cliOpts.OutputFile
	}
	if len(cliOpts.Format) != 0 {
		cfg.OutputFormat = cliOpts.Format
	}
	return cfg, err
}

func NewLogger() fxevent.Logger {
	trace := false
	for _, arg := range os.Args {
		if arg == "--trace" {
			trace = true
			break
		}
	}
	if trace {
		return &fxevent.ConsoleLogger{W: os.Stderr}
	}

	return fxevent.NopLogger
}
