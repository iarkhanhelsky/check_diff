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
	name        string `json:"name"`
	code        string `json:"code"`
	severity    string `json:"severity"`
	path        string `json:"path"`
	description string `json:"description"`
	line        int    `json:"line"`
	char        int    `json:"char"`
}

type Formatter struct{}

var _ core.Formatter = &Formatter{}

func (*Formatter) Supports() []core.Format {
	return []core.Format{core.Phabricator}
}

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
		name:        iss.Tag,
		code:        iss.Source,
		path:        iss.File,
		line:        iss.Line,
		char:        iss.Column,
		severity:    mapSeverity(iss.Severity),
		description: iss.Message,
	}
}

func mapSeverity(s string) string {
	return s
}
