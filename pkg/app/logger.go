package app

import (
	"github.com/iarkhanhelsky/check_diff/pkg/app/command"
	"go.uber.org/zap"
)

func NewLogger(env command.Env) (*zap.Logger, error) {
	trace := false
	for _, arg := range env.Args {
		if arg == "--trace" {
			trace = true
			break
		}
	}
	if trace {
		cfg := zap.NewDevelopmentConfig()
		cfg.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
		cfg.OutputPaths = []string{"stderr"}
		return cfg.Build()
	}

	return zap.NewNop(), nil
}
