package app

import "go.uber.org/config"

func NewYaml(options CliOptions) (*config.YAML, error) {
	return config.NewYAML(config.File(options.ConfigFile))
}
