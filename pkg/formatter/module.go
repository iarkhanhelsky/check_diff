package formatter

import (
	"fmt"
	"github.com/iarkhanhelsky/check_diff/pkg/core"
	"github.com/iarkhanhelsky/check_diff/pkg/formatter/codeclimate"
	"github.com/iarkhanhelsky/check_diff/pkg/formatter/phabricator"
	"github.com/iarkhanhelsky/check_diff/pkg/formatter/stdout"
	"go.uber.org/fx"
)

type Options struct {
	Format     string
	OutputFile string
}

var Module = fx.Options(fx.Provide(NewFormatter))

func NewFormatter(opts Options) (core.Formatter, error) {
	switch opts.Format {
	case "stdout":
		return stdout.NewFormatter(), nil
	case "gitlab":
		return &codeclimate.Formatter{}, nil
	case "codeclimate":
		return &codeclimate.Formatter{}, nil
	case "phabricator":
		return &phabricator.Formatter{}, nil
	}

	return nil, fmt.Errorf("unknown formatter type: '%s'", opts.Format)
}
