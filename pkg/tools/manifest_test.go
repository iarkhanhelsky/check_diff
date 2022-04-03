package tools

import (
	"bytes"
	"crypto/sha256"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/buffer"
	"io/ioutil"
	"os/exec"
	"path"
	"runtime"
	"testing"
)

func TestManifest_isDiffer(t *testing.T) {
	assert := assert.New(t)
	dir := t.TempDir()
	cp(t, "testdata", dir)

	m := manifest{
		dir:          path.Join(dir, "testdata"),
		binaryDigest: sha256.Sum256([]byte("ABC")),
		logger:       zap.NewNop().Sugar(),
	}

	err := m.saveManifest()
	assert.NoError(err)

	assert.False(m.isDiffer())

	m.binaryDigest = sha256.Sum256([]byte("XXX"))
	assert.True(m.isDiffer())
}

func TestManifest_writeManifest(t *testing.T) {
	m := manifest{
		dir:          "testdata",
		binaryDigest: sha256.Sum256([]byte("ABC")),
		logger:       zap.NewNop().Sugar(),
	}

	var buf buffer.Buffer
	err := m.writeManifest(&buf)
	assert.NoError(t, err)
}

func TestManifest_saveManifest(t *testing.T) {
	assert := assert.New(t)
	dir := t.TempDir()
	cp(t, "testdata", dir)

	m := manifest{
		dir:          path.Join(dir, "testdata"),
		binaryDigest: sha256.Sum256([]byte("ABC")),
		logger:       zap.NewNop().Sugar(),
	}

	err := m.saveManifest()
	assert.NoError(err)

	data, _ := ioutil.ReadFile(path.Join(dir, "testdata", manifestFile))
	var buf bytes.Buffer
	err = m.writeManifest(&buf)
	assert.NoError(err)
	assert.Equal(buf.Bytes(), data)
}

func cp(t *testing.T, src string, dst string) {
	// Assuming mainstream Linux, MacOS, Windows. This should fail on windows
	// machines only.
	assert.False(t, runtime.GOOS == "windows", "cp is not supported")
	cmd := exec.Command("cp", "-r", src, dst)
	err := cmd.Run()
	assert.NoError(t, err, cmd.String())
}
