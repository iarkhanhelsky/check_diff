package kubelinter

import (
	"github.com/iarkhanhelsky/check_diff/pkg/core"
	"go.uber.org/config"
)

type Settings struct {
	core.Settings `yaml:",inline"`
}

func ReadSettings(yaml *config.YAML) (Settings, error) {
	v := Settings{}
	err := yaml.Get("KubeLinter").Populate(&v)
	return v, err
}
