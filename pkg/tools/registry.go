package tools

import (
	"github.com/pkg/errors"
	"go.uber.org/multierr"
	"go.uber.org/zap"
	"os"
	"path"
	"path/filepath"
	"sync"
)

type Registry interface {
	Install(binaries ...*Binary) error
}

func NewRegistry(logger *zap.SugaredLogger, vendorDir string) Registry {
	return &registry{
		logger:    logger.Named("registry"),
		vendorDir: vendorDir,
	}
}

type registry struct {
	logger    *zap.SugaredLogger
	vendorDir string
}

var _ Registry = &registry{}

func (registry *registry) Install(binaries ...*Binary) error {
	var wg sync.WaitGroup
	wg.Add(len(binaries))
	errchan := make(chan error)
	done := make(chan bool)
	for _, binary := range binaries {
		bin := binary
		go func() {
			errchan <- registry.ensureBinary(bin)
			wg.Done()
		}()
	}

	var err error
	go func() {
		for {
			select {
			case <-done:
				return
			case e := <-errchan:
				err = multierr.Append(err, e)
			}
		}
	}()

	wg.Wait()
	done <- true

	return err
}

func (registry *registry) ensureBinary(binary *Binary) error {
	logger := registry.logger.With("binary", binary.Name)

	logger.Debug("check if we already have this binary")
	if registry.uptodate(binary) {
		logger.Debug("binary is up to date")
		goto install
	}

	if err := os.MkdirAll(registry.binHome(binary), 0755); err != nil {
		logger.With("err", err, "dir", registry.binHome(binary)).Error("failed to create directory")
		return err
	}

	logger.Debug("ensuring binary")
	if err := registry.fetch(binary); err != nil {
		logger.With("err", err).Errorf("fetch failed")
		return errors.Wrapf(err, "fetch failed; binary = %s", binary.Name)
	}

	if err := registry.unpack(binary); err != nil {
		logger.With("err", err).Errorf("unpack failed")
		return errors.Wrapf(err, "unpack failed; binary = %s", binary.Name)
	}

	if err := registry.writeManifest(binary); err != nil {
		logger.With("err", err).Errorf("failed to write digest")
		return errors.Wrapf(err, "digest failed; binary = %s", binary.Name)
	}

install:
	if err := registry.install(binary); err != nil {
		logger.With("err", err).Errorf("install failed")
		return errors.Wrapf(err, "install failed; binary = %s", binary.Name)
	}

	return nil
}

func (registry registry) uptodate(binary *Binary) bool {
	logger := registry.logger.With("binary", binary.Name)
	logger.Debug("checking digest")
	return newManifest(registry.binHome(binary), binary, registry.logger).isDiffer()
}

func (registry registry) binHome(binary *Binary, other ...string) string {
	return path.Join(append([]string{registry.vendorDir, binary.Name}, other...)...)
}

func (registry registry) fetch(binary *Binary) error {
	logger := registry.logger.With("binary", binary.Name)
	logger.Debug("fetching binary")
	targetSource, err := binary.selectSource()
	if err != nil {
		return err
	}
	downloader := newHTTPDownloader(targetSource.MD5, targetSource.SHA256, targetSource.Urls...)

	dstFolder := path.Join(registry.binHome(binary), "_dist")
	logger.With("dist", dstFolder).Debug("create dist folder")
	if err := os.MkdirAll(dstFolder, 0755); err != nil {
		logger.With("dist", dstFolder).With("err", err).Error("failed to create dist")
		return errors.Wrap(err, "failed to create dist")
	}

	if err := downloader.Download(dstFolder); err != nil {
		logger.With("dist", dstFolder).With("err", err).Error("download failed")
		return errors.Wrap(err, "download failed")
	}

	files, _ := filepath.Glob(path.Join(dstFolder, "*"))
	for _, file := range files {
		source, _ := filepath.Abs(file)
		fileName := targetSource.DstFile
		if fileName == "" {
			fileName = binary.DstFile
		}

		target := path.Join(registry.binHome(binary), fileName)
		if _, err := os.Stat(target); err == nil {
			logger.With("source", file).With("target", target).Debug("exists removing")
			continue
		}
		logger.With("source", file).With("target", target).Debug("create symlink")
		if err := os.Symlink(source, target); err != nil {
			return errors.Wrap(err, "failed to create symlink")
		}
	}

	return nil
}

func (registry registry) unpack(binary *Binary) error {
	logger := registry.logger.With("binary", binary.Name)
	logger.Debug("unpacking binary")
	return unpack{logger: logger, dir: registry.binHome(binary)}.unpackAll()
}

func (registry registry) writeManifest(binary *Binary) error {
	logger := registry.logger.With("binary", binary.Name)
	logger.Debug("writing manifest")
	return newManifest(registry.binHome(binary), binary, registry.logger).saveManifest()
}

func (registry registry) install(binary *Binary) error {
	logger := registry.logger.With("binary", binary.Name)
	logger.Debug("installing binary")
	binary.executable = path.Join(registry.binHome(binary), binary.Path)
	logger.With("executable", binary.executable).Debug("binary path updated")
	return nil
}
