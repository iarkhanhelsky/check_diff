package app

import (
	"github.com/stretchr/testify/assert"
	"io"
	"testing"
)

type _writer struct {
}

func (n2 _writer) Write(p []byte) (n int, err error) {
	return len(p), nil
}

var devNull io.Writer = &_writer{}

func args(args ...string) []string {
	return append([]string{"check_diff"}, args...)
}

func TestParseArgs(t *testing.T) {
	testcases := map[string]struct {
		args []string
		opts CliOptions
		err  bool
	}{
		"no args": {
			args: args(),
			opts: CliOptions{
				VendorDir:   defaultVendorDir,
				ConfigFile:  defaultConfigName,
				FailOnError: true,
			},
		},
		"unknown flag": {
			args: args("--unknown-flag"),
			err:  true,
		},
		"show help": {
			args: args("--help"),
			err:  true,
		},

		"read input": {
			args: args("-i", "test.diff"),
			opts: CliOptions{
				VendorDir:   defaultVendorDir,
				ConfigFile:  defaultConfigName,
				FailOnError: true,

				InputFile: "test.diff",
			},
		},
		"disable color": {
			args: args("--no-color"),
			opts: CliOptions{
				VendorDir:   defaultVendorDir,
				ConfigFile:  defaultConfigName,
				FailOnError: true,

				NoColor: func() *bool { x := true; return &x }(),
			},
		},
	}

	for name, tc := range testcases {
		t.Run(name, func(t *testing.T) {
			assert := assert.New(t)

			tc.opts.output = devNull
			opts := newCliOptions(devNull)
			err := opts.parseArgs(tc.args)

			if err != nil {
				assert.Equal(tc.err, err != nil)
			} else {
				assert.Equal(tc.opts, *opts)
			}
		})
	}
}
