package app

import (
	"fmt"
	flag "github.com/spf13/pflag"
	"os"
)

const (
	defaultConfigName = "check_diff.yaml"
	defaultVendorDir  = ".check_diff/vendor"
)

type CliOptions struct {
	InputFile   string
	Format      string
	OutputFile  string
	ConfigFile  string
	VendorDir   string
	FailOnError bool
	// Trace is not really used, but we generate flag for help entry
	// --trace is checked in NewLogger function, as CliOptions can't be provided
	// before Logger.
	Trace   bool
	Version bool
}

func parseArgs(args []string) CliOptions {
	var outputfile, format, configfile, vendordir, inputFile string

	flagset := flag.NewFlagSet(args[0], flag.ContinueOnError)
	flagset.StringVarP(&format, "format", "f", "", "Output format. One of: stdout,phabricator,codeclimate,gitlab")
	flagset.StringVarP(&outputfile, "output-file", "o", "", "Output file path")
	flagset.StringVarP(&configfile, "config", "c", defaultConfigName, "Config file path")
	flagset.StringVarP(&vendordir, "vendor-dir", "", defaultVendorDir, "vendor directory to store intermediate data")
	flagset.StringVarP(&inputFile, "input", "i", "", "Input file. Read from STDIN if not set")
	noFailOnError := flagset.BoolP("no-fail", "", false, "")
	trace := flagset.BoolP("trace", "", false, "Enable debug logs")
	version := flagset.BoolP("version", "", false, "Print version and exit")

	err := flagset.Parse(args[1:])
	if err == flag.ErrHelp {
		os.Exit(0)
	} else if err != nil {
		fmt.Printf("error: %v\n", err)
		os.Exit(2)
	} else if *version {
		fmt.Printf("check_diff v%s\n", Version)
		os.Exit(0)
	}

	return CliOptions{
		InputFile:   inputFile,
		Format:      format,
		OutputFile:  outputfile,
		ConfigFile:  configfile,
		VendorDir:   vendordir,
		FailOnError: !(*noFailOnError),
		Trace:       *trace,
	}
}

func NewCliOptions() CliOptions {
	return parseArgs(os.Args)
}
