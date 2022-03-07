package tools

import (
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"os"
	"path"
	"path/filepath"
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
	for _, binary := range binaries {
		if err := registry.ensureBinary(binary); err != nil {
			return err
		}
	}

	return nil
}

func (registry *registry) ensureBinary(binary *Binary) error {
	logger := registry.logger.With("binary", binary.Name)
	logger.Debug("ensuring binary")
	if err := registry.fetch(binary); err != nil {
		logger.With("err", err).Errorf("fetch failed")
		return errors.Wrapf(err, "fetch failed; binary = %s", binary.Name)
	}
	if err := registry.unpack(binary); err != nil {
		logger.With("err", err).Errorf("unpack failed")
		return errors.Wrapf(err, "unpack failed; binary = %s", binary.Name)
	}
	if err := registry.install(binary); err != nil {
		logger.With("err", err).Errorf("install failed")
		return errors.Wrapf(err, "install failed; binary = %s", binary.Name)
	}

	return nil
}

func (registry registry) fetch(binary *Binary) error {
	logger := registry.logger.With("binary", binary.Name)
	logger.Debug("fetching binary")
	handler := func(_ string) error { return nil }
	targetSource, err := binary.selectSource()
	if err != nil {
		return err
	}
	downloader := NewHTTPDownloader(handler, binary.DstFile,
		targetSource.MD5, targetSource.SHA256, targetSource.Urls...)

	dstFolder := path.Join(registry.vendorDir, binary.Name, "_dist")
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
		target := path.Join(registry.vendorDir, binary.Name, path.Base(file))
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
	return unpack{logger: logger, dir: path.Join(registry.vendorDir, binary.Name)}.unpackAll()
}

func (registry registry) install(binary *Binary) error {
	logger := registry.logger.With("binary", binary.Name)
	logger.Debug("installing binary")
	binary.executable = path.Join(registry.vendorDir, binary.Name, binary.Path)
	logger.With("executable", binary.executable).Debug("binary path updated")
	return nil
}
