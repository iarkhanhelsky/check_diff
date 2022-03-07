package tools

import (
	"fmt"
	"github.com/pkg/errors"
	"os"
	"path"
	"runtime"
)

type TargetTuple string

const (
	Any         TargetTuple = "any"
	LinuxAMD64  TargetTuple = "linux-amd64"
	LinuxARM64  TargetTuple = "linux-arm64"
	DarwinAMD64 TargetTuple = "darwin-amd64"
	DarwinARM64 TargetTuple = "darwin-arm64"
	Current                 = TargetTuple(runtime.GOOS + "-" + runtime.GOARCH)
)

type TargetSource struct {
	Urls    []string
	MD5     string
	SHA256  string
	DstFile string
}

type Binary struct {
	Name    string
	DstFile string
	Path    string
	Targets map[TargetTuple]TargetSource

	executable string
}

var _ Downloader = &Binary{}

func (binary *Binary) Download(vendorDir string) error {

	return nil
}

func (binary *Binary) Executable() string {
	return binary.executable
}

func (binary Binary) selectSource() (TargetSource, error) {
	var targetSource TargetSource
	targetSource, ok := binary.Targets[Any]
	if ok {
		return targetSource, nil
	}

	targetSource, ok = binary.Targets[Current]
	if ok {
		return targetSource, nil
	}

	return targetSource, fmt.Errorf("failed to find target source; binary = %s", binary.Name)
}

func (binary *Binary) handleDownload(dstFolder string) error {
	executable := path.Join(dstFolder, binary.DstFile)

	if _, err := os.Stat(executable); err == os.ErrNotExist {
		return errors.Wrapf(err, "failed to setup %s", binary.Name)
	}

	if err := os.Chmod(executable, 0755); err != nil {
		return errors.Wrapf(err, "failed to setup %s", binary.Name)
	}
	binary.executable = executable
	return nil
}
