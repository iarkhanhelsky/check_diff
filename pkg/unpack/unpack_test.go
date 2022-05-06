package unpack

import (
	"github.com/golang/mock/gomock"
	"github.com/iarkhanhelsky/check_diff/mocks"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"os/exec"
	"path"
	"path/filepath"
	"runtime"
	"testing"
)

func TestCompositeUnpacker_UnpackAll(t *testing.T) {
	ctrl := gomock.NewController(t)
	assert := assert.New(t)
	unpackers := []Unpacker{}
	for i := 0; i < 3; i++ {
		u := mocks.NewMockUnpacker(ctrl)
		u.EXPECT().UnpackAll("test/")
		unpackers = append(unpackers, u)
	}
	unpacker := &compositeUnpacker{unpackers: unpackers}
	assert.NoError(unpacker.UnpackAll("test/"))

	errUnpacker := mocks.NewMockUnpacker(ctrl)
	errUnpacker.EXPECT().UnpackAll("test/").Return(errors.New("fail"))

	unpacker = &compositeUnpacker{unpackers: []Unpacker{errUnpacker}}
	assert.EqualError(unpacker.UnpackAll("test/"), "fail")
}

func TestNewUnpacker(t *testing.T) {
	assert := assert.New(t)

	unpacker := NewUnpacker(zap.NewNop().Sugar())
	tempDir := t.TempDir()
	files, err := filepath.Glob("testdata/*")
	assert.NoError(err)
	for _, file := range files {
		cp(t, file, path.Join(tempDir, path.Base(file)))
	}

	assert.NoError(unpacker.UnpackAll(tempDir))
	assert.FileExists(path.Join(tempDir, "a.txt"))
}

func cp(t *testing.T, src string, dst string) {
	// Assuming mainstream Linux, MacOS, Windows. This should fail on windows
	// machines only.
	assert.False(t, runtime.GOOS == "windows", "cp is not supported")
	cmd := exec.Command("cp", "-r", src, dst)
	err := cmd.Run()
	assert.NoError(t, err, cmd.String())
}
