package core

import (
	"io"
)

type Formatter interface {
	Print(issues []Issue, w io.Writer) error
}
