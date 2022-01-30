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
	"os"
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

	dstFile := path.Join(dstFolder, d.dstFile)
	if !d.isUptodate(dstFile) {
		for _, url := range d.urls {
			err := d.downloadFrom(url, dstFile)
			if err == nil {
				break
			}

			accumulatedErrors = append(accumulatedErrors, err)
		}

		if len(accumulatedErrors) > 0 {
			// TODO: Pretty error output
			return fmt.Errorf("%v", accumulatedErrors)
		}
	}

	return d.handler(dstFolder)
}

func (d *downloader) downloadFrom(url string, outputFile string) error {
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

	return nil
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

func (d *downloader) isUptodate(dstFile string) bool {
	if len(d.md5) == 0 && len(d.sha256) == 0 {
		// No way to check, always outdated
		return false
	}

	if _, err := os.Stat(dstFile); errors.Is(err, os.ErrNotExist) {
		return false
	}

	file, err := os.Open(dstFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "WARN: Failed to check %s: %v", dstFile, err)
		return false
	}
	defer file.Close()

	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Fprintf(os.Stderr, "WARN: Failed to check %s: %v", dstFile, err)
		return false
	}

	return checkMD5(bytes, d.md5) == nil && checkSHA256(bytes, d.sha256) == nil
}
