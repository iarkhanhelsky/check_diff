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
	writer        writer
	currentFile   string
	contextReader contextReader
}

func NewFormatter() formatter.Formatter {
	return &Formatter{
		contextReader: newCachedFileContext(5),
	}
}

func (*Formatter) Supports() []formatter.Format {
	return []formatter.Format{formatter.STDOUT}
}

func (formatter *Formatter) Print(issues []core.Issue, w io.Writer) error {
	formatter.writer = writer{w: w}

	if len(issues) == 0 {
		return nil
	}

	for _, issue := range issues {
		formatter.writer.color(color.FgWhite, color.Bold).printf("%s:%d\n", issue.File, issue.Line).reset()
		formatter.writer.color(color.FgHiMagenta, color.Bold).printf("[%s] %s\n", issue.Severity, issue.Message).reset()
		formatter.fileBanner(issue)
		formatter.writer.printf("\n")
	}

	uniqFiles := countUniqFiles(issues)
	formatter.writer.color().printf("Total: %d in %d files\n", len(issues), uniqFiles).reset()
	return formatter.writer.err()
}

func (formatter *Formatter) fileBanner(issue core.Issue) {
	contextLines, offset, err := formatter.readContext(issue.File, issue.Line)
	if err != nil {
		fmt.Fprintf(os.Stderr, "err reading file: %s: %v\n", issue.File, err)
		return
	}

	margin := int(math.Ceil(math.Log10(float64(offset + len(contextLines)))))
	for i, line := range contextLines {
		l := offset + i + 1
		w := formatter.writer.color(color.FgWhite)
		if l == issue.Line {
			w = formatter.writer.color(color.Bold, color.BgMagenta)
		}
		w.printf("%"+strconv.Itoa(margin)+"d:", l)

		w = formatter.writer
		if l == issue.Line {
			w = formatter.writer.color(color.FgHiWhite)
		}
		w = w.printf(" %s\n", line).reset()

		if l == issue.Line {
			w.color(color.FgWhite, color.Bold).printf(rjust("", ' ', margin+2+issue.Column) + "^\n")
		}
	}
}

func (formatter *Formatter) readContext(file string, line int) ([]string, int, error) {
	return formatter.contextReader.readContext(file, line)
}

func rjust(s string, ch byte, size int) string {
	if len(s) < size {
		return strings.Repeat(string(ch), size-len(s))
	}

	return s
}

func countUniqFiles(issues []core.Issue) int {
	fileset := make(map[string]int, 0)
	for _, issue := range issues {
		fileset[issue.File] = 1
	}
	return len(fileset)
}
