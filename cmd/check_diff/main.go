package main

import (
	"github.com/iarkhanhelsky/check_diff/pkg/app"
	"github.com/iarkhanhelsky/check_diff/pkg/app/command"
	"os"
)

func main() {
	err := app.Main(command.Env{
		Args: os.Args, OutWriter: os.Stdout, ErrWriter: os.Stderr,
	})

	if err != nil {
		// TODO: Dispatch different types of errors
		panic(err)
	}
}
