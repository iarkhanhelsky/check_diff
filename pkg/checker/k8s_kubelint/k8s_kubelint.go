package k8s_kubelint

import (
	"bytes"
	"errors"
	"github.com/iarkhanhelsky/check_diff/pkg/checker"
	"github.com/iarkhanhelsky/check_diff/pkg/core"
	"github.com/owenrumney/go-sarif/sarif"
	"os"
	"os/exec"
	"path"
)

type K8sKubeLint struct {
	kubeLint string
}

var _ checker.Checker = &K8sKubeLint{}

func (linter *K8sKubeLint) Setup() {

}

func (linter *K8sKubeLint) Check(ranges []core.LineRange) ([]core.Issue, error) {
	args := []string{"lint", "--format", "sarif"}
	for _, r := range ranges {
		args = append(args, r.File)
	}
	cmd := exec.Command(linter.kubeLint, args...)

	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Run()

	report, err := sarif.FromBytes(out.Bytes())
	if err != nil {
		return nil, err
	}

	return toIssues(report), nil
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

func toIssues(report *sarif.Report) []core.Issue {
	var issues []core.Issue
	for _, run := range report.Runs {
		for _, result := range run.Results {
			file, line, column := extractLocation(result.Locations)

			issues = append(issues, core.Issue{
				Tag:     *result.RuleID,
				Message: *result.Message.Text,
				File:    file,
				Line:    line,
				Column:  column,
			})
		}
	}

	return issues
}

func extractLocation(locations []*sarif.Location) (string, int, int) {
	if len(locations) == 0 {
		return "", 0, 0
	}

	location := locations[0].PhysicalLocation

	file := *location.ArtifactLocation.URI

	return file, 0, 0
}

func NewK8KubeLint() checker.Checker {
	return &K8sKubeLint{}
}
