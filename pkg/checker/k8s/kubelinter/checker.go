package kubelinter

import (
	"errors"
	"fmt"
	"github.com/iarkhanhelsky/check_diff/pkg/core"
	"github.com/iarkhanhelsky/check_diff/pkg/mapper"
	"go.uber.org/config"
	"os"
	"path"
)

var defaultCliArgs = []string{"lint", "--format", "sarif"}

type KubeLinter struct {
	core.Settings `yaml:",inline"`
	kubeLint      string
}

var _ core.Checker = &KubeLinter{}

func (checker *KubeLinter) Tag() string {
	return "KubeLinter"
}

func (checker *KubeLinter) Check(ranges []core.LineRange) ([]core.Issue, error) {
	return core.NewFlow("kube-linter", checker.Settings,
		core.WithCommand(checker.kubeLint, defaultCliArgs...),
		core.WithFileExtensions(".yaml", ".yml"),
		core.WithConverter(mapper.SarifBytesToIssues),
	).Run(ranges)
}

func (checker *KubeLinter) handleDownload(dstPath string) error {
	kubelint := path.Join(dstPath, "kube-linter", "kube-linter")
	if _, err := os.Stat(kubelint); errors.Is(err, os.ErrNotExist) {
		return err
	}

	if err := os.Chmod(kubelint, 0755); err != nil {
		return err
	}
	checker.kubeLint = kubelint
	return nil
}

func NewKubeLint(yaml *config.YAML) (core.Checker, error) {
	kubeLinter := KubeLinter{}
	if err := yaml.Get("KubeLinter").Populate(&kubeLinter); err != nil {
		return nil, fmt.Errorf("can't create kube-linter: %v", err)
	}
	return &kubeLinter, nil
}
