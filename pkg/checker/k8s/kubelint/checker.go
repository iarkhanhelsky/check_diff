package kubelint

import (
	"bytes"
	"errors"
	"github.com/iarkhanhelsky/check_diff/pkg/core"
	"os"
	"os/exec"
	"path"
)

type K8sKubeLint struct {
	kubeLint string
}

var _ core.Checker = &K8sKubeLint{}

func (linter *K8sKubeLint) Setup() {

}

func (linter *K8sKubeLint) Check(ranges []core.LineRange) ([]core.Issue, error) {
	args := []string{"lint", "--format", "sarif", "k8s/deployment.yaml"}
	for _, r := range ranges {
		args = append(args, r.File)
	}
	cmd := exec.Command(linter.kubeLint, args...)

	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Run()

	return parseSarif(out.Bytes())
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

func NewK8KubeLint() core.Checker {
	return &K8sKubeLint{}
}
