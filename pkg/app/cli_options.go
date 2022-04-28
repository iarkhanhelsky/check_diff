package app

import (
	"github.com/iarkhanhelsky/check_diff/pkg/app/command"
	"github.com/iarkhanhelsky/check_diff/pkg/formatter"
	flag "github.com/spf13/pflag"
	"io"
	"strings"
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
	NoColor     *bool
	// Trace is not really used, but we generate flag for help entry
	// --trace is checked in NewLogger function, as CliOptions can't be provided
	// before Logger.
	Trace bool

	// Print version and exit
	version bool
	output  io.Writer
}

func (opts *CliOptions) parseArgs(args []string) error {
	flagset := flag.NewFlagSet(args[0], flag.ContinueOnError)
	flagset.SetOutput(opts.output)

	flagset.StringVarP(&opts.Format, "format", "f", "", "Output format. One of: "+strings.Join(formatter.FormatNames(), ", "))
	flagset.StringVarP(&opts.OutputFile, "output-file", "o", "", "Output file path")
	flagset.StringVarP(&opts.ConfigFile, "config", "c", defaultConfigName, "Config file path")
	flagset.StringVarP(&opts.VendorDir, "vendor-dir", "", defaultVendorDir, "vendor directory to store intermediate data")
	flagset.StringVarP(&opts.InputFile, "input", "i", "", "Input file. Read from STDIN if not set")

	noFailOnError := flagset.BoolP("no-fail", "", false, "")
	trace := flagset.BoolP("trace", "", false, "Enable debug logs")
	noColor := flagset.BoolP("no-color", "", false, "Disable colors")
	version := flagset.BoolP("version", "", false, "Print version and exit")

	err := flagset.Parse(args[1:])

	opts.FailOnError = !(*noFailOnError)
	opts.Trace = *trace

	// Assign this value only if flag was set
	for _, a := range args[1:] {
		if a == "--no-color" {
			opts.NoColor = noColor
			break
		}
	}

	opts.version = *version

	return err
}

func newCliOptions(output io.Writer) *CliOptions {
	return &CliOptions{output: output}
}

func NewCliOptions(env command.Env) (CliOptions, command.Type, error) {
	var commandType = command.RunNone
	opts := newCliOptions(env.OutWriter)
	err := opts.parseArgs(env.Args)
	if err == flag.ErrHelp {
		err = nil // Showed help
	} else if err != nil {

	} else if opts.version {
		commandType = command.RunVersion
	} else {
		commandType = command.RunCheck
	}

	return *opts, commandType, err
}
