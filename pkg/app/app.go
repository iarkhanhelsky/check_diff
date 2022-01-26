package app

import (
	"fmt"
	"github.com/iarkhanhelsky/check_diff/pkg/core"
	"github.com/iarkhanhelsky/check_diff/pkg/formatter/stdout"
	"os"
)

func Run() {
	cliOptions := ParseArgs(os.Args)
	config, err := ParseConfig(cliOptions.ConfigFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(2)
	}

	fmt.Printf("%v", config)
	os.MkdirAll(config.VendorDir, os.ModePerm)

	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v", err)
		os.Exit(2)
	}

	var issues []core.Issue
	for _, checker := range []core.Checker{} {
		checked, err := checker.Check([]core.LineRange{})
		if err != nil {
			panic(err)
		}
		issues = append(issues, checked...)
	}

	err = stdout.NewFormatter().Print(issues, os.Stdout)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v", err)
		os.Exit(2)
	}
}
