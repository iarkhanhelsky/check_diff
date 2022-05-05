package core

import (
	"github.com/golang/mock/gomock"
	"github.com/iarkhanhelsky/check_diff/mocks"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFlow_Run(t *testing.T) {
	testCases := map[string]struct {
		flowOpts   []FlowOption
		ranges     []LineRange
		settings   Settings
		shellSetup func(shell *mocks.MockShell)

		err      string
		expected []Issue
	}{
		"empty line range": {
			ranges:   []LineRange{},
			expected: []Issue{},
		},
		"no matching files": {
			ranges: []LineRange{
				{File: "a.cpp", Start: 1, End: 2},
			},
			flowOpts: []FlowOption{
				WithFileExtensions(".txt"),
			},
			expected: []Issue{},
		},
		"cmd error": {
			ranges: []LineRange{
				{File: "a.txt", Start: 1, End: 2},
				{File: "b.txt", Start: 3, End: 5},
			},
			flowOpts: []FlowOption{
				WithFileExtensions(".txt"),
			},
			settings: Settings{
				Command: "test-lint",
			},
			shellSetup: func(shell *mocks.MockShell) {
				shell.EXPECT().Execute("test-lint", "a.txt", "b.txt").
					Return(nil, errors.New("cmd failed"))
			},

			err: "failed to execute test issues: cmd failed",
		},
		"with command": {
			ranges: []LineRange{
				{File: "a.txt", Start: 1, End: 2},
				{File: "b.txt", Start: 3, End: 5},
			},
			flowOpts: []FlowOption{
				WithFileExtensions(".txt"),
				WithCommand("test-lint-v2"),
			},
			shellSetup: func(shell *mocks.MockShell) {
				shell.EXPECT().Execute("test-lint-v2", "a.txt", "b.txt").
					Return([]byte{}, nil)
			},

			expected: []Issue{},
		},
		"convert error": {
			ranges: []LineRange{
				{File: "a.txt", Start: 1, End: 2},
			},
			flowOpts: []FlowOption{
				WithFileExtensions(".txt"),
				WithConverter(func(bytes []byte) ([]Issue, error) {
					return []Issue{}, errors.New("convert failed")
				}),
			},
			shellSetup: func(shell *mocks.MockShell) {
				shell.EXPECT().Execute("test-lint", "a.txt").
					Return([]byte{}, nil)
			},

			err: "failed to convert test issues: convert failed",
		},
		"no matching issues": {
			ranges: []LineRange{
				{File: "a.txt", Start: 1, End: 2},
			},
			flowOpts: []FlowOption{
				WithFileExtensions(".txt"),
				WithConverter(func(bytes []byte) ([]Issue, error) {
					return []Issue{
						{Line: 10, File: "a.txt"},
					}, nil
				}),
			},
			shellSetup: func(shell *mocks.MockShell) {
				shell.EXPECT().Execute("test-lint", "a.txt").
					Return([]byte{}, nil)
			},

			expected: []Issue{},
		},
		"filter issues": {
			ranges: []LineRange{
				{File: "a.txt", Start: 1, End: 5},
			},
			flowOpts: []FlowOption{
				WithFileExtensions(".txt"),
				WithConverter(func(bytes []byte) ([]Issue, error) {
					return []Issue{
						{Line: 1, File: "a.txt"},
						{Line: 10, File: "a.txt"},
					}, nil
				}),
			},
			shellSetup: func(shell *mocks.MockShell) {
				shell.EXPECT().Execute("test-lint", "a.txt").
					Return([]byte{}, nil)
			},

			expected: []Issue{
				{Line: 1, File: "a.txt"},
			},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			assert := assert.New(t)
			ctrl := gomock.NewController(t)

			// Build test defaults, everything can be overriden by
			// tc.flowOpts
			opts := []FlowOption{
				WithConverter(func(bytes []byte) ([]Issue, error) {
					return []Issue{}, nil
				}),
				WithCommand("test-lint"),
			}
			opts = append(opts, tc.flowOpts...)
			opts = append(opts, withShellFactory(func(option ...ShellOption) Shell {
				shell := mocks.NewMockShell(ctrl)
				assert.NotNil(tc.shellSetup)
				tc.shellSetup(shell)
				return shell
			}))

			flow := NewFlow("test", tc.settings, opts...)
			result, err := flow.Run(tc.ranges)
			if tc.err != "" {
				assert.EqualError(err, tc.err)
			} else {
				assert.NoError(err)
				assert.Equal(tc.expected, result)
			}
		})
	}
}

func TestNewFlow(t *testing.T) {

}
