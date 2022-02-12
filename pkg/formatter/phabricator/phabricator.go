package phabricator

import (
	"encoding/json"
	"io"

	"github.com/iarkhanhelsky/check_diff/pkg/core"
)

//   {
//  	"name": "Comment Name",
//  	"code": "Haiku", "severity": "error",
//    	"path": "example/path", "line": 2, "char": 0,
//    	"description": "Line is not a Haiku"
//   }
type issue struct {
	Name        string `json:"name"`
	Code        string `json:"code"`
	Severity    string `json:"severity"`
	Path        string `json:"path"`
	Description string `json:"description"`
	Line        int    `json:"line"`
	Char        int    `json:"char"`
}

type Formatter struct{}

var _ core.Formatter = &Formatter{}

func (*Formatter) Print(issues []core.Issue, w io.Writer) error {
	for _, issue := range issues {
		i := convert(issue)
		bytes, err := json.Marshal(&i)
		if err != nil {
			return err
		}
		bytes = append(bytes, []byte("\n")...)
		_, err = w.Write(bytes)
		if err != nil {
			return err
		}
	}

	return nil
}

func convert(iss core.Issue) issue {
	return issue{
		Name:        iss.Tag,
		Code:        iss.Source,
		Path:        iss.File,
		Line:        iss.Line,
		Char:        iss.Column,
		Severity:    mapSeverity(iss.Severity),
		Description: iss.Message,
	}
}

func mapSeverity(s string) string {
	return s
}
