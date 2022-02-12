package core

import "github.com/iarkhanhelsky/check_diff/pkg/downloader"

type Checker interface {
	Tag() string
	Downloads() []downloader.Interface
	Check(ranges []LineRange) ([]Issue, error)
}
