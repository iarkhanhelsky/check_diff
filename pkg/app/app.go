package app

import (
	"context"
	"github.com/iarkhanhelsky/check_diff/pkg/app/command"
	"github.com/iarkhanhelsky/check_diff/pkg/checker"
	"github.com/iarkhanhelsky/check_diff/pkg/formatter"
	"go.uber.org/fx"
)

func Main() {
	var check command.Check
	app := fx.New(fx.Options(
		Module,
		checker.Module,
		formatter.Module,
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
