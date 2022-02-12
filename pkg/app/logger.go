package app

import (
	"go.uber.org/zap"
	"os"
)

func NewLogger() (*zap.Logger, error) {
	trace := false
	for _, arg := range os.Args {
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
