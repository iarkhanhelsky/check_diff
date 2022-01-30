package checkstyle

import (
	"github.com/iarkhanhelsky/check_diff/pkg/core"
	"github.com/iarkhanhelsky/check_diff/pkg/downloader"
)

type Settings struct {
	core.Settings `yaml:"JavaCheckstyle"`
}

type JavaCheckstyle struct {
}

func (j JavaCheckstyle) Downloads() []downloader.Interface {
	return []downloader.Interface{
		downloader.NewHTTPDownloader(func(path string) error {
			return nil
		}, "checkstyle-all.jar",
			"ab3891e43b4bc41371d66b2ec615375d",
			"f97302b2d7f139a6cb0e9ebaa5142d61e96b2d438c0969d373729b88e95f5732",
			"https://github.com/checkstyle/checkstyle/releases/download/checkstyle-8.41/checkstyle-8.41-all.jar"),
	}
}

func (j JavaCheckstyle) Check(ranges []core.LineRange) ([]core.Issue, error) {
	return []core.Issue{}, nil
}

var _ core.Checker = &JavaCheckstyle{}
