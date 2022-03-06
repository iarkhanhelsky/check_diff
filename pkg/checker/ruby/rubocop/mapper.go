package rubocop

import (
	"encoding/json"
	"fmt"
	"github.com/iarkhanhelsky/check_diff/pkg/core"
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
