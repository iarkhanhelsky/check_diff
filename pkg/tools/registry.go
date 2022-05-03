package tools

import (
	errs "errors"
	"fmt"
	"github.com/iarkhanhelsky/check_diff/pkg/executors"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"os"
	"path"
	"path/filepath"
)

type Registry interface {
	Install(binaries ...*Binary) error
}

func NewRegistry(vendorDir string, logger *zap.SugaredLogger) Registry {
	return newRegistry(vendorDir,
		newUnpacker(logger.Named("unpacker")), logger)
}

func newRegistry(vendorDir string, unpacker unpacker, logger *zap.SugaredLogger) *registry {
	return &registry{
		logger:    logger.Named("registry"),
		unpacker:  unpacker,
		vendorDir: vendorDir,
	}
}

type registry struct {
	logger    *zap.SugaredLogger
	unpacker  unpacker
	vendorDir string
}

var _ Registry = &registry{}

func (registry *registry) Install(binaries ...*Binary) error {
	executor := executors.NewParallel()
	for _, b := range binaries {
		bx := b
		executor.Add(func() error {
			return registry.ensureBinary(bx)
		})
	}
	return executor.Run()
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
	return !newManifest(registry.binHome(binary), binary, registry.logger).isDiffer()
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

	dstFolder := path.Join(registry.binHome(binary), "_dist")
	logger.With("dist", dstFolder).Debug("create dist folder")
	if err := os.MkdirAll(dstFolder, 0755); err != nil {
		logger.With("dist", dstFolder).With("err", err).Error("failed to create dist")
		return errors.Wrap(err, "failed to create dist")
	}

	downloader := newHTTPDownloader(targetSource.MD5, targetSource.SHA256, targetSource.Urls...)
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
			if err := os.Remove(target); err != nil {
				logger.With("source", file).With("target", target).Error(err.Error())
				return fmt.Errorf("failed to remove %s: %w", target, err)
			}
		}
		logger.With("source", file).With("target", target).Debug("create symlink")
		if err := os.Symlink(source, target); err != nil {
			return fmt.Errorf("failed to create symlink: %w", err)
		}
	}

	return nil
}

func (registry registry) unpack(binary *Binary) error {
	logger := registry.logger.With("binary", binary.Name)
	logger.Debug("unpacking binary")
	return registry.unpacker.unpackAll(registry.binHome(binary))
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
	if _, err := os.Stat(binary.executable); errs.Is(err, os.ErrNotExist) {
		return fmt.Errorf("unexpected %s: %w", binary.Name, err)
	}
	logger.With("executable", binary.executable).Debug("binary path updated")
	return nil
}
