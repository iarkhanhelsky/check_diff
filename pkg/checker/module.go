package checker

import (
	"fmt"
	golangci_lint "github.com/iarkhanhelsky/check_diff/pkg/checker/golang/golangci-lint"
	"github.com/iarkhanhelsky/check_diff/pkg/checker/java/checkstyle"
	"github.com/iarkhanhelsky/check_diff/pkg/checker/k8s/kubelinter"
	"github.com/iarkhanhelsky/check_diff/pkg/checker/ruby/rubocop"
	"github.com/iarkhanhelsky/check_diff/pkg/core"
	"go.uber.org/config"
	"go.uber.org/fx"
)

type provider func(yaml *config.YAML) (core.Checker, error)

func defaultProvider(obj core.Checker) provider {
	return func(yaml *config.YAML) (core.Checker, error) {
		v := obj
		if err := yaml.Get(v.Tag()).Populate(v); err != nil {
			return nil, fmt.Errorf("can't create %s: %v", v.Tag(), err)
		}
		return v, nil
	}
}

var Module = ProvideCheckers(
	kubelinter.NewChecker(),
	checkstyle.NewChecker(),
	rubocop.NewChecker(),
	golangci_lint.NewChecker(),
)

func ProvideCheckers(checkers ...core.Checker) fx.Option {
	var annotated []interface{}
	for _, ch := range checkers {
		annotated = append(annotated, fx.Annotated{Group: "checkers", Target: defaultProvider(ch)})
	}

	return fx.Provide(annotated...)
}
