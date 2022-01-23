package stdout

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/iarkhanhelsky/check_diff/pkg/core"
	"github.com/iarkhanhelsky/check_diff/pkg/formatter"
	"io"
	"math"
	"os"
	"strconv"
	"strings"
)

type Formatter struct {
	writer      io.Writer
	currentFile string
	fileReader  core.FileReader
}

func (*Formatter) Supports() []formatter.Format {
	return []formatter.Format{formatter.STDOUT}
}

func (formatter *Formatter) Print(issues []core.Issue, w io.Writer) error {
	formatter.writer = w

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
	_, err := fmt.Fprintf(formatter.writer, "Total: %d in %d files\n", len(issues), uniqFiles)
	if err != nil {
		return err
	}

	return nil
}

func (formatter *Formatter) printIssue(issue core.Issue) error {

	_, err := color.New(color.FgHiWhite, color.Bold).Fprintf(formatter.writer, "%s:%d\n", issue.File, issue.Line)
	if err != nil {
		return err
	}

	_, err = color.New(color.FgHiMagenta, color.Bold).Fprintf(formatter.writer, "[%s] %s\n", issue.Severity, issue.Message)
	if err != nil {
		return err
	}

	err = formatter.fileBanner(issue)
	if err != nil {
		return err
	}

	_, err = fmt.Fprintln(formatter.writer, "")
	if err != nil {
		return err
	}

	return nil
}

func (formatter *Formatter) fileBanner(issue core.Issue) error {
	contextLines, offset, err := formatter.readContext(issue.File)
	if err != nil {
		fmt.Fprintf(os.Stderr, "err reading file: %s: %v\n", issue.File, err)
		// log to stderr
		return nil
	}

	margin := int(math.Ceil(math.Log10(float64(offset + len(contextLines)))))
	for i, line := range contextLines {
		l := offset + i
		_, err := color.New(color.FgWhite).Fprintf(formatter.writer, "%"+strconv.Itoa(margin)+"d:%s\n", l, line)
		if err != nil {
			return err
		}
		if l == issue.Line {
			// TODO: print line number with bg_magenta, bright
			_, err = color.New(color.FgWhite).Fprintf(formatter.writer, rjust("^", ' ', margin+1+issue.Column)) // white, bright
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
	fileset := make(map[string]int, 0)
	for _, issue := range issues {
		fileset[issue.File] = 1
	}
	return len(fileset)
}
