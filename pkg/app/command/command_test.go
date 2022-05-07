package command_test

import (
	"github.com/iarkhanhelsky/check_diff/pkg/app/command"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewCommand(t *testing.T) {
	testCases := map[string]struct {
		params   command.Params
		expected interface{}
	}{
		"version": {
			params: command.Params{
				Type: command.RunVersion,
			},
			expected: &command.Version{},
		},
		"check": {
			params: command.Params{
				Type: command.RunCheck,
			},
			expected: &command.Check{},
		},
		"none": {
			params: command.Params{
				Type: command.RunNone,
			},
			expected: &command.None{},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			assert := assert.New(t)
			cmd, err := command.NewCommand(tc.params)
			assert.NoError(err)
			assert.IsType(tc.expected, cmd)
		})
	}
}
