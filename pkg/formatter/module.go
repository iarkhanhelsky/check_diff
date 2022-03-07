package formatter

import (
	"github.com/iarkhanhelsky/check_diff/pkg/core"
	"go.uber.org/fx"
)

type Params struct {
	Format string
}

var Module = fx.Options(fx.Provide(NewFormatter))

func NewFormatter(p Params) (core.Formatter, error) {
	return Format(p.Format).Formatter()
}
