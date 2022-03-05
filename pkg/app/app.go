package app

import (
	"context"
	"github.com/iarkhanhelsky/check_diff/pkg/app/command"
	"go.uber.org/fx"
)

func Main(env command.Env) error {
	var check command.Command
	app := fx.New(
		fx.Provide(func() command.Env { return env }),
		fx.Options(Module, fx.Populate(&check)))

	if err := app.Start(context.Background()); err != nil {
		return err
	}
	if err := check.Run(); err != nil {
		return err
	}
	if err := app.Stop(context.Background()); err != nil {
		return err
	}

	return nil
}
