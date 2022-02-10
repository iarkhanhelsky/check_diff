package rubocop

import (
	"encoding/json"
	"fmt"
	"github.com/iarkhanhelsky/check_diff/pkg/core"
	"github.com/iarkhanhelsky/check_diff/pkg/downloader"
	"go.uber.org/config"
)

type report struct {
	Files []file `json:"files"`
}

type file struct {
	Path     string    `json:"path"`
	Offences []offence `json:"offenses"`
}

func (f file) issues() []core.Issue {
	var issues []core.Issue
	for _, o := range f.Offences {
		issues = append(issues, core.Issue{
			File:     f.Path,
			Line:     o.Location.Start,
			Column:   o.Location.Column,
			Message:  o.Message,
			Severity: o.Severity,
			Source:   o.CopName,
		})
	}
	return issues
}

type offence struct {
	Location location `json:"location"`
	Severity string   `json:"severity"`
	Message  string   `json:"message"`
	CopName  string   `json:"cop_name"`
}

type location struct {
	Start  int `json:"start_line"`
	End    int `json:"last_line"`
	Column int `json:"start_column"`
}

type Rubocop struct {
	core.Settings `yaml:",inline"`
}

var _ core.Checker = &Rubocop{}

func NewRubocop(yaml *config.YAML) (core.Checker, error) {
	v := Rubocop{}
	if err := yaml.Get("Rubocop").Populate(&v); err != nil {
		return nil, fmt.Errorf("can't create Rubocop: %v", err)
	}
	return &v, nil
}

func (r Rubocop) Downloads() []downloader.Interface {
	return downloader.Empty
}

func (r Rubocop) Check(ranges []core.LineRange) ([]core.Issue, error) {
	command := r.Command
	if len(command) == 0 {
		command = "rubocop"
	}

	args := []string{"-f", "json"}

	if len(r.Config) != 0 {
		args = append(args, "-c", r.Config)
	}

	return core.NewFlow("rubocop", r.Settings,
		core.WithCommand(command, args...),
		core.WithFileExtensions(".rb", ".erb", "Rakefile", ".rake"),
		core.WithConverter(parseReport),
	).Run(ranges)
}

func parseReport(bytes []byte) ([]core.Issue, error) {
	r := report{}
	err := json.Unmarshal(bytes, &r)
	if err != nil {
		err = fmt.Errorf("failed to parse rubocop report: %v", err)
	}

	var issues []core.Issue
	for _, f := range r.Files {
		issues = append(issues, f.issues()...)
	}
	return issues, err
}
