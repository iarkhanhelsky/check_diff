package app

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/iarkhanhelsky/check_diff/pkg/app/command"
	"github.com/iarkhanhelsky/check_diff/pkg/checker"
	"github.com/iarkhanhelsky/check_diff/pkg/core"
	"github.com/iarkhanhelsky/check_diff/pkg/formatter"
	"github.com/iarkhanhelsky/check_diff/pkg/tools"
	"go.uber.org/config"
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
	"go.uber.org/zap"
)

var Module = fx.Options(
	checker.Module,
	formatter.Module,
	tools.Module,
	fx.WithLogger(func() fxevent.Logger { return fxevent.NopLogger }),
	fx.Provide(
		NewCliOptions,
		NewConfig,
		NewYaml,
		NewFormatterParams,
		NewToolsParams,
		NewLogger,
		func(logger *zap.Logger) *zap.SugaredLogger { return logger.Sugar() },
		command.NewCommand,
	),
	fx.Invoke(func(cfg core.Config) {
		color.NoColor = !cfg.WithColors()
	}),
)

func NewFormatterParams(config core.Config) formatter.Params {
	return formatter.Params{Format: config.OutputFormat}
}

func NewToolsParams(config core.Config, logger *zap.SugaredLogger) tools.Params {
	return tools.Params{VendorDir: config.VendorDir, Logger: logger}
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
	if cliOpts.NoColor != nil {
		v := !*cliOpts.NoColor
		cfg.Color = &v
	}

	cfg.InputFile = cliOpts.InputFile

	return cfg, err
}

func NewYaml(options CliOptions) (*config.YAML, error) {
	return config.NewYAML(config.File(options.ConfigFile))
}
