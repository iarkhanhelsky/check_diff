package rubocop

import (
	"github.com/iarkhanhelsky/check_diff/pkg/core"
	"github.com/iarkhanhelsky/check_diff/pkg/downloader"
)

type Checker struct {
	core.Settings `yaml:",inline"`
}

var _ core.Checker = &Checker{}

func (checker *Checker) Tag() string {
	return "Rubocop"
}

func (checker Checker) Downloads() []downloader.Interface {
	return downloader.Empty
}

func (checker Checker) Check(ranges []core.LineRange) ([]core.Issue, error) {
	args := []string{"-f", "json"}

	if len(checker.Config) != 0 {
		args = append(args, "-c", checker.Config)
	}

	return core.NewFlow(checker.Tag(), checker.Settings,
		core.WithCommand(checker.CommandOrDefault("rubocop"), args...),
		core.WithFileExtensions(".rb", ".erb", "Rakefile", ".rake"),
		core.WithConverter(parseReport),
	).Run(ranges)
}
