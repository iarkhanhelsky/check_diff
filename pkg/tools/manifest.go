package tools

import (
	"bufio"
	"bytes"
	"crypto/md5"
	"crypto/sha256"
	"encoding/gob"
	"fmt"
	"go.uber.org/zap"
	"io"
	"io/fs"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
)

const (
	manifestFile = "_manifest"
)

type manifestEntry struct {
	Path       string
	FileDigest []byte
}

func fileDigest(path string) []byte {
	file, err := os.Open(path)
	if err != nil {
		return []byte{}
	}
	data, err := ioutil.ReadAll(file)
	digest := sha256.Sum256(data)
	return digest[:]
}

func newManifestEntry(path string) manifestEntry {
	digest := fileDigest(path)
	return manifestEntry{Path: path, FileDigest: digest}
}

type manifest struct {
	dir          string
	binaryDigest [32]byte
	logger       *zap.SugaredLogger
}

func newManifest(dir string, binary *Binary, logger *zap.SugaredLogger) manifest {
	digest, _ := binary.digest()
	return manifest{
		dir:          dir,
		binaryDigest: digest,
		logger:       logger.Named("manifest").With("binary", binary.Name),
	}
}

func (m manifest) saveManifest() error {
	logger := m.logger
	manifestFile := path.Join(m.dir, manifestFile)
	logger.With("path", manifestFile).Debug(manifestFile, "saving manifest file")
	file, err := os.Create(manifestFile)
	defer file.Close()
	if err != nil {
		logger.With("err", err).Error("failed to save manifest")
		return fmt.Errorf("writing manifest: %w", err)
	}

	writer := bufio.NewWriter(file)
	if err := m.writeManifest(writer); err != nil {
		return fmt.Errorf("writing manifest: %w", err)
	}

	return writer.Flush()
}

func (m manifest) writeManifest(w io.Writer) error {
	encoder := gob.NewEncoder(w)
	logger := m.logger
	digest := m.binaryDigest
	logger.With("digest", fmt.Sprintf("%x", digest)).Debug("binary digest")
	if err := encoder.Encode(digest); err != nil {
		return err
	}
	return filepath.Walk(m.dir, func(prefix string, info fs.FileInfo, err error) error {
		if info.Name() == manifestFile || info.IsDir() {
			return nil
		}
		entry := newManifestEntry(prefix)
		logger.
			With("path", entry.Path).
			With("digest", fmt.Sprintf("%x", entry.FileDigest)).
			Debug("new manifest entry")

		return encoder.Encode(entry)
	})
}

func (m manifest) isDiffer() bool {
	logger := m.logger

	manifestFile := path.Join(m.dir, manifestFile)
	logger.With("path", manifestFile).Debug("reading manifest file")
	data, err := ioutil.ReadFile(manifestFile)
	if err != nil {
		logger.With("err", err).Error("failed to read manifest file: differ = true")
		return true
	}

	var buffer bytes.Buffer
	logger.Debug("creating uptodate manifest blob")
	if err := m.writeManifest(&buffer); err != nil {
		logger.With("err", err).Error("failed create digest blob: differ = true")
		return true
	}

	if bytes.Compare(data, buffer.Bytes()) == 0 {
		return false
	}

	logger.
		With("md5", fmt.Sprintf("%x", md5.Sum(data))).
		With("size", len(data)).
		Debug("actual")
	logger.
		With("md5", fmt.Sprintf("%x", md5.Sum(buffer.Bytes()))).
		With("size", len(buffer.Bytes())).
		Debug("expected")

	return true
}
