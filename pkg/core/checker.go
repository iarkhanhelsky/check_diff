package core

import "github.com/iarkhanhelsky/check_diff/pkg/downloader"

type Checker interface {
	Downloads() []downloader.Interface
	Setup()
	Check(ranges []LineRange) ([]Issue, error)
}
