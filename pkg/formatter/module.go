package formatter

import (
	"github.com/iarkhanhelsky/check_diff/pkg/core"
	"go.uber.org/fx"
)

type Options struct {
	Format     string
	OutputFile string
}

var Module = fx.Options(fx.Provide(NewFormatter))

func NewFormatter(opts Options) (core.Formatter, error) {
	return Format(opts.Format).Formatter()
}
