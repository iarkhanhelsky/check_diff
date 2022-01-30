package kubelinter

import (
	"go.uber.org/fx"
)

var Module = fx.Options(
	fx.Provide(ReadSettings),
	fx.Provide(NewK8KubeLint),
)
