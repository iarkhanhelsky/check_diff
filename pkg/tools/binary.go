package tools

import (
	"crypto/sha256"
	stbinary "encoding/binary"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"runtime"
)

type TargetTuple string

// When this version changes all binaries will be re-downloaded
const serialVersion uint64 = 1

const (
	Any         TargetTuple = "any"
	LinuxAMD64  TargetTuple = "linux-amd64"
	LinuxARM64  TargetTuple = "linux-arm64"
	DarwinAMD64 TargetTuple = "darwin-amd64"
	DarwinARM64 TargetTuple = "darwin-arm64"
	Current                 = TargetTuple(runtime.GOOS + "-" + runtime.GOARCH)
)

type TargetSource struct {
	Urls    []string `json:"urls"`
	MD5     string   `json:"md5"`
	SHA256  string   `json:"sha256"`
	DstFile string   `json:"dstFile"`
}

type Binary struct {
	Name    string                       `json:"name"`
	Path    string                       `json:"path"`
	Targets map[TargetTuple]TargetSource `json:"targets"`
	DstFile string                       `json:"dstFile"`

	executable string
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

func (binary Binary) digest() ([32]byte, error) {
	data, err := json.Marshal(&binary)
	if err != nil {
		return [32]byte{}, errors.Wrapf(err, "failed to calculate digest binary=%s", binary.Name)
	}
	stbinary.LittleEndian.PutUint64(data, serialVersion)

	return sha256.Sum256(data), nil
}
