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

func (linter *KubeLinter) Check(ranges []core.LineRange) ([]core.Issue, error) {
	return core.NewFlow("kube-linter", linter.Settings,
		core.WithCommand(linter.kubeLint, defaultCliArgs...),
		core.WithFileExtensions(".yaml", ".yml"),
		core.WithConverter(mapper.SarifBytesToIssues),
	).Run(ranges)
}

func (linter *KubeLinter) handleDownload(dstPath string) error {
	kubelint := path.Join(dstPath, "kube-linter", "kube-linter")
	if _, err := os.Stat(kubelint); errors.Is(err, os.ErrNotExist) {
		return err
	}

	if err := os.Chmod(kubelint, 0755); err != nil {
		return err
	}
	linter.kubeLint = kubelint
	return nil
}

func NewKubeLint(yaml *config.YAML) (core.Checker, error) {
	kubeLinter := KubeLinter{}
	if err := yaml.Get("KubeLinter").Populate(&kubeLinter); err != nil {
		return nil, fmt.Errorf("can't create kube-linter: %v", err)
	}
	return &kubeLinter, nil
}
