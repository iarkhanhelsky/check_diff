package checkstyle

import (
	"github.com/iarkhanhelsky/check_diff/pkg/core"
	"github.com/iarkhanhelsky/check_diff/pkg/mapper"
	"github.com/iarkhanhelsky/check_diff/pkg/tools"
	"strings"
)

type Checker struct {
	core.Settings `yaml:",inline"`

	checkstyleJar *tools.Binary
	checkstyle    string
}

var _ core.Checker = &Checker{}
var _ core.Binaries = &Checker{}

const (
	checkstyleVersion     = "9.3"
	checkstyleUrlTemplate = "https://github.com/checkstyle/checkstyle/releases/download/checkstyle-@VERSION/checkstyle-@VERSION-all.jar"
)

func NewChecker() core.Checker {
	return &Checker{
		checkstyleJar: &tools.Binary{
			Name:    "checkstyle",
			DstFile: "checkstyle-all.jar",
			Path:    "checkstyle-all.jar",
			Targets: map[tools.TargetTuple]tools.TargetSource{
				tools.Any: {
					Urls:   []string{strings.ReplaceAll(checkstyleUrlTemplate, "@VERSION", checkstyleVersion)},
					MD5:    "970092a4271e5388b13055db1df485dd",
					SHA256: "02ad3307e46059a7c4f8af6c5f61f477bc5fd910e56afb145c45904c95d213ac",
				},
			},
		},
	}
}

func (checker *Checker) Tag() string {
	return "Checkstyle"
}

func (checker *Checker) Check(ranges []core.LineRange) ([]core.Issue, error) {
	// java -jar .check_diff/vendor/checkstyle-all.jar \
	//  -c java/google_checks.xml \
	//  -f sarif \
	// ...files
	args := []string{"-jar", checker.checkstyleJar.Executable(), "-c", checker.Config, "-f", "sarif"}
	return core.NewFlow(checker.Tag(), checker.Settings,
		core.WithCommand("java", args...),
		core.WithFileExtensions(".java"),
		core.WithConverter(mapper.SarifBytesToIssues),
	).Run(ranges)
}

func (checker *Checker) Binaries() []*tools.Binary {
	return []*tools.Binary{checker.checkstyleJar}
}
