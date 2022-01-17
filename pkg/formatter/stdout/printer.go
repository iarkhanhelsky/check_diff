package stdout

import "io"

type printer struct {
	writer  io.Writer
	newLine []byte
}

func newPrinter(w io.Writer) printer {
	return printer{
		writer: w, newLine: []byte("\n"),
	}
}

func (p printer) println(str string) error {
	if _, err := p.writer.Write([]byte(str)); err != nil {
		return err
	}
	if _, err := p.writer.Write(p.newLine); err != nil {
		return err
	}

	return nil
}

func (p printer) w() func(string) error {
	return func(s string) error {
		return p.println(s)
	}
}
