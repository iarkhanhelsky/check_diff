package stdout

import (
	"github.com/fatih/color"
	"io"
)

type writer struct {
	w io.Writer

	attributes []color.Attribute
	errors     []error
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
		writer.errors = append(writer.errors, err)
	}
	return writer
}

func (writer writer) err() error {
	if len(writer.errors) > 0 {
		return writer.errors[0]
	}

	return nil
}
