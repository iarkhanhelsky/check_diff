package command

import (
	"fmt"
	"io"
)

type Version struct {
	version string
	commit  string
	date    string

	outWriter io.Writer
}

var _ Command = &Version{}

func NewVersion(env Env) Command {
	return &Version{
		version:   env.Version,
		commit:    env.Commit,
		date:      env.Date,
		outWriter: env.OutWriter,
	}
}

func (version *Version) Run() error {
	fmt.Fprintf(version.outWriter, "check_diff v%s (%s, %s)\n",
		version.version, version.commit, version.date)
	return nil
}
