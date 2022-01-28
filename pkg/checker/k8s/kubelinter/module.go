package kubelinter

import (
	"go.uber.org/config"
	"go.uber.org/fx"
)

var Module = fx.Options(
	fx.Provide(func(yaml *config.YAML) (Settings, error) {
		v := Settings{}
		err := yaml.Get("KubeLinter").Populate(&v)
		return v, err
	}),
	fx.Provide(NewK8KubeLint),
)
