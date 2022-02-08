package app

import (
	"context"
	"github.com/iarkhanhelsky/check_diff/pkg/app/command"
	"go.uber.org/fx"
)

func Main() {
	var check command.Check
	app := fx.New(fx.Options(Module, fx.Populate(&check)))
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
