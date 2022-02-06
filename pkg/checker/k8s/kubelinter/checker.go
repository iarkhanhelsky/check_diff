package kubelinter

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/iarkhanhelsky/check_diff/pkg/core"
	"github.com/iarkhanhelsky/check_diff/pkg/mapper"
	"go.uber.org/config"
	"os"
	"os/exec"
	"path"
)

var defaultCliArgs = []string{"lint", "--format", "sarif"}

type KubeLinter struct {
	core.Settings `yaml:",inline"`
	kubeLint      string
}

var _ core.Checker = &KubeLinter{}

func (linter *KubeLinter) Check(ranges []core.LineRange) ([]core.Issue, error) {
	args := append(make([]string, 0), defaultCliArgs...)

	matchedRanges := linter.Filter(ranges, ".yaml", ".yml")
	if len(matchedRanges) == 0 {
		return []core.Issue{}, nil
	}
	for _, r := range matchedRanges {
		args = append(args, r.File)
	}

	cmd := exec.Command(linter.kubeLint, args...)

	var stdout, stderr bytes.Buffer

	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil && cmd.ProcessState.ExitCode() != 1 {
		return nil, fmt.Errorf("failed to run kube-lint: %v: %s", err, string(stderr.Bytes()))
	}

	issues, err := mapper.SarifBytesToIssues(stdout.Bytes())
	if err != nil {
		return nil, fmt.Errorf("failed to convert kube-linter issues: %v", err)
	}
	return issues, nil
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
