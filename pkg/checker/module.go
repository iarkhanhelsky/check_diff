package checker

import (
	"github.com/iarkhanhelsky/check_diff/pkg/checker/java/checkstyle"
	"github.com/iarkhanhelsky/check_diff/pkg/checker/k8s/kubelinter"
	"go.uber.org/fx"
)

var Module = fx.Options(
	fx.Provide(kubelinter.ReadSettings, fx.Annotated{Group: "checkers", Target: kubelinter.NewKubeLint}),
	fx.Provide(checkstyle.ReadSettings, fx.Annotated{Group: "checkers", Target: checkstyle.NewCheckstyle}),
)
