package kubelinter

import (
	"github.com/iarkhanhelsky/check_diff/pkg/core"
	"github.com/iarkhanhelsky/check_diff/pkg/mapper"
	"github.com/iarkhanhelsky/check_diff/pkg/tools"
)

type Checker struct {
	core.Settings `yaml:",inline"`
	kubeLint      *tools.Binary
}

var _ core.Checker = &Checker{}
var _ core.Binaries = &Checker{}

func NewChecker() core.Checker {
	return &Checker{
		kubeLint: &tools.Binary{
			Name:    "kube-linter",
			Path:    "kube-linter",
			DstFile: "kube-linter.zip",
			Targets: map[tools.TargetTuple]tools.TargetSource{
				tools.DarwinAMD64: {
					Urls: []string{
						"https://github.com/stackrox/kube-linter/releases/download/0.2.5/kube-linter-darwin.zip",
					},
					MD5:    "58b4a9b8d55c1997c866471c14bbcb3a",
					SHA256: "dd75ba0a35db6ee12f36e8e36dac0e3e361e9a43103196962da86458092f9ab7",
				},
				tools.DarwinARM64: {
					Urls: []string{
						"https://github.com/stackrox/kube-linter/releases/download/0.2.5/kube-linter-darwin.zip",
					},
					MD5:    "58b4a9b8d55c1997c866471c14bbcb3a",
					SHA256: "dd75ba0a35db6ee12f36e8e36dac0e3e361e9a43103196962da86458092f9ab7",
				},
				tools.LinuxAMD64: {
					Urls: []string{
						"https://github.com/stackrox/kube-linter/releases/download/0.2.5/kube-linter-linux.zip",
					},
					MD5:    "05c8a6c57cb6d84ebae6a09efc9f46c2",
					SHA256: "a858572d7b673574855ce8cb84476b6d4690e79905d3fbaf303fe0a70eb8798e",
				},
			},
		},
	}
}

func (checker *Checker) Check(ranges []core.LineRange) ([]core.Issue, error) {
	args := []string{"lint", "--format", "sarif"}

	if len(checker.Config) != 0 {
		args = append(args, "--config", checker.Config)
	}

	return core.NewFlow(checker.Tag(), checker.Settings,
		core.WithCommand(checker.CommandOrDefault(checker.kubeLint.Executable()), args...),
		core.WithFileExtensions(".yaml", ".yml"),
		core.WithConverter(mapper.SarifBytesToIssues),
	).Run(ranges)
}

func (checker *Checker) Tag() string {
	return "KubeLinter"
}

func (checker *Checker) Binaries() []*tools.Binary {
	return []*tools.Binary{checker.kubeLint}
}
