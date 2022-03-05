package command

import (
	"fmt"
	"io"
)

type Version struct {
	version string

	outWriter io.Writer
}

var _ Command = &Version{}

func NewVersion(env Env) Command {
	return &Version{
		version:   "0.0.1",
		outWriter: env.OutWriter,
	}
}

func (version *Version) Run() error {
	fmt.Fprintf(version.outWriter, "check_diff v%s\n", version.version)
	return nil
}
