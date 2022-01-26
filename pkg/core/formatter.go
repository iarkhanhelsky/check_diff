package core

import (
	"io"
)

type Formatter interface {
	Supports() []Format
	Print(issues []Issue, w io.Writer) error
}

func Formats() []Format {
	return []Format{
		STDOUT, Codeclimate, Gitlab, Phabricator,
	}
}
