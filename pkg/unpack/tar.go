package unpack

import (
	"archive/tar"
	"compress/gzip"
	"errors"
	"fmt"
	"go.uber.org/zap"
	"io"
	"os"
	"path"
	"path/filepath"
	"strings"
)

type untar struct {
	logger *zap.SugaredLogger
	dir    string
}

func (unpack *untar) UnpackAll(dir string) error {
	files, _ := filepath.Glob(path.Join(dir, "*.tar.gz"))
	for _, file := range files {
		unpack.logger.With("file", file).Debug("found archive")
		reader, err := os.Open(file)
		defer reader.Close()

		if err != nil {
			return fmt.Errorf("reading archive file=%s: %w", file, err)
		}
		var stream io.Reader
		if strings.HasSuffix(file, ".gz") {
			stream, err = gzip.NewReader(reader)
			if err != nil {
				return fmt.Errorf("reading archive file=%s: %w", file, err)
			}
		}
		if err := untarArchive(stream, dir); err != nil {
			return err
		}
	}

	return nil
}

func untarArchive(stream io.Reader, dst string) error {
	tarReader := tar.NewReader(stream)
	for true {
		header, err := tarReader.Next()

		if err == io.EOF {
			break
		}

		if err != nil {
			return fmt.Errorf("reading tar: %w", err)
		}

		dstPath := path.Join(dst, header.Name)

		switch header.Typeflag {
		case tar.TypeDir:
			if err := os.Mkdir(dstPath, 0755); err != nil {
				return fmt.Errorf("reading tar: mkdir %s: %w", dstPath, err)
			}
		case tar.TypeReg:
			if _, err := os.Stat(path.Dir(dstPath)); errors.Is(err, os.ErrNotExist) {
				if err := os.MkdirAll(path.Dir(dstPath), 0755); err != nil {
					return fmt.Errorf("reading tar: mkdir %s: %w", path.Dir(dstPath), err)
				}

			}
			outFile, err := os.Create(dstPath)
			defer outFile.Close()
			if err != nil {
				return fmt.Errorf("reading tar: create file %s: %w", dstPath, err)
			}
			if _, err := io.Copy(outFile, tarReader); err != nil {
				return fmt.Errorf("reading tar: write file %s: %w", dstPath, err)
			}

		default:
			return fmt.Errorf("reading tar: unknown entry type %d of %s", header.Typeflag, header.Name)
		}
	}

	return nil
}
