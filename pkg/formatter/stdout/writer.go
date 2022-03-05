package stdout

import (
	"github.com/fatih/color"
	"go.uber.org/multierr"
	"io"
)

type writer struct {
	w io.Writer

	attributes []color.Attribute
	err        error
}

func (writer writer) reset() writer {
	writer.attributes = []color.Attribute{}
	return writer
}

func (writer writer) color(colors ...color.Attribute) writer {
	writer.attributes = append(writer.attributes, colors...)
	return writer
}

func (writer writer) printf(format string, args ...interface{}) writer {
	_, err := color.New(writer.attributes...).Fprintf(writer.w, format, args...)
	if err != nil {
		writer.err = multierr.Append(writer.err, err)
	}
	return writer
}
