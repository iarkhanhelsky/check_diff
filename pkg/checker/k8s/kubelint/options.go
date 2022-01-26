package kubelint

import (
	"github.com/iarkhanhelsky/check_diff/pkg/core"
)

type Options struct {
	core.Options `yaml:"KubeLinter"`
}
