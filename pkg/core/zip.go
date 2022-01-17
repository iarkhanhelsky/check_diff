package core

import (
	"context"
	"github.com/saracen/fastzip"
	"path"
)

var _ Downloader = &zipDownloader{}

type zipDownloader struct {
	inner   Downloader
	zipDir  string
	dstFile string
	handler DownloadHandler
}

func NewZipDownloader(handler DownloadHandler, dstFile string, md5 string, sha256 string, urls ...string) Downloader {
	d := zipDownloader{dstFile: dstFile, handler: handler}
	d.inner = NewDownloader(d.handleDownload, dstFile+".zip", md5, sha256, urls...)
	return &d
}

func (downloader *zipDownloader) Download(dstFolder string) error {
	return downloader.inner.Download(dstFolder)
}

func (downloader *zipDownloader) handleDownload(resultDir string) error {
	err := unzip(path.Join(resultDir, downloader.dstFile+".zip"), path.Join(resultDir, downloader.dstFile))
	if err != nil {
		return err
	}
	return downloader.handler(resultDir)
}

func unzip(src string, dst string) error {
	// Create new extractor
	e, err := fastzip.NewExtractor(src, dst)
	if err != nil {
		panic(err)
	}
	defer e.Close()

	// Extract archive files
	return e.Extract(context.Background())
}
