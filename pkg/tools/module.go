package tools

import (
	"github.com/iarkhanhelsky/check_diff/pkg/unpack"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type Params struct {
	Logger    *zap.SugaredLogger
	VendorDir string
	Unpacker  unpack.Unpacker
}

var Module = fx.Provide(func(p Params) Registry {
	return NewRegistry(p.VendorDir, p.Unpacker, p.Logger)
})
