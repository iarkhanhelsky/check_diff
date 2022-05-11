package unpack_test

import (
	"fmt"
	"github.com/golang/mock/gomock"
	unpackmocks "github.com/iarkhanhelsky/check_diff/mocks/pkg/unpack"
	"github.com/iarkhanhelsky/check_diff/pkg/unpack"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"os/exec"
	"path"
	"runtime"
	"testing"
)

func TestCompositeUnpacker_UnpackAll(t *testing.T) {
	ctrl := gomock.NewController(t)
	assert := assert.New(t)
	unpackers := []unpack.Unpacker{}
	for i := 0; i < 3; i++ {
		u := unpackmocks.NewMockUnpacker(ctrl)
		u.EXPECT().UnpackAll("test/")
		unpackers = append(unpackers, u)
	}
	unpacker := &unpack.CompositeUnpacker{Unpackers: unpackers}
	assert.NoError(unpacker.UnpackAll("test/"))

	errUnpacker := unpackmocks.NewMockUnpacker(ctrl)
	errUnpacker.EXPECT().UnpackAll("test/").Return(errors.New("fail"))

	unpacker = &unpack.CompositeUnpacker{Unpackers: []unpack.Unpacker{errUnpacker}}
	assert.EqualError(unpacker.UnpackAll("test/"), "fail")
}

func TestNewUnpacker(t *testing.T) {
	assert := assert.New(t)

	testCases := map[string]struct {
		files       []string
		expectFiles []string
		errf        string
	}{
		"a.zip": {
			files:       []string{"testdata/a.zip"},
			expectFiles: []string{"a.txt"},
		},
		"notzip.zip": {
			files: []string{"testdata/notzip.zip"},
			errf:  "unpacking %s/notzip.zip: zip: not a valid zip file",
		},
		"a.tar.gz": {
			files:       []string{"testdata/a.tar.gz"},
			expectFiles: []string{"a.txt", "b/b.txt"},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			unpacker := unpack.NewUnpacker(zap.NewNop().Sugar())
			tempDir := t.TempDir()
			for _, file := range tc.files {
				cp(t, file, path.Join(tempDir, path.Base(file)))
			}

			if err := unpacker.UnpackAll(tempDir); tc.errf == "" {
				assert.NoError(err)
				for _, expect := range tc.expectFiles {
					assert.FileExists(path.Join(tempDir, expect))
				}
			} else {
				assert.EqualError(err, fmt.Sprintf(tc.errf, tempDir))
			}
		})
	}
}

func cp(t *testing.T, src string, dst string) {
	// Assuming mainstream Linux, MacOS, Windows. This should fail on windows
	// machines only.
	assert.False(t, runtime.GOOS == "windows", "cp is not supported")
	cmd := exec.Command("cp", "-r", src, dst)
	err := cmd.Run()
	assert.NoError(t, err, cmd.String())
}
