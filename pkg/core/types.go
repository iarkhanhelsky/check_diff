package core

import (
	"github.com/iarkhanhelsky/check_diff/pkg/tools"
	"io"
)

type Checker interface {
	Tag() string
	Check(ranges []LineRange) ([]Issue, error)
}

type Binaries interface {
	Binaries() []*tools.Binary
}

type LineRange struct {
	File  string `json:"file"`
	Start int    `json:"start"`
	End   int    `json:"end"`
}

type Issue struct {
	Tag      string
	File     string
	Line     int
	Column   int
	Severity string
	Message  string
	Source   string
}

type Formatter interface {
	Print(issues []Issue, w io.Writer) error
}
