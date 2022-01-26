package app

import (
	"context"
	"github.com/iarkhanhelsky/check_diff/pkg/checker/k8s/kubelint"
	"github.com/iarkhanhelsky/check_diff/pkg/core"
	"go.uber.org/fx"
)

func Main() {
	app := fx.New(fx.Options(
		Module,
		kubelint.Module,
		fx.Invoke(Run),
	))

	app.Start(context.Background())
	app.Stop(context.Background())
}

func Run(lint core.Checker, formatter core.Formatter) error {
	return nil
}
