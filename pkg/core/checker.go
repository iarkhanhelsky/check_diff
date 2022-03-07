package core

import "github.com/iarkhanhelsky/check_diff/pkg/tools"

type Checker interface {
	Tag() string
	Check(ranges []LineRange) ([]Issue, error)
}

type Binaries interface {
	Binaries() []*tools.Binary
}
