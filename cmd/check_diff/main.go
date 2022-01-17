package main

import (
	"fmt"
	"github.com/iarkhanhelsky/check_diff/pkg/checker"
	"github.com/iarkhanhelsky/check_diff/pkg/checker/k8s_kubelint"
	"github.com/iarkhanhelsky/check_diff/pkg/command"
	"github.com/iarkhanhelsky/check_diff/pkg/core"
	"os"
	"path"
)

func main() {
	options := command.ParseArgs(os.Args)
	_, err := command.ParseConfig(options.ConfigFile)

	enabledCheckers := []checker.Checker{
		k8s_kubelint.NewK8KubeLint(),
	}
	os.MkdirAll(path.Join(".check_diff", "vendor"), os.ModePerm)
	for _, checker := range enabledCheckers {
		for _, d := range checker.Downloads() {
			err := d.Download(path.Join(".check_diff", "vendor"))
			if err != nil {
				panic(err)
			}
		}
	}

	for _, checker := range enabledCheckers {
		checker.Check([]core.LineRange{})
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v", err)
		os.Exit(2)
	}
}
