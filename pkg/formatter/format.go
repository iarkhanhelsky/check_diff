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
	switch string(format) {
	case "stdout":
		return stdout.NewFormatter(), nil
	case "gitlab":
		return &codeclimate.Formatter{}, nil
	case "codeclimate":
		return &codeclimate.Formatter{}, nil
	case "phabricator":
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
