package core

type ConfigReader interface {
	Read() (error, Config)
}

type Config struct {
	OutputFormat string `yaml:"OutputFormat"`
	OutputFile   string `yaml:"OutputFile"`
	VendorDir    string `yaml:"VendorDir"`
	Color        bool   `yaml:"Color"`
}

func NewDefaultConfig() Config {
	return Config{
		OutputFormat: "stdout",
	}
}
