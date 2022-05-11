package command_test

import (
	"bytes"
	"github.com/golang/mock/gomock"
	mockcore "github.com/iarkhanhelsky/check_diff/mocks/pkg/core"
	mocktools "github.com/iarkhanhelsky/check_diff/mocks/pkg/tools"
	"github.com/iarkhanhelsky/check_diff/pkg/app/command"
	"github.com/iarkhanhelsky/check_diff/pkg/core"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"path"
	"testing"
)

func TestCheck_Run(t *testing.T) {
	stubCheckers := func(ctrl *gomock.Controller) []core.Checker {
		checker := mockcore.NewMockChecker(ctrl)
		checker.EXPECT().
			Check([]core.LineRange{
				{File: "testdata/java/src/main/java/Main.java", Start: 1, End: 5},
			}).
			Return([]core.Issue{
				{Line: 1, Message: "test", File: "testdata/java/src/main/java/Main.java"},
			}, nil)
		checker.EXPECT().Tag().Return("test")

		return []core.Checker{checker}
	}

	testCases := map[string]struct {
		env      command.Env
		config   core.Config
		checkers func(ctrl *gomock.Controller) []core.Checker

		err string
	}{
		"empty input": {
			env: command.Env{
				OutWriter: &bytes.Buffer{},
			},
		},
		"file input": {
			config: core.Config{
				InputFile:  "testdata/input.diff",
				OutputFile: "test.out",
			},
		},
		"unknown format": {
			config: core.Config{
				InputFile:    "testdata/input.diff",
				OutputFile:   "test.out",
				OutputFormat: "unknown",
			},
			checkers: stubCheckers,
			err:      "can't create formatter: unknown formatter type: 'unknown'",
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			assert := assert.New(t)

			ctrl := gomock.NewController(t)
			registry := mocktools.NewMockRegistry(ctrl)
			registry.EXPECT().Install(gomock.Any()).Return(nil)

			tmp := t.TempDir()
			tc.config.VendorDir = path.Join(tmp, ".check_diff")
			if tc.config.OutputFormat == "" {
				tc.config.OutputFormat = "stdout"
			}
			if tc.config.OutputFile != "" {
				tc.config.OutputFile = path.Join(tmp, tc.config.OutputFile)
			}

			checkers := []core.Checker{}
			if tc.checkers != nil {
				checkers = tc.checkers(ctrl)
			}

			cmd := command.NewCheck(
				tc.env,
				tc.config,
				checkers,
				zap.NewNop().Sugar(),
				registry)
			if err := cmd.Run(); tc.err == "" {
				assert.NoError(err)
			} else {
				assert.EqualError(err, tc.err)
			}
		})
	}
}
