package kubelinter

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/iarkhanhelsky/check_diff/pkg/core"
	"os"
	"os/exec"
	"path"
)

var defaultCliArgs = []string{"lint", "--format", "sarif"}

type KubeLinter struct {
	kubeLint string
	settings Settings
}

var _ core.Checker = &KubeLinter{}

func (linter *KubeLinter) Setup() {

}

func (linter *KubeLinter) Check(ranges []core.LineRange) ([]core.Issue, error) {
	args := append(make([]string, 0), defaultCliArgs...)

	matchedRanges := linter.settings.Filter(ranges, ".yaml", ".yml")
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

	issues, err := parseSarif(stdout.Bytes())
	if err != nil {
		return nil, fmt.Errorf("failed to convert issues: %v", err)
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

func NewKubeLint(settings Settings) core.Checker {
	return &KubeLinter{settings: settings}
}
