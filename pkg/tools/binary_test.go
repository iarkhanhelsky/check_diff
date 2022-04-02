package tools

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestBinary_Digest(t *testing.T) {
	x := Binary{
		DstFile: "a",
		Path:    "a",
		Targets: map[TargetTuple]TargetSource{
			Any: {
				Urls: []string{"https//a.com/v1.0.0/a.tar.gz"},
				MD5:  "68b329da9893e34099c7d8ad5cb9c940",
			},
		},
	}

	y := Binary{
		DstFile: "a",
		Path:    "a",
		Targets: map[TargetTuple]TargetSource{
			Any: {
				Urls: []string{"https//a.com/v1.1.0/a.tar.gz"},
				MD5:  "b026324c6904b2a9cb4b88d6d61c81d1",
			},
		},
	}

	xdigest, err := x.digest()
	assert.NoError(t, err)
	ydigest, err := y.digest()
	assert.NotEqual(t, xdigest, ydigest)
}