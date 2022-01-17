package stdout

import (
	"fmt"
	"github.com/iarkhanhelsky/check_diff/pkg/core"
	"github.com/iarkhanhelsky/check_diff/pkg/formatter"
	"io"
	"math"
	"strconv"
	"strings"
)

type Formatter struct {
	printer     printer
	currentFile string
	fileReader  core.FileReader
}

func (*Formatter) Supports() []formatter.Format {
	return []formatter.Format{formatter.STDOUT}
}

func (formatter *Formatter) Print(issues []core.Issue, w io.Writer) error {
	if len(issues) == 0 {
		return nil
	}

	for _, issue := range issues {
		err := formatter.printIssue(issue)
		if err != nil {
			return err
		}
	}

	uniqFiles := countUniqFiles(issues)
	return formatter.println(fmt.Sprintf("Total: %d in %d files", len(issues), uniqFiles))
}

func (formatter *Formatter) println(str string) error {
	return formatter.printer.println(str)
}

func (formatter *Formatter) printIssue(issue core.Issue) error {
	w := formatter.printer.w()
	err := w(fmt.Sprintf("%s:%d", issue.File, issue.Line)) // white, bright
	if err != nil {
		return err
	}
	err = w(fmt.Sprintf("[%s] %s", issue.Severity, issue.Message)) // magenta, bright
	if err != nil {
		return err
	}
	err = formatter.fileBanner(issue)
	if err != nil {
		return err
	}
	err = formatter.println("")
	if err != nil {
		return err
	}

	return nil
}

func (formatter *Formatter) fileBanner(issue core.Issue) error {
	w := formatter.printer.w()
	contextLines, offset, err := formatter.readContext(issue.File)
	if err != nil {
		// log to stderr
		return nil
	}

	margin := int(math.Ceil(math.Log10(float64(offset + len(contextLines)))))
	for i, line := range contextLines {
		l := offset + i
		err := w(fmt.Sprintf("%"+strconv.Itoa(margin)+"d:%s", l, line)) // white, dim
		if err != nil {
			return err
		}
		if l == issue.Line {
			// TODO: print line number with bg_magenta, bright
			err := w(rjust("^", ' ', margin+1+issue.Column)) // white, bright
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (formatter *Formatter) readContext(file string) ([]string, int, error) {
	return nil, 0, nil
}

func rjust(s string, ch byte, size int) string {
	return strings.Repeat(string(ch), len(s)-size)
}

func countUniqFiles(issues []core.Issue) int {
	var fileset map[string]int
	for _, issue := range issues {
		fileset[issue.File] = 1
	}
	return len(fileset)
}
