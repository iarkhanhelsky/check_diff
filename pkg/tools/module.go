package tools

import (
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type Params struct {
	Logger    *zap.SugaredLogger
	VendorDir string
}

var Module = fx.Provide(func(p Params) Registry {
	return NewRegistry(p.Logger, p.VendorDir)
})
