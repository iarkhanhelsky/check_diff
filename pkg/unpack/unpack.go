package unpack

import (
	"context"
	"fmt"
	"github.com/saracen/fastzip"
	"go.uber.org/multierr"
	"go.uber.org/zap"
	"os"
	"path"
	"path/filepath"
)

type Unpacker interface {
	UnpackAll(dir string) error
}

// for tests only
var _ Unpacker = &noopUnpacker{}

var _ Unpacker = &CompositeUnpacker{}
var _ Unpacker = &unzip{}

type noopUnpacker struct {
}

func NoopUnpacker() Unpacker {
	return &noopUnpacker{}
}

func (*noopUnpacker) UnpackAll(dir string) error {
	return nil
}

type CompositeUnpacker struct {
	Unpackers []Unpacker
	logger    zap.Logger
}

func NewUnpacker(logger *zap.SugaredLogger) Unpacker {
	return &CompositeUnpacker{
		Unpackers: []Unpacker{
			&unzip{logger: logger},
		},
	}
}

func (c CompositeUnpacker) UnpackAll(dir string) error {
	var err error
	for _, u := range c.Unpackers {
		err = multierr.Append(err, u.UnpackAll(dir))
	}

	return err
}

type unzip struct {
	logger *zap.SugaredLogger
	dir    string
}

func (unpack *unzip) UnpackAll(dir string) error {
	files, _ := filepath.Glob(path.Join(dir, "*.zip"))
	for _, file := range files {
		if err := unzipArchive(file, dir); err != nil {
			return err
		}
	}

	return nil
}

func unzipArchive(src string, dst string) error {
	defer os.Remove(src)
	// Create new extractor
	e, err := fastzip.NewExtractor(src, dst)
	if err != nil {
		return fmt.Errorf("unpacking %s: %w", src, err)
	}
	defer e.Close()

	// Extract archive files
	if err := e.Extract(context.Background()); err != nil {
		return fmt.Errorf("unpacking %s: %w", src, err)
	}
	return nil
}
