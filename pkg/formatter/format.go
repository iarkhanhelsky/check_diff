package formatter

import (
	"fmt"
	"github.com/iarkhanhelsky/check_diff/pkg/core"
	"github.com/iarkhanhelsky/check_diff/pkg/formatter/codeclimate"
	"github.com/iarkhanhelsky/check_diff/pkg/formatter/phabricator"
	"github.com/iarkhanhelsky/check_diff/pkg/formatter/stdout"
)

type Format string

const (
	STDOUT      Format = "stdout"
	Codeclimate Format = "codeclimate"
	Gitlab      Format = "gitlab"
	Phabricator Format = "phabricator"
)

func (format Format) String() string {
	return string(format)
}

func (format Format) Formatter() (core.Formatter, error) {
	switch format {
	case STDOUT:
		return stdout.NewFormatter(), nil
	case Gitlab, Codeclimate:
		return &codeclimate.Formatter{}, nil
	case Phabricator:
		return &phabricator.Formatter{}, nil
	}

	return nil, fmt.Errorf("unknown formatter type: '%s'", format)
}

func Formats() []Format {
	return []Format{STDOUT, Codeclimate, Gitlab, Phabricator}
}

func FormatNames() []string {
	formats := Formats()
	var list = make([]string, len(formats))
	for i, f := range formats {
		list[i] = string(f)
	}
	return list
}
