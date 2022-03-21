package command

import "io"

type Env struct {
	Args      []string
	OutWriter io.Writer
	ErrWriter io.Writer

	Version string
	Commit  string
	Date    string
}
