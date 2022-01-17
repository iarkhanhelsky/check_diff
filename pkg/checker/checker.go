package checker

import "github.com/iarkhanhelsky/check_diff/pkg/core"

type Checker interface {
	Downloads() []core.Downloader
	Setup()
	Check(ranges []core.LineRange) ([]core.Issue, error)
}
