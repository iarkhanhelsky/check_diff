package kubelint

import (
	"go.uber.org/config"
	"go.uber.org/fx"
)

var Module = fx.Options(
	fx.Provide(func(yaml *config.YAML) (Settings, error) {
		v := Settings{}
		err := yaml.Get("KubeLint").Populate(&v)
		return v, err
	}),
	fx.Provide(NewK8KubeLint),
)
