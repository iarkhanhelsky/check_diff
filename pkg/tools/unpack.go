package tools

import (
	"context"
	"github.com/saracen/fastzip"
	"go.uber.org/zap"
	"os"
	"path"
	"path/filepath"
)

type unpack struct {
	logger *zap.SugaredLogger
	dir    string
}

func (unpack unpack) unpackAll() error {
	files, _ := filepath.Glob(path.Join(unpack.dir, "*.zip"))
	for _, file := range files {
		if err := unzip(file, unpack.dir); err != nil {
			return err
		}
	}

	return nil
}

func unzip(src string, dst string) error {
	defer os.Remove(src)
	// Create new extractor
	e, err := fastzip.NewExtractor(src, dst)
	if err != nil {
		panic(err)
	}
	defer e.Close()

	// Extract archive files
	return e.Extract(context.Background())
}
