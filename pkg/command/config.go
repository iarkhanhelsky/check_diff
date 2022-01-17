package command

import (
	"io/ioutil"
	"os"

	yaml "gopkg.in/yaml.v2"
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
}

func ParseConfig(path string) (Config, error) {
	file, err := os.Open(path)
	if err != nil {
		return Config{}, err
	}
	defer file.Close()

	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		return Config{}, err
	}

	var config Config
	err = yaml.Unmarshal(bytes, &config)
	return config, err
}
