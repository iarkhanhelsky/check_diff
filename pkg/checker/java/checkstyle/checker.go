package checkstyle

import (
	"bytes"
	"fmt"
	"github.com/iarkhanhelsky/check_diff/pkg/core"
	"github.com/iarkhanhelsky/check_diff/pkg/downloader"
	"github.com/iarkhanhelsky/check_diff/pkg/mapper"
	"go.uber.org/config"
	"os/exec"
	"path"
)

type Checkstyle struct {
	core.Settings `yaml:",inline"`

	checkstyle string
}

func (j *Checkstyle) Downloads() []downloader.Interface {
	return []downloader.Interface{
		downloader.NewHTTPDownloader(j.handleDownload, "checkstyle-all.jar",
			"970092a4271e5388b13055db1df485dd",
			"02ad3307e46059a7c4f8af6c5f61f477bc5fd910e56afb145c45904c95d213ac",
			"https://github.com/checkstyle/checkstyle/releases/download/checkstyle-9.3/checkstyle-9.3-all.jar"),
	}
}

func (j *Checkstyle) Check(ranges []core.LineRange) ([]core.Issue, error) {
	// java -jar .check_diff/vendor/checkstyle-all.jar -c java/google_checks.xml java -f sarif
	args := []string{"-jar", j.checkstyle, "-c", j.Config, "-f", "sarif"}
	matchedRanges := j.Filter(ranges, ".java")
	if len(matchedRanges) == 0 {
		return []core.Issue{}, nil
	}
	for _, r := range matchedRanges {
		args = append(args, r.File)
	}

	cmd := exec.Command("java", args...)

	var stdout, stderr bytes.Buffer

	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil || cmd.ProcessState.ExitCode() > 1 {
		return nil, fmt.Errorf("failed to run checkstyle: %v: %s", err, string(stderr.Bytes()))
	}

	issues, err := mapper.SarifBytesToIssues(stdout.Bytes())
	if err != nil {
		return nil, fmt.Errorf("failed to convert checkstyle issues: %v", err)
	}
	return issues, nil

	return []core.Issue{}, nil
}

func (j *Checkstyle) handleDownload(p string) error {
	j.checkstyle = path.Join(p, "checkstyle-all.jar")
	return nil
}

var _ core.Checker = &Checkstyle{}

func NewCheckstyle(yaml *config.YAML) (core.Checker, error) {
	v := Checkstyle{}
	if err := yaml.Get("Checkstyle").Populate(&v); err != nil {
		return nil, fmt.Errorf("can't create Checkstyle: %v", err)
	}
	return &v, nil
}
