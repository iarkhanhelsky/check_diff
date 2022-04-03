package tools

import (
	"github.com/pkg/errors"
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

func Test_SelectSource(t *testing.T) {
	testCases := map[string]struct {
		targets       map[TargetTuple]TargetSource
		target        TargetTuple
		expected      TargetSource
		expectedError error
	}{
		"Any always wins": {
			targets: map[TargetTuple]TargetSource{
				Any:     {MD5: "xxx"},
				Current: {MD5: "zzz"},
			},
			target:   Current,
			expected: TargetSource{MD5: "xxx"},
		},
		"select by target platform": {
			targets: map[TargetTuple]TargetSource{
				LinuxAMD64: {MD5: "xxx"},
				LinuxARM64: {MD5: "zzz"},
			},
			target:   LinuxAMD64,
			expected: TargetSource{MD5: "xxx"},
		},
		"missing target": {
			targets:       map[TargetTuple]TargetSource{},
			target:        LinuxAMD64,
			expected:      TargetSource{},
			expectedError: errors.New("failed to find target source; platform = " + string(LinuxAMD64)),
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			assert := assert.New(t)
			src, err := selectSource(tc.targets, tc.target)
			if tc.expectedError != nil {
				assert.EqualError(tc.expectedError, err.Error())
			}
			assert.Equal(tc.expected, src)
		})
	}
}
