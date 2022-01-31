package checkstyle

import (
	"fmt"
	"github.com/iarkhanhelsky/check_diff/pkg/core"
	"github.com/iarkhanhelsky/check_diff/pkg/downloader"
	"go.uber.org/config"
)

type Checkstyle struct {
	core.Settings `yaml:",inline"`
}

func (j Checkstyle) Downloads() []downloader.Interface {
	return []downloader.Interface{
		downloader.NewHTTPDownloader(func(path string) error {
			return nil
		}, "checkstyle-all.jar",
			"ab3891e43b4bc41371d66b2ec615375d",
			"f97302b2d7f139a6cb0e9ebaa5142d61e96b2d438c0969d373729b88e95f5732",
			"https://github.com/checkstyle/checkstyle/releases/download/checkstyle-8.41/checkstyle-8.41-all.jar"),
	}
}

func (j Checkstyle) Check(ranges []core.LineRange) ([]core.Issue, error) {
	return []core.Issue{}, nil
}

var _ core.Checker = &Checkstyle{}

func NewCheckstyle(yaml *config.YAML) (core.Checker, error) {
	v := Checkstyle{}
	if err := yaml.Get("Checkstyle").Populate(&v); err != nil {
		return nil, fmt.Errorf("can't create Checkstyle: %v", err)
	}
	return &v, nil
}
