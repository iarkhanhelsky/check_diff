package tools

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

type unpacker interface {
	unpackAll(dir string) error
}

// for tests only
var _ unpacker = &noopUnpacker{}

var _ unpacker = &compositeUnpacker{}
var _ unpacker = &unzip{}

type noopUnpacker struct {
}

func (*noopUnpacker) unpackAll(dir string) error {
	return nil
}

type compositeUnpacker struct {
	unpackers []unpacker
	logger    zap.Logger
}

func newUnpacker(logger *zap.SugaredLogger) unpacker {
	return &compositeUnpacker{
		unpackers: []unpacker{
			&unzip{logger: logger},
		},
	}
}

func (c compositeUnpacker) unpackAll(dir string) error {
	var err error
	for _, u := range c.unpackers {
		err = multierr.Append(err, u.unpackAll(dir))
	}

	return err
}

type unzip struct {
	logger *zap.SugaredLogger
	dir    string
}

func (unpack *unzip) unpackAll(dir string) error {
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
