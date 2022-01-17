package command

import (
	"fmt"
	flag "github.com/spf13/pflag"
	"os"
)

type Options struct {
	Format      string
	OutputFile  string
	ConfigFile  string
	VendorDir   string
	FailOnError bool
}

func ParseArgs(args []string) Options {
	var outputfile, format, configfile, vendordir string

	flagset := flag.NewFlagSet(args[0], flag.ContinueOnError)
	flagset.StringVarP(&format, "format", "f", "stdout", "Output format. One of: stdout,phabricator,codeclimate,gitlab")
	flagset.StringVarP(&outputfile, "output-file", "o", "", "Output file path")
	flagset.StringVarP(&configfile, "config", "c", defaultConfigName, "Config file path")
	flagset.StringVarP(&vendordir, "vendor-dir", "", defaultVendorDir, "vendor directory to store intermediate data")
	noFailOnError := flag.BoolP("no-fail", "", false, "")

	err := flagset.Parse(args[1:])
	if err == flag.ErrHelp {

	} else if err != nil {
		fmt.Printf("error: %v\n", err)
		os.Exit(2)
	}

	return Options{
		Format:      format,
		OutputFile:  outputfile,
		ConfigFile:  configfile,
		VendorDir:   vendordir,
		FailOnError: !(*noFailOnError),
	}
}
