package app

import (
	"go.uber.org/fx"
	"os"
)

var Module = fx.Options(
	fx.Provide(
		func() Options { return ParseArgs(os.Args) },
		newConfig,
		newYaml,
		newFormatter,
	),
)
