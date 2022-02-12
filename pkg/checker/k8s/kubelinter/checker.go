package kubelinter

import (
	"errors"
	"github.com/iarkhanhelsky/check_diff/pkg/core"
	"github.com/iarkhanhelsky/check_diff/pkg/mapper"
	"os"
	"path"
	"runtime"
)

var defaultCliArgs = []string{"lint", "--format", "sarif"}

type Checker struct {
	core.Settings `yaml:",inline"`
	kubeLint      string
}

var _ core.Checker = &Checker{}

func (checker *Checker) Tag() string {
	return "KubeLinter"
}

func (checker *Checker) Check(ranges []core.LineRange) ([]core.Issue, error) {
	return core.NewFlow("kube-linter", checker.Settings,
		core.WithCommand(checker.CommandOrDefault(checker.kubeLint), defaultCliArgs...),
		core.WithFileExtensions(".yaml", ".yml"),
		core.WithConverter(mapper.SarifBytesToIssues),
	).Run(ranges)
}

func (checker *Checker) handleDownload(dstPath string) error {
	exe := "kube-linter"
	if runtime.GOOS == "windows" {
		exe = "kube-linter.exe"
	}
	kubelint := path.Join(dstPath, "kube-linter", exe)
	if _, err := os.Stat(kubelint); errors.Is(err, os.ErrNotExist) {
		return err
	}

	if err := os.Chmod(kubelint, 0755); err != nil {
		return err
	}
	checker.kubeLint = kubelint
	return nil
}
