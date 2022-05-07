package core_test

import (
	"github.com/golang/mock/gomock"
	mockcore "github.com/iarkhanhelsky/check_diff/mocks/pkg/core"
	"github.com/iarkhanhelsky/check_diff/pkg/core"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFlow_Run(t *testing.T) {
	testCases := map[string]struct {
		flowOpts   []core.FlowOption
		ranges     []core.LineRange
		settings   core.Settings
		shellSetup func(shell *mockcore.MockShell)

		err      string
		expected []core.Issue
	}{
		"empty line range": {
			ranges:   []core.LineRange{},
			expected: []core.Issue{},
		},
		"no matching files": {
			ranges: []core.LineRange{
				{File: "a.cpp", Start: 1, End: 2},
			},
			flowOpts: []core.FlowOption{
				core.WithFileExtensions(".txt"),
			},
			expected: []core.Issue{},
		},
		"cmd error": {
			ranges: []core.LineRange{
				{File: "a.txt", Start: 1, End: 2},
				{File: "b.txt", Start: 3, End: 5},
			},
			flowOpts: []core.FlowOption{
				core.WithFileExtensions(".txt"),
			},
			settings: core.Settings{
				Command: "test-lint",
			},
			shellSetup: func(shell *mockcore.MockShell) {
				shell.EXPECT().Execute("test-lint", "a.txt", "b.txt").
					Return(nil, errors.New("cmd failed"))
			},

			err: "failed to execute test issues: cmd failed",
		},
		"with command": {
			ranges: []core.LineRange{
				{File: "a.txt", Start: 1, End: 2},
				{File: "b.txt", Start: 3, End: 5},
			},
			flowOpts: []core.FlowOption{
				core.WithFileExtensions(".txt"),
				core.WithCommand("test-lint-v2"),
			},
			shellSetup: func(shell *mockcore.MockShell) {
				shell.EXPECT().Execute("test-lint-v2", "a.txt", "b.txt").
					Return([]byte{}, nil)
			},

			expected: []core.Issue{},
		},
		"convert error": {
			ranges: []core.LineRange{
				{File: "a.txt", Start: 1, End: 2},
			},
			flowOpts: []core.FlowOption{
				core.WithFileExtensions(".txt"),
				core.WithConverter(func(bytes []byte) ([]core.Issue, error) {
					return []core.Issue{}, errors.New("convert failed")
				}),
			},
			shellSetup: func(shell *mockcore.MockShell) {
				shell.EXPECT().Execute("test-lint", "a.txt").
					Return([]byte{}, nil)
			},

			err: "failed to convert test issues: convert failed",
		},
		"no matching issues": {
			ranges: []core.LineRange{
				{File: "a.txt", Start: 1, End: 2},
			},
			flowOpts: []core.FlowOption{
				core.WithFileExtensions(".txt"),
				core.WithConverter(func(bytes []byte) ([]core.Issue, error) {
					return []core.Issue{
						{Line: 10, File: "a.txt"},
					}, nil
				}),
			},
			shellSetup: func(shell *mockcore.MockShell) {
				shell.EXPECT().Execute("test-lint", "a.txt").
					Return([]byte{}, nil)
			},

			expected: []core.Issue{},
		},
		"filter issues": {
			ranges: []core.LineRange{
				{File: "a.txt", Start: 1, End: 5},
			},
			flowOpts: []core.FlowOption{
				core.WithFileExtensions(".txt"),
				core.WithConverter(func(bytes []byte) ([]core.Issue, error) {
					return []core.Issue{
						{Line: 1, File: "a.txt"},
						{Line: 10, File: "a.txt"},
					}, nil
				}),
			},
			shellSetup: func(shell *mockcore.MockShell) {
				shell.EXPECT().Execute("test-lint", "a.txt").
					Return([]byte{}, nil)
			},

			expected: []core.Issue{
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
			opts := []core.FlowOption{
				core.WithConverter(func(bytes []byte) ([]core.Issue, error) {
					return []core.Issue{}, nil
				}),
				core.WithCommand("test-lint"),
			}
			opts = append(opts, tc.flowOpts...)
			opts = append(opts, core.WithShellFactory(func(option ...core.ShellOption) core.Shell {
				shell := mockcore.NewMockShell(ctrl)
				assert.NotNil(tc.shellSetup)
				tc.shellSetup(shell)
				return shell
			}))

			flow := core.NewFlow("test", tc.settings, opts...)
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
