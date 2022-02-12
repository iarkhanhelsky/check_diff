package checkstyle

import (
	"github.com/iarkhanhelsky/check_diff/pkg/core"
	"github.com/iarkhanhelsky/check_diff/pkg/downloader"
	"github.com/iarkhanhelsky/check_diff/pkg/mapper"
	"path"
)

type Checker struct {
	core.Settings `yaml:",inline"`

	checkstyle string
}

func (checker *Checker) Tag() string {
	return "Checkstyle"
}

func (checker *Checker) Check(ranges []core.LineRange) ([]core.Issue, error) {
	// java -jar .check_diff/vendor/checkstyle-all.jar \
	//  -c java/google_checks.xml \
	//  -f sarif \
	// ...files
	args := []string{"-jar", checker.checkstyle, "-c", checker.Config, "-f", "sarif"}
	return core.NewFlow("checkstyle", checker.Settings,
		core.WithCommand("java", args...),
		core.WithFileExtensions(".java"),
		core.WithConverter(mapper.SarifBytesToIssues),
	).Run(ranges)
}

func (checker *Checker) Downloads() []downloader.Interface {
	return []downloader.Interface{
		downloader.NewHTTPDownloader(checker.handleDownload, "checkstyle-all.jar",
			"970092a4271e5388b13055db1df485dd",
			"02ad3307e46059a7c4f8af6c5f61f477bc5fd910e56afb145c45904c95d213ac",
			"https://github.com/checkstyle/checkstyle/releases/download/checkstyle-9.3/checkstyle-9.3-all.jar"),
	}
}

func (checker *Checker) handleDownload(p string) error {
	checker.checkstyle = path.Join(p, "checkstyle-all.jar")
	return nil
}

var _ core.Checker = &Checker{}
