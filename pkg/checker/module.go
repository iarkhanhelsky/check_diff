package checker

import (
	"github.com/iarkhanhelsky/check_diff/pkg/checker/java/checkstyle"
	"github.com/iarkhanhelsky/check_diff/pkg/checker/k8s/kubelinter"
	"github.com/iarkhanhelsky/check_diff/pkg/checker/ruby/rubocop"
	"github.com/iarkhanhelsky/check_diff/pkg/core"
	"go.uber.org/config"
	"go.uber.org/fx"
)

type provider func(yaml *config.YAML) (core.Checker, error)

var Module = ProvideCheckers(
	kubelinter.NewKubeLint,
	checkstyle.NewCheckstyle,
	rubocop.NewRubocop,
)

func ProvideCheckers(checkers ...provider) fx.Option {
	var annotated []interface{}
	for _, ch := range checkers {
		annotated = append(annotated, fx.Annotated{Group: "checkers", Target: ch})
	}

	return fx.Provide(annotated...)
}