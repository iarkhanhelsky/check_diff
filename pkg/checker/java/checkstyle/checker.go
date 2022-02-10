package checkstyle

import (
	"fmt"
	"github.com/iarkhanhelsky/check_diff/pkg/core"
	"github.com/iarkhanhelsky/check_diff/pkg/downloader"
	"github.com/iarkhanhelsky/check_diff/pkg/mapper"
	"go.uber.org/config"
	"path"
)

type Checkstyle struct {
	core.Settings `yaml:",inline"`

	checkstyle string
}

func (linter *Checkstyle) Downloads() []downloader.Interface {
	return []downloader.Interface{
		downloader.NewHTTPDownloader(linter.handleDownload, "checkstyle-all.jar",
			"970092a4271e5388b13055db1df485dd",
			"02ad3307e46059a7c4f8af6c5f61f477bc5fd910e56afb145c45904c95d213ac",
			"https://github.com/checkstyle/checkstyle/releases/download/checkstyle-9.3/checkstyle-9.3-all.jar"),
	}
}

func (linter *Checkstyle) Check(ranges []core.LineRange) ([]core.Issue, error) {
	// java -jar .check_diff/vendor/checkstyle-all.jar \
	//  -c java/google_checks.xml \
	//  -f sarif \
	// ...files
	args := []string{"-jar", linter.checkstyle, "-c", linter.Config, "-f", "sarif"}
	return core.NewFlow("checkstyle", linter.Settings,
		core.WithCommand("java", args...),
		core.WithFileExtensions(".java"),
		core.WithConverter(mapper.SarifBytesToIssues),
	).Run(ranges)
}

func (linter *Checkstyle) handleDownload(p string) error {
	linter.checkstyle = path.Join(p, "checkstyle-all.jar")
	return nil
}

var _ core.Checker = &Checkstyle{}

func NewCheckstyle(yaml *config.YAML) (core.Checker, error) {
	v := Checkstyle{}
	if err := yaml.Get("Checkstyle").Populate(&v); err != nil {
		return nil, fmt.Errorf("can't create Checkstyle: %v", err)
	}
	return &v, nil
}
