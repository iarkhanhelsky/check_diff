package app

import (
	"context"
	"github.com/iarkhanhelsky/check_diff/pkg/app/command"
	"github.com/iarkhanhelsky/check_diff/pkg/checker/k8s/kubelinter"
	"github.com/iarkhanhelsky/check_diff/pkg/core"
	"github.com/iarkhanhelsky/check_diff/pkg/formatter"
	"go.uber.org/fx"
)

func Main() {
	var check command.Check
	app := fx.New(fx.Options(
		Module,
		kubelinter.Module,
		formatter.Module,
		fx.Invoke(Downloads),
		fx.Populate(&check),
	))
	if err := app.Start(context.Background()); err != nil {
		panic(err)
	}
	if err := check.Run(); err != nil {
		panic(err)
	}
	if err := app.Stop(context.Background()); err != nil {
		panic(err)
	}
}

func Downloads(checker core.Checker, config core.Config) error {
	for _, download := range checker.Downloads() {
		if err := download.Download(config.VendorDir); err != nil {
			panic(err)
		}
	}

	return nil
}
