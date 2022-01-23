package formatter

import (
	"github.com/iarkhanhelsky/check_diff/pkg/core"
	"io"
)

type Formatter interface {
	Supports() []Format
	Print(issues []core.Issue, w io.Writer) error
}

//func Formatters() []Formatter {
//	return []Formatter{
//		&codeclimate.Formatter{},
//		&phabricator.Formatter{},
//		&stdout.Formatter{},
//	}
//}

func Formats() []Format {
	return []Format{
		STDOUT, Codeclimate, Gitlab, Phabricator,
	}
}
