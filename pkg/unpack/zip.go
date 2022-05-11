package unpack

import (
	"context"
	"fmt"
	"github.com/saracen/fastzip"
	"go.uber.org/zap"
	"os"
	"path"
	"path/filepath"
)

type unzip struct {
	logger *zap.SugaredLogger
	dir    string
}

func (unpack *unzip) UnpackAll(dir string) error {
	files, _ := filepath.Glob(path.Join(dir, "*.zip"))
	for _, file := range files {
		unpack.logger.With("file", file).Debug("found archive")
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
