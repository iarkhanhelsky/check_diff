package kubelint

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/iarkhanhelsky/check_diff/pkg/core"
	"os"
	"os/exec"
	"path"
	"strings"
)

type K8sKubeLint struct {
	kubeLint string
	settings Settings
}

var _ core.Checker = &K8sKubeLint{}

func (linter *K8sKubeLint) Setup() {

}

func (linter *K8sKubeLint) Check(ranges []core.LineRange) ([]core.Issue, error) {
	args := []string{"lint", "--format", "sarif"}
	var hasInput bool
	for _, r := range ranges {
		if strings.HasSuffix(r.File, ".yaml") {
			args = append(args, r.File)
			hasInput = true
		}
	}

	if !hasInput {
		return []core.Issue{}, nil
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

func (linter *K8sKubeLint) handleDownload(dstPath string) error {
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

func NewK8KubeLint(settings Settings) core.Checker {
	return &K8sKubeLint{settings: settings}
}
