package command

import (
	"github.com/iarkhanhelsky/check_diff/pkg/core"
	"github.com/iarkhanhelsky/check_diff/pkg/tools"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type Type int

const (
	RunCheck Type = iota
	RunVersion
	RunNone
)

type Command interface {
	Run() error
}

type Params struct {
	fx.In

	Type     Type
	Env      Env
	Checkers []core.Checker `group:"checkers"`
	Registry tools.Registry
	Config   core.Config
	Logger   *zap.SugaredLogger
}

func NewCommand(params Params) (Command, error) {
	var cmd Command
	var err error
	switch params.Type {
	case RunCheck:
		cmd = NewCheck(params.Env, params.Config, params.Checkers, params.Logger, params.Registry)
	case RunVersion:
		cmd = NewVersion(params.Env)
	case RunNone:
		cmd = &None{}
	}

	return cmd, err
}
