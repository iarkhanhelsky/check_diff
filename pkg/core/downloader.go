package core

import (
	"crypto/md5"
	"crypto/sha256"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"io/ioutil"
	"net/http"
	"path"
)

type Downloader interface {
	Download(dstFolder string) error
}

type DownloadHandler func(path string) error

type downloader struct {
	urls    []string
	md5     string
	sha256  string
	dstFile string
	handler DownloadHandler
}

func NewDownloader(handler DownloadHandler, dstFile string, md5 string, sha256 string, urls ...string) Downloader {
	return &downloader{urls, md5, sha256, dstFile, handler}
}

func (d *downloader) Download(dstFolder string) error {
	var accumulatedErrors []error

	for _, url := range d.urls {
		err := d.downloadFrom(url, dstFolder)
		if err == nil {
			return nil
		}

		accumulatedErrors = append(accumulatedErrors, err)
	}

	// TODO: Pretty error output
	return fmt.Errorf("%v", accumulatedErrors)
}

func (d *downloader) downloadFrom(url string, dstFolder string) error {
	outputFile := path.Join(dstFolder, d.dstFile)

	resp, err := http.Get(url)
	defer resp.Body.Close()
	if err != nil {
		return err
	}

	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if err := checkMD5(bytes, d.md5); err != nil {
		return err
	}

	if err := checkSHA256(bytes, d.sha256); err != nil {
		return err
	}

	if err := ioutil.WriteFile(outputFile, bytes, fs.FileMode(int(0777))); err != nil {
		return err
	}

	return d.handler(dstFolder)
}

func (d *downloader) ensurePath() (string, error) {
	return "", nil
}

func checkMD5(bytes []byte, md5sum string) error {
	// "" means skip
	if len(md5sum) == 0 {
		return nil
	}

	hex := fmt.Sprintf("%x", md5.Sum(bytes))
	if hex != md5sum {
		return errors.New("md5sum mismatch")
	}

	return nil
}

func checkSHA256(bytes []byte, sha1sum string) error {
	// "" means skip
	if len(sha1sum) == 0 {
		return nil
	}

	hex := fmt.Sprintf("%x", sha256.Sum256(bytes))
	if hex != sha1sum {
		return errors.New("sha1sum mismatch")
	}

	return nil
}
