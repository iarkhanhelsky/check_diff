package core

import (
	"strconv"
	"strings"
)

type DiffParser interface {
	ParseNextLine(line string)
	Result() []LineRange
}

var _ DiffParser = &diffParser{}

type diffParser struct {
	ranges []LineRange

	file          string
	currentLine   int
	changeStarted bool
}

func NewDiffParser() DiffParser {
	return &diffParser{}
}

func (p *diffParser) ParseNextLine(line string) {
	if strings.HasPrefix(line, "+++ b/") {
		end := len(line)
		tabIndex := strings.Index(line, "\t")
		if tabIndex > 0 {
			end = tabIndex
		}
		p.file = line[len("+++ b/"):end]
	} else if strings.HasPrefix(line, "+++ /dev/nil") {
		p.file = ""
	} else if strings.HasPrefix(line, "@@") && strings.HasSuffix(line, "@@") {
		// Don't have a file, skip this line
		if len(p.file) == 0 {
			return
		}

		// Case like this:
		// 		@@ -1,10 +1,10 @@
		// Our main interest is a last part, i.e. +1,10. Which means 10 lines
		// starting from line number 1.

		// > If a hunk contains just one line, only its start line number
		//   appears. Otherwise its line numbers look like ‘start,count’. An
		//  empty hunk is considered to start at the line that follows the hunk.
		//  If a hunk and its context contain two or more lines, its line
		//  numbers look like ‘start,count’. Otherwise only its end line number
		//  appears. An empty hunk is considered to end at the line that
		//  precedes the hunk.
		// Link: https://www.gnu.org/software/diffutils/manual/html_node/Detailed-Unified.html#Detailed-Unified
		lineRange := line[3 : len(line)-3] // strip '@@ ' from the beginning and ' @@' from the end
		// "-1,10 +1,10".split(" ") => ["-1,10", "+1,10"] => "+1,10"
		startAndCount := strings.Split(strings.Split(lineRange, " ")[1], ",")

		start, _ := strconv.Atoi(startAndCount[0])
		p.currentLine = start - 1 // We'll add 1 when start processing change lines
	} else if strings.HasPrefix(line, "+") {
		p.currentLine++

		if !p.changeStarted {
			p.changeStarted = true
			lineRange := LineRange{File: p.file, Start: p.currentLine, End: p.currentLine}
			p.ranges = append(p.ranges, lineRange)
		} else {
			p.ranges[len(p.ranges)-1].End = p.currentLine
		}
	} else if strings.HasPrefix(line, "-") {
		// Do nothing, not interested in removals and they do not affect line
		// number
	} else {
		// End of Continuous sequence of changes inside hunk. Waiting for new.
		p.changeStarted = false
		// Even if there is no changes we still need to update line number
		p.currentLine++
	}
}

func (p *diffParser) Result() []LineRange {
	return p.ranges
}
