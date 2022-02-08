package checker

import (
	"github.com/iarkhanhelsky/check_diff/pkg/checker/java/checkstyle"
	"github.com/iarkhanhelsky/check_diff/pkg/checker/k8s/kubelinter"
	"github.com/iarkhanhelsky/check_diff/pkg/checker/ruby/rubocop"
	"go.uber.org/fx"
)

var Module = fx.Options(
	fx.Provide(fx.Annotated{Group: "checkers", Target: kubelinter.NewKubeLint}),
	fx.Provide(fx.Annotated{Group: "checkers", Target: checkstyle.NewCheckstyle}),
	fx.Provide(fx.Annotated{Group: "checkers", Target: rubocop.NewRubocop}),
)
