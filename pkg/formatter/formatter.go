package formatter

import (
	"github.com/iarkhanhelsky/check_diff/pkg/core"
	"github.com/iarkhanhelsky/check_diff/pkg/formatter/codeclimate"
	"github.com/iarkhanhelsky/check_diff/pkg/formatter/phabricator"
	"io"
)

type Formatter interface {
	Supports() []Format
	Print(issues []core.Issue, w io.Writer) error
}

func Formatters() []Formatter {
	return []Formatter{
		&codeclimate.Formatter{},
		&phabricator.Formatter{},
	}
}

func Formats() []Format {
	return []Format{
		STDOUT, Codeclimate, Gitlab, Phabricator,
	}
}
