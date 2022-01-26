package app

import (
	"github.com/iarkhanhelsky/check_diff/pkg/core"
	"go.uber.org/config"
)

const defaultConfigName = "check_diff.yaml"
const defaultVendorDir = ".check_diff/vendor"

type ConfigReader interface {
	Read() (error, Config)
}

type Config struct {
	OutputFormat string `yaml:"OutputFormat"`
	OutputFile   string `yaml:"OutputFile"`
	VendorDir    string `yaml:"VendorDir"`
	Color        bool   `yaml:"Color"`
}

func newConfig() Config {
	return Config{
		OutputFormat: string(core.STDOUT),
		VendorDir:    defaultVendorDir,
	}
}

func newYaml(options Options) (*config.YAML, error) {
	return config.NewYAML(config.File(options.ConfigFile))
}
