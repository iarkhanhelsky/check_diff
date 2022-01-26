package core

type Options struct {
	Exclude []string `yaml:"Exclude"`
	Include []string `yaml:"Include"`
	Command string   `yaml:"Command"`
	Config  string   `yaml:"Config"`
}
