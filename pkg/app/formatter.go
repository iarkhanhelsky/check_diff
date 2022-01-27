package app

import (
	"github.com/iarkhanhelsky/check_diff/pkg/core"
	"github.com/iarkhanhelsky/check_diff/pkg/formatter/codeclimate"
	"github.com/iarkhanhelsky/check_diff/pkg/formatter/phabricator"
	"github.com/iarkhanhelsky/check_diff/pkg/formatter/stdout"
)

func newFormatter(options Options) core.Formatter {
	switch options.Format {
	case "stdout":
		return stdout.NewFormatter()
	case "gitlab":
		return &codeclimate.Formatter{}
	case "phabricator":
		return &phabricator.Formatter{}
	}

	return nil
}
