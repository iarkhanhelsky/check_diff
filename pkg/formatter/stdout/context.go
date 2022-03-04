package stdout

import (
	"bufio"
	"math"
	"os"
)

type contextReader interface {
	readContext(path string, line int) ([]string, int, error)
}

var _ contextReader = &fileReader{}

type fileReader struct {
	defaultOffset int

	path    string
	line    int
	scanner *bufio.Scanner
	file    *os.File
}

func newCachedFileContext(defaultOffset int) contextReader {
	return &fileReader{defaultOffset: defaultOffset}
}

func (reader *fileReader) readContext(path string, line int) ([]string, int, error) {
	scanner, err := reader.prepareScanner(path, line)
	if err != nil {
		return nil, 0, err
	}

	currentLine := reader.line
	startLine := int(math.Max(1, float64(line-reader.defaultOffset)))
	endLine := line + reader.defaultOffset
	lines := make([]string, endLine-startLine+1)

	sz := 0
	for scanner.Scan() {
		if startLine <= currentLine && currentLine <= endLine {
			lines[currentLine-startLine] = scanner.Text()
			sz++
		}
		if currentLine == endLine {
			break
		}
		currentLine += 1
	}

	reader.line = currentLine
	return lines[:sz], line - startLine, scanner.Err()
}

func (reader *fileReader) prepareScanner(path string, line int) (*bufio.Scanner, error) {
	if reader.path != path || reader.line > line {
		if reader.scanner != nil {
			if err := reader.file.Close(); err != nil {
				return nil, err
			}
		}
		file, err := os.Open(path)
		reader.file = file
		if err != nil {
			return nil, err
		}

		reader.scanner = bufio.NewScanner(reader.file)
		reader.line = 1
	}

	return reader.scanner, nil
}
