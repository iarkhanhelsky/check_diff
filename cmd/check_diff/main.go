package main

import (
	"github.com/iarkhanhelsky/check_diff/pkg/app"
	"github.com/iarkhanhelsky/check_diff/pkg/app/command"
	"os"
)

// These values are populated by go-releaser
var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func main() {
	err := app.Main(command.Env{
		Args: os.Args, OutWriter: os.Stdout, ErrWriter: os.Stderr,
		Version: version, Commit: commit, Date: date,
	})

	if err != nil {
		// TODO: Dispatch different types of errors
		panic(err)
	}
}
